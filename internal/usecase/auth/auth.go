package auth

import (
	"context"
	"errors"
	"fmt"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
)

type authUseCase struct {
	repository Repository
	logger     ditzap.Logger

	accessList entityAuth.AccessList
}

func NewAuthUseCase(repository Repository, logger ditzap.Logger, accessList entityAuth.AccessList) *authUseCase {
	return &authUseCase{
		repository: repository,
		logger:     logger,
		accessList: accessList,
	}
}

// GetAuthURL
//
//	URL для перенаправления пользователя для авторизации в СУДИР
func (a *authUseCase) GetAuthURL(ctx context.Context, callbackURI string) (string, error) {
	redirectURL, err := a.repository.GetRedirectURL(ctx, callbackURI)
	if err != nil {
		if errors.As(err, new(diterrors.LocalizedError)) {
			a.logger.Error("cant get auth url", zap.Error(errors.Unwrap(err)))
		} else {
			a.logger.Error("cant get auth url", zap.Error(err))
		}
		return "", fmt.Errorf("cant get auth url: %w", err)
	}
	return redirectURL, nil
}

func (a *authUseCase) Auth(ctx context.Context, code, state, callbackURI string) (*entityAuth.Auth, error) {
	info, err := a.repository.Auth(ctx, code, state, callbackURI)
	if err != nil {
		if errors.Is(err, diterrors.ErrNotFound) {
			return nil, ErrEmployeesNotFound
		}
		if !errors.As(err, new(diterrors.ValidationError)) && errors.As(err, new(diterrors.LocalizedError)) {
			a.logger.Error("cant authenticate user",
				zap.String("code", code),
				zap.String("state", state),
				zap.Error(err),
			)
		} else {
			a.logger.Debug("cant authenticate user",
				zap.String("code", code),
				zap.String("state", state),
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("cant authenticate user: %w", err)
	}

	if info.GetUser() == nil {
		return nil, diterrors.ErrPermissionDenied
	}

	if a.accessList != nil && len(a.accessList) > 0 {
		if !a.accessList.Have(info.GetUser().Email) {
			return nil, ErrUserAccessDenied
		}
	}

	// TODO: Отдать этот корнер-кейс на пересмотр аналитику
	/*
		if info.GetUser().CloudID == "" {
			return nil, ErrSUDIRNoCloudID
		}
	*/

	if len(info.GetUser().Portals) == 0 {
		a.logger.Error("нет порталов связанных с сотрудником",
			entityAuth.LogModuleUE,
			entityAuth.LogCode("UE_030"),
			zap.String("login", info.GetUser().Login),
			zap.String("email", info.GetUser().Email),
			zap.String("user", info.GetUser().FIO),
			zap.Error(err),
		)
		return nil, ErrPortalsNotFound
	}
	auth1C, err := a.repository.AuthPortal(ctx, entityAuth.AuthPortalParams{
		// TODO: временно аутентифицируем на первом портале
		//  в дальнейшем брать основное место работы по умолчанию
		//  если несколько мест работы
		PortalURL: info.GetUser().Portals[0].URL,
		User: entityAuth.User1C{
			CloudID: info.GetUser().CloudID,
			SNILS:   info.GetUser().SNILS,
			Email:   info.GetUser().Email,
		},
	})
	if err != nil {
		// TODO: сделать рефакторинг.
		if !errors.As(err, new(diterrors.ValidationError)) && errors.As(err, new(diterrors.LocalizedError)) {
			a.logger.Error("cant authenticate user into portal",
				zap.String("user_id", info.GetUser().CloudID),
				zap.String("portal_url", info.GetUser().Portals[0].URL),
				zap.Error(err),
			)
		} else if errors.Is(err, diterrors.ErrNotFound) {
			return nil, err //nolint:wrapcheck
		} else {
			a.logger.Debug("cant authenticate user into portal",
				zap.String("user_id", info.GetUser().CloudID),
				zap.String("portal_url", info.GetUser().Portals[0].URL),
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("cant authenticate user: %w", err)
	}

	// ошибка заглушена, потому что IP может быть не передан в заголовке
	clientIP, _ := entity.ClientIPFromContext(ctx)
	device, err := entity.DeviceFromContext(ctx)
	if err != nil || device.UserAgent == "" {
		return nil, ErrInvalidDevice
	}

	tokensPair, err := a.repository.CreateSession(ctx, info.GetUser(), clientIP, device, auth1C)
	if err != nil {
		if errors.Is(err, repositories.ErrNoUserPortals) {
			return nil, ErrPortalsNotFound
		}
		a.logger.Error("can't create session",
			zap.String("user_id", info.GetUser().CloudID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("can't create session: %w", err)
	}

	return &entityAuth.Auth{
		JWTToken:      tokensPair.AccessToken,
		RefreshToken:  tokensPair.RefreshToken,
		PortalSession: auth1C.PortalSession,
		Portals:       info.GetUser().Portals,
	}, nil
}

func (a *authUseCase) GetSession(ctx context.Context, accessToken string) (*entityAuth.Session, error) {
	session, err := a.repository.GetSession(ctx, accessToken)
	if err != nil {
		switch {
		case errors.As(err, new(diterrors.ValidationError)):
			return nil, err
		case errors.Is(err, diterrors.ErrFailedPrecondition):
			return nil, err
		case errors.Is(err, diterrors.ErrNotFound):
			return nil, err
		case errors.Is(err, diterrors.ErrUnauthenticated):
			return nil, err
		default:
			a.logger.Warn("can't get session in repository",
				ditzap.JWTField("access_token", accessToken),
				zap.Error(err),
			)
			return nil, fmt.Errorf("can't get session in repository: %w", err)
		}
	}
	return session, nil
}

func (a *authUseCase) Logout(ctx context.Context, accessToken, refreshToken string) error {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		a.logger.Error("can't get session from context", zap.Error(err))
		return usecase.ErrGetSessionFromContext
	}

	err = a.repository.Logout(ctx, session, accessToken, refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrNilSession) || errors.Is(err, repositories.ErrEmptyAccessToken):
			a.logger.Warn("invalid session",
				zap.String("session_id", session.GetID().String()),
				zap.String("login", session.GetUser().GetLogin()),
				zap.String("portal_url", session.GetActivePortal().GetPortal().URL),
				zap.Error(err),
			)
			return ErrInvalidSession
		case errors.As(err, new(diterrors.ValidationError)):
			fallthrough
		case errors.Is(err, diterrors.ErrNotFound):
			a.logger.Warn("cant Logout",
				zap.String("session_id", session.GetID().String()),
				zap.String("login", session.GetUser().GetLogin()),
				zap.String("portal_url", session.GetActivePortal().GetPortal().URL),
				zap.Error(err),
			)
			return err
		default:
			a.logger.Warn("can't logout in repository",
				zap.String("session_id", session.GetID().String()),
				zap.String("login", session.GetUser().GetLogin()),
				zap.String("portal_url", session.GetActivePortal().GetPortal().URL),
				zap.Error(err),
			)
			return fmt.Errorf("can't logout in repository: %w", err)
		}
	}

	return nil
}

func (a *authUseCase) ChangePortal(ctx context.Context, selectedPortalID int) ([]*entityAuth.Portal, string, error) {
	session, err := entity.SessionFromContext(ctx)
	if err != nil {
		a.logger.Error("can't get session from context", zap.Error(err))
		return nil, "", usecase.ErrGetSessionFromContext
	}

	portals, portalID, err := a.repository.ChangePortal(ctx, selectedPortalID, session)
	if err != nil {
		return nil, "", fmt.Errorf("can't change portal in repository: %w", err)
	}

	if len(portals) == 0 {
		return nil, "", ErrPortalsNotFound
	}

	var selectedPortal *entityAuth.Portal
	for _, portal := range portals {
		if portal.ID == selectedPortalID {
			if portal.URL == "" {
				return nil, "", ErrEmptyPortalURL
			}
			selectedPortal = portal
			selectedPortal.IsSelected = true
			break
		}
	}

	if selectedPortal == nil {
		return nil, "", ErrUnavailablePortal
	}

	return portals, portalID, nil
}

func (a *authUseCase) RefreshTokensPair(ctx context.Context, accessToken, refreshToken string) (*entityAuth.TokensPair, error) {
	tokensPair, err := a.repository.RefreshTokensPair(ctx, accessToken, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("can't change portal in repository: %w", err)
	}

	return tokensPair, nil
}
