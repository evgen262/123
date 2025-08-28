package sudir

import "time"

type OAuthResponse struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	Expiry       *time.Time
}

type AuthURLOptions struct {
	State               string
	IsOffline           bool
	RedirectURI         string
	ClientID            string
	CodeChallengeMethod string
	CodeChallenge       string
}

type CodeExchangeOptions struct {
	RedirectURI  string
	ClientID     string
	ClientSecret string
	CodeVerifier string
}

type LoginOptions struct {
	ClientID     string
	ClientSecret string
}
