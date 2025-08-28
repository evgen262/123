package files

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrSessionIsEmpty = diterrors.StringError("session is empty")
	ErrFileIdIsEmpty  = diterrors.StringError("file id is empty")
	ErrVisitorIsEmpty = diterrors.StringError("visitor is empty")
	ErrFileIsEmpty    = diterrors.StringError("file is empty")
)
