package auth

import (
	"context"
	"fmt"
	"net/url"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

type redirectSessionInteractor struct {
	repository RedirectSessionRepository

	authLink string
}

func NewRedirectSessionInteractor(repository RedirectSessionRepository, authLink string) *redirectSessionInteractor {
	return &redirectSessionInteractor{
		repository: repository,

		authLink: authLink,
	}
}

func (rhs redirectSessionInteractor) CreateSession(ctx context.Context, userInfo *entitySession.RedirectSessionUserInfo) (string, error) {
	if userInfo == nil {
		return "", fmt.Errorf("redirectSessionInteractor.Create: %w", ErrUserInfoRequired)
	}

	sessionID, err := rhs.repository.CreateSession(ctx, userInfo)
	if err != nil {
		return "", fmt.Errorf("repository.CreateSession: %w", err)
	}

	redirectURL := url.URL{
		Scheme: "https",
		Host:   userInfo.PortalURL,
		Path:   rhs.authLink,
	}

	urlValues := redirectURL.Query()
	urlValues.Set("state", sessionID)
	redirectURL.RawQuery = urlValues.Encode()

	return redirectURL.String(), nil
}
