package grpc

import (
	"context"
	"errors"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/view"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/kadry"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/client/http/sudir"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/usecase"
)

type authServer struct {
	interactor AuthInteractor
	presenter  AuthPresenter
	authv1.UnimplementedAuthAPIServer
}

func NewAuthServer(interactor AuthInteractor, presenter AuthPresenter) *authServer {
	return &authServer{
		interactor: interactor,
		presenter:  presenter,
	}
}

// GetURL
//
//	URL для перенаправления пользователя для авторизации в СУДИР.
func (a *authServer) GetURL(ctx context.Context, request *authv1.GetURLRequest) (*authv1.GetURLResponse, error) {
	callbackURL := request.GetCallbackUrl()
	clientID := ""
	clientSecret := ""

	if request.GetCredentials() != nil {
		clientID = request.GetCredentials().GetClientId()
		clientSecret = request.GetCredentials().GetClientSecret()

		if clientID == "" || clientSecret == "" {
			return nil, NewApiError(codes.InvalidArgument, "bad credentials provided").WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			dataDeviceID := md.Get("deviceid")
			if len(dataDeviceID) > 0 {
				deviceID := dataDeviceID[0]
				ctx = context.WithValue(ctx, "deviceid", deviceID)
			}

			dataXCFCUserAgent := md.Get("x-cfc-useragent")
			if len(dataXCFCUserAgent) > 0 {
				ua := dataXCFCUserAgent[0]
				ctx = context.WithValue(ctx, "x-cfc-useragent", ua)
			}
		}

	}

	authURL, err := a.interactor.GetAuthURL(ctx, callbackURL, clientID, clientSecret)
	if err != nil {
		return nil, NewApiError(codes.Internal, "auth service error", err).WithLocalizedMessage(view.ErrStatusTextServiceError)
	}

	return &authv1.GetURLResponse{
		RedirectUrl: authURL,
	}, nil
}

// Auth авторизация пользователя для сервиса
//
//	метод возвращает информацию о пользователе в СУДИР и сотрудниках из СКС
func (a *authServer) Auth(ctx context.Context, request *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	code := request.GetCode()
	state := request.GetState()
	callbackURL := request.GetCallbackUrl()

	if code == "" {
		return nil, NewApiError(codes.InvalidArgument, "code not provided").WithLocalizedMessage(view.ErrStatusTextCodeNotProvided)
	}

	if callbackURL == "" {
		return nil, NewApiError(codes.InvalidArgument, "callback url not provided").WithLocalizedMessage(view.ErrStatusTextCallbackURLNotProvided)
	}

	authInfo, err := a.interactor.Auth(ctx, code, state, callbackURL)
	if err != nil {
		switch {
		case errors.Is(err, kadry.ErrEmployeesService):
			break
		case errors.Is(err, usecase.ErrCallbackURLMismatch):
			return nil, NewApiError(codes.InvalidArgument, "callback url mismatch").WithLocalizedMessage(view.ErrStatusTextCallbackURLMismatch)
		case errors.Is(err, sudir.ErrInvalidGrant):
			return nil, NewApiError(codes.InvalidArgument, "invalid code provided").WithLocalizedMessage(view.ErrStatusTextBadCodeProvided)
		case errors.Is(err, repositories.ErrInvalidState):
			fallthrough
		case errors.Is(err, repositories.ErrStateNotFound):
			return nil, NewApiError(codes.InvalidArgument, "invalid state provided").WithLocalizedMessage(view.ErrStatusTextStateInvalid)
		default:
			return nil, NewApiError(codes.Internal, "auth service error", err).WithLocalizedMessage(view.ErrStatusTextServiceError)
		}
	}

	return &authv1.AuthResponse{
		User:        a.presenter.UserToPb(authInfo.User),
		AccessToken: authInfo.OAuth.AccessToken,
		Device:      a.presenter.DeviceToPb(authInfo.Device),
	}, nil
}

func (a *authServer) Login(ctx context.Context, request *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	var (
		info *entity.AuthInfo
		err  error
	)

	switch loginBy := request.LoginBy.(type) {
	case *authv1.LoginRequest_Credentials:
		credentials := loginBy.Credentials
		if credentials == nil {
			return nil, NewApiError(codes.InvalidArgument, "credentials not provided").WithLocalizedMessage(view.ErrStatusTextCredentialsNotProvided)
		}
		info, err = a.interactor.LoginByCredentials(ctx, credentials.GetClientId(), credentials.GetClientSecret())
	default:
		return nil, NewApiError(codes.InvalidArgument, "credentials not provided").WithLocalizedMessage(view.ErrStatusTextCredentialsNotProvided)
	}

	if err != nil {
		switch true {
		case errors.Is(err, sudir.ErrInvalidClient):
			return nil, NewApiError(codes.Unauthenticated, "bad credentials provided").WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
		case errors.Is(err, sudir.ErrAccessDenied):
			return nil, NewApiError(codes.PermissionDenied, "bad credentials provided").WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
		default:
			return nil, NewApiError(codes.Internal, "auth service error", err).WithLocalizedMessage(view.ErrStatusTextServiceError)
		}
	}

	return &authv1.LoginResponse{AccessToken: info.OAuth.AccessToken}, nil
}

