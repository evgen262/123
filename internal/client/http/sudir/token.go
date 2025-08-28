package sudir

import "time"

type Token struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}
