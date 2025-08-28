package users

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrEmptySessionPersonID      diterrors.StringError = "person id in session is empty"
	ErrEmptySessionEmployeeExtID diterrors.StringError = "employee ext id in session is empty"
	ErrZeroSessionActivePortalID diterrors.StringError = "active portal id in session is zero"
)