func (a *authServer) RefreshToken(ctx context.Context, request *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	id := ""
	sid := ""
	switch requestID := request.GetId().(type) {
	case *authv1.RefreshTokenRequest_CloudId:
		id = requestID.CloudId
		// Для WEB необходимо передавать session id СУДИР для обновления токенов
		sid = request.GetSid()
	case *authv1.RefreshTokenRequest_ClientId:
		id = requestID.ClientId
	}

	if id == "" {
		return nil, NewApiError(codes.InvalidArgument, "id not provided").WithLocalizedMessage(view.ErrStatusTextCloudIDNotProvided)
	}

	token, err := a.interactor.RefreshToken(ctx, id, sid)
	if err != nil {
		switch true {
		case errors.Is(err, sudir.ErrInvalidGrant):
			// refresh token устарел или был уже использован ранее (каким-либо сервисом)
			return nil, NewApiError(codes.Internal, "sudir service error").WithLocalizedMessage(view.ErrStatusTextServiceError)
		case errors.Is(err, repositories.ErrNotFound):
			// refresh token отсутствует
			return nil, NewApiError(codes.NotFound, "refresh token not found").WithLocalizedMessage(view.ErrStatusTextRefreshTokenNotFound)
		default:
			return nil, NewApiError(codes.Internal, "auth service error", err).WithLocalizedMessage(view.ErrStatusTextServiceError)
		}
	}

	return &authv1.RefreshTokenResponse{AccessToken: token}, nil
}

func (a *authServer) Logout(ctx context.Context, request *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	id := ""
	token := ""
	sid := ""

	switch requestID := request.GetId().(type) {
	case *authv1.LogoutRequest_CloudId:
		id = requestID.CloudId
		// Для WEB необходимо передавать session id СУДИР для осуществления logout
		sid = request.GetSid()
	case *authv1.LogoutRequest_Credentials:
		if requestID.Credentials == nil {
			return nil, NewApiError(codes.InvalidArgument, "credentials not provided").WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
		}

		id = requestID.Credentials.GetClientId()
		token = requestID.Credentials.GetRegistrationAccessToken()

		if token == "" {
			return nil, NewApiError(codes.InvalidArgument, "registration token not provided").WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
		}
	}

	if id == "" {
		return nil, NewApiError(codes.InvalidArgument, "id not provided").WithLocalizedMessage(view.ErrStatusTextBadCredentialsProvided)
	}

	if err := a.interactor.Logout(ctx, id, sid, token); err != nil {
		if errors.Is(err, sudir.ErrAccessDenied) {
			return nil, NewApiError(codes.PermissionDenied, "access denied").WithLocalizedMessage(view.ErrStatusTextAccessDenied)
		}
		return nil, NewApiError(codes.Internal, "auth service error").WithLocalizedMessage(view.ErrStatusTextServiceError)
	}

	return &authv1.LogoutResponse{}, nil
}

// GetUser получение информации о пользователе по access_token
func (a *authServer) GetUser(ctx context.Context, request *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	accessToken := request.GetAccessToken()

	if accessToken == "" {
		return nil, NewApiError(codes.InvalidArgument, "access_token not provided").WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
	}

	uInfo, err := a.interactor.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return &authv1.GetUserResponse{
		UserInfo: a.presenter.UserInfoToPb(uInfo),
	}, nil
}

func (a *authServer) GetEmployees(ctx context.Context, request *authv1.GetEmployeesRequest) (*authv1.GetEmployeesResponse, error) {
	var params entity.EmployeeGetParams

	switch request.GetKey().(type) {
	case *authv1.GetEmployeesRequest_CloudId:
		params = entity.EmployeeGetParams{
			Key:     entity.EmployeeGetKeyCloudID,
			CloudID: request.GetCloudId(),
		}
	case *authv1.GetEmployeesRequest_Email:
		params = entity.EmployeeGetParams{
			Key:   entity.EmployeeGetKeyEmail,
			Email: request.GetEmail(),
		}
	default:
		return nil, diterrors.NewApiError(codes.InvalidArgument, "invalid key type").
			WithLocalizedMessage("Передан некорректный ключ")
	}

	employees, err := a.interactor.GetEmployees(ctx, params)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, diterrors.NewApiError(codes.InvalidArgument, "invalid request", err).
				WithLocalizedMessage("Неверные параметры запроса")
		case errors.Is(err, diterrors.ErrNotFound):
			return nil, diterrors.NewApiError(codes.NotFound, "employee not found", err).
				WithLocalizedMessage("Сотрудник не найден")
		default:
			return nil, diterrors.NewApiError(codes.Internal, "auth service error", err).
				WithLocalizedMessage(view.ErrStatusTextServiceError)
		}
	}

	return &authv1.GetEmployeesResponse{
		Employees: a.presenter.EmployeesToPb(employees),
	}, nil
}

// Validate валидация access token в СУДИР.
func (a *authServer) Validate(ctx context.Context, request *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	accessToken := request.GetAccessToken()

	if accessToken == "" {
		return nil, NewApiError(codes.InvalidArgument, "access_token not provided").WithLocalizedMessage(view.ErrStatusTextAccessTokenNotProvided)
	}

	info, err := a.interactor.IsValidToken(ctx, accessToken)
	if err != nil {
		if errors.Is(err, usecase.ErrTokenIsExpire) {
			return nil, NewApiError(codes.PermissionDenied, "token expired").WithLocalizedMessage(view.ErrStatusTextAccessTokenExpired)
		}
		return nil, NewApiError(codes.Internal, "auth service error").WithLocalizedMessage(view.ErrStatusTextServiceError)
	}

	return &authv1.ValidateResponse{
		Info: a.presenter.TokenInfoToPb(info),
	}, nil
}
