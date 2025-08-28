package sudir

import (
	"context"

	"golang.org/x/oauth2"
)

//go:generate mockgen -source=interfaces.go -destination=./oauth2_mock.go -package=sudir

type OAuth2 interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}

type TokenSource interface {
	Token() (*oauth2.Token, error)
}
