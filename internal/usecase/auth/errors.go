package auth

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrPortalsNotFound   diterrors.StringError = "portals not found"
	ErrEmployeesNotFound diterrors.StringError = "employees not found"
	ErrSUDIRNoCloudID    diterrors.StringError = "user does not have cloud_id"
	ErrInvalidSession    diterrors.StringError = "session is invalid"
	ErrInvalidDevice     diterrors.StringError = "device is invalid"
	ErrUnavailablePortal diterrors.StringError = "selected portal unavailable for user"
	ErrEmptyPortalURL    diterrors.StringError = "selected portal url is empty"
	ErrUserInfoRequired  diterrors.StringError = "user info required"
	ErrUserAccessDenied  diterrors.StringError = "user access denied"
)
