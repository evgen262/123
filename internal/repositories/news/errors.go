package news

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
	ErrCategoryAlreadyExists diterrors.StringError = "category already exists"
	ErrCategoryNotFound      diterrors.StringError = "category not found"
)
