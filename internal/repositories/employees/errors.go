package employees

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrEmptyPersonID    diterrors.StringError = "person id is empty"
	ErrEmployeeNotFound diterrors.StringError = "employee not found"
)
