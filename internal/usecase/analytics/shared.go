package analytics

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrBodyIsEmpty     = diterrors.StringError("body is empty")
	ErrUnauthenticated = diterrors.StringError("unauthenticated")
	ErrHeadersAreAmpty = diterrors.StringError("headers are empty")
)
