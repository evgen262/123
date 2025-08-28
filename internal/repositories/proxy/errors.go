package proxy

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrUnauthorized     diterrors.StringError = "unauthorized"
	ErrBadRequest       diterrors.StringError = "bad request"
	ErrInternal         diterrors.StringError = "internal error"
	ErrPermissionDenied diterrors.StringError = "permission denied"
	ErrTransportRequest diterrors.StringError = "transport request error"
	ErrNotFound         diterrors.StringError = "not found"
)
