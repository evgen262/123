package auth

import (
	"context"
	"net"
	"net/url"
	"time"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/auth/v1"
	authErrorsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/errors/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
)

type authRepository struct {
	client      authv1.AuthAPIClient
	callbackURL url.URL
	mapper      MapperAuth
	appName     string
	accessTTL   time.Duration
	refreshTTL  time.Duration
	tu          timeUtils.TimeUtils
	logger      ditzap.Logger
}

func NewAuthRepository(
	client authv1.AuthAPIClient,
	mapper MapperAuth,
	redirectURL url.URL,
	appName string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	tu timeUtils.TimeUtils,
	logger ditzap.Logger,
) *authRepository {
	return &authRepository{
		client:      client,
		callbackURL: redirectURL,
		mapper:      mapper,
		appName:     appName,
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		tu:          tu,
		logger:      logger,
	}
}

func (a *authRepository) buildURL(callbackUri string) string {
	u := a.callbackURL
	uri, err := url.Parse(callbackUri)
	if err != nil {
		a.logger.Warn("cant parse uri", zap.String("callback_uri", callbackUri), zap.Error(err))
		return u.String()
	}

	vals := u.Query()
	for k, vs := range uri.Query() {
		for _, v := range vs {
			vals.Add(k, v)
		}
	}

	newURL := u.JoinPath(uri.Path)
	newURL.RawQuery = vals.Encode()
	newURL.Fragment = uri.Fragment

	return newURL.String()
}

func (a *authRepository) GetRedirectURL(ctx context.Context, callbackURI string) (string, error) {
	resp, err := a.client.GetURL(ctx, &authv1.GetURLRequest{
		CallbackUrl: a.buildURL(callbackURI),
	})
	if err != nil {
		return "", diterrors.GrpcErrorToError(err)
	}

	return resp.GetRedirectUrl(), nil
}

func (a *authRepository) Auth(ctx context.Context, code, state, callbackURI string) (*entityAuth.AuthSudir, error) {
	resp, err := a.client.Auth(ctx, &authv1.AuthRequest{
		State:       state,
		Code:        code,
		CallbackUrl: a.buildURL(callbackURI),
	})
	if err != nil {
		return nil, diterrors.GrpcErrorToError(err) //nolint:wrapcheck
	}

	return &entityAuth.AuthSudir{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		User:         a.mapper.UserToEntity(resp.GetUser()),
	}, nil
}

func (a *authRepository) AuthPortal(ctx context.Context, params entityAuth.AuthPortalParams) (*entityAuth.Auth1C, error) {
	resp, err := a.client.Auth1C(
		ctx,
		&authv1.Auth1CRequest{
			PortalUrl: params.PortalURL,
			// SessionId = zero value указывает, что пользователя надо аутентифицировать на портале
			// SessionId: params.PortalSession,
			User: &authv1.User1C{
				CloudId:  params.User.CloudID,
				Snils:    params.User.SNILS,
				Email:    params.User.Email,
				UserType: authv1.User1C_USER_TYPE_WEB,
			},
		},
	)
	if err != nil {
		return nil, diterrors.GrpcErrorToError(err) //nolint:wrapcheck
	}

	return &entityAuth.Auth1C{
		PortalSession: resp.GetSessionId(),
		EmployeeID:    resp.GetEmployeeId(),
		PersonID:      resp.GetPersonId(),
	}, nil
}

func (a *authRepository) CreateSession(ctx context.Context, user *entityAuth.UserSudir, clientIP net.IP, device *entityAuth.Device, auth *entityAuth.Auth1C) (entityAuth.TokensPair, error) {
	if len(user.Portals) == 0 {
		return entityAuth.TokensPair{}, repositories.ErrNoUserPortals
	}

	resp, err := a.client.CreateSession(ctx, &authv1.CreateSessionRequest{
		User: &authv1.CreateSessionRequest_User{
			// TODO Временно передаем nil uuid, позднее будет заменен на идентификатор пользователя
			Id:        uuid.Nil.String(),
			CloudId:   user.CloudID,
			LogonName: user.Login,
			Email:     user.Email,
			Snils:     user.SNILS,
			Portal: &authv1.CreateSessionRequest_UserPortal{
				Id:   int32(user.Portals[0].ID),
				Name: user.Portals[0].Name,
				Url:  user.Portals[0].URL,
				Sid:  auth.GetPortalSession(),
			},
			Employee: &authv1.CreateSessionRequest_Employee{
				Id: auth.GetEmployeeID(),
			},
			Person: &authv1.CreateSessionRequest_Person{
				Id: auth.GetPersonID(),
			},
		},
		UserType:           authv1.UserType_USER_TYPE_AUTH,
		UserIp:             clientIP.String(),
		Device:             &authv1.CreateSessionRequest_Device{UserAgent: device.GetUserAgent()},
		SudirInfo:          &authv1.CreateSessionRequest_SudirInfo{Sid: user.SID},
		Issuer:             a.appName,
		Subject:            user.Login,
		AccessTtlDuration:  durationpb.New(a.accessTTL),
		RefreshTtlDuration: durationpb.New(a.refreshTTL),
	})
	if err != nil {
		return entityAuth.TokensPair{}, diterrors.GrpcErrorToError(err)
	}

	accessExpiredTime := resp.GetTokens().GetAccessToken().GetExpiredTime().AsTime()
	refreshExpiredTime := resp.GetTokens().GetRefreshToken().GetExpiredTime().AsTime()

	return entityAuth.TokensPair{
		AccessToken: entityAuth.Token{
			Value:     resp.GetTokens().GetAccessToken().GetValue(),
			ExpiredAt: &accessExpiredTime,
		},
		RefreshToken: entityAuth.Token{
			Value:     resp.GetTokens().GetRefreshToken().GetValue(),
			ExpiredAt: &refreshExpiredTime,
		},
	}, nil
}

