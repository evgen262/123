package auth

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrUserInfoRequired  = diterrors.StringError("session required")
)