func (a *authRepository) GetSession(ctx context.Context, accessToken string) (*entityAuth.Session, error) {
	if accessToken == "" {
		return nil, diterrors.NewValidationError(repositories.ErrEmptyAccessToken)
	}

	resp, err := a.client.GetSession(ctx, &authv1.GetSessionRequest{AccessToken: accessToken})
	if err != nil {
		return nil, diterrors.GrpcErrorToError(err)
	}

	session := a.mapper.SessionToEntity(resp.GetSession())
	if session == nil {
		return nil, diterrors.ErrNotFound
	}

	return session, nil
}

func (a *authRepository) Logout(ctx context.Context, session *entityAuth.Session, accessToken, refreshToken string) error {
	if session == nil {
		return diterrors.NewValidationError(repositories.ErrNilSession)
	}
	if accessToken == "" {
		return diterrors.NewValidationError(repositories.ErrEmptyAccessToken)
	}
	if refreshToken == "" {
		return diterrors.NewValidationError(repositories.ErrEmptyRefreshToken)
	}
	if session.GetActivePortal() == nil {
		return repositories.ErrNilSessionActivePortal
	}
	if session.GetUser() == nil {
		return repositories.ErrNilSessionUser
	}
	if session.GetDevice().GetSudirInfo() == nil {
		return repositories.ErrNilSessionSudirInfo
	}

	// Ответ заглушен, потому что отсутствует
	_, err := a.client.Logout(ctx, &authv1.LogoutRequest{
		PortalUrl: session.GetActivePortal().GetPortal().URL,
		SessionId: session.GetID().String(),
		LogoutBy: &authv1.LogoutRequest_CloudId{
			CloudId: session.GetUser().CloudID,
		},
		SudirSid:      session.GetDevice().GetSudirInfo().SID,
		PortalSession: session.GetActivePortal().SID,
		Tokens: &authv1.LogoutRequest_TokensPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
	if err != nil {
		return diterrors.GrpcErrorToError(err)
	}

	return nil
}

func (a *authRepository) ChangePortal(ctx context.Context, portalID int, session *entityAuth.Session) ([]*entityAuth.Portal, string, error) {
	if session == nil {
		return nil, "", repositories.ErrNilSession
	}

	resp, err := a.client.ChangePortal(ctx, &authv1.ChangePortalRequest{
		SelectedPortalId: int32(portalID),
		Session:          a.mapper.SessionToPb(session),
	})
	if err != nil {
		return nil, "", diterrors.GrpcErrorToError(err)
	}

	return a.mapper.PortalsToEntity(resp.GetPortals()), resp.GetPortalSid(), nil
}

func (a *authRepository) RefreshTokensPair(ctx context.Context, accessToken, refreshToken string) (*entityAuth.TokensPair, error) {
	if accessToken == "" {
		return nil, repositories.NewDetailsError("access_token", repositories.ErrEmptyAccessToken.Error(), false)
	}
	if refreshToken == "" {
		return nil, repositories.NewDetailsError("refresh_token", repositories.ErrEmptyRefreshToken.Error(), false)
	}

	resp, err := a.client.RefreshTokensPair(ctx, &authv1.RefreshTokensPairRequest{
		Tokens: &authv1.RefreshTokensPairRequest_TokensPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})
	if err != nil {
		st, _ := status.FromError(err)
		for _, detail := range st.Details() {
			if detailErr, ok := detail.(*authErrorsv1.AuthInvalidArgument); ok {
				return nil, repositories.NewDetailsError(detailErr.GetField(), detailErr.GetMessage(), detailErr.GetReauthRequired())
			}
		}
		return nil, diterrors.GrpcErrorToError(err)
	}
	return &entityAuth.TokensPair{
		AccessToken: entityAuth.Token{
			Value:     resp.GetTokens().GetAccessToken().GetValue(),
			ExpiredAt: a.tu.TimestampToTime(resp.GetTokens().GetAccessToken().GetExpiredTime()),
		},
		RefreshToken: entityAuth.Token{
			Value:     resp.GetTokens().GetRefreshToken().GetValue(),
			ExpiredAt: a.tu.TimestampToTime(resp.GetTokens().GetRefreshToken().GetExpiredTime()),
		},
	}, nil
}
