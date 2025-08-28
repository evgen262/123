package employees_search

import (
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

//go:generate ditgen -source=absence.go -nil=true -zero=true -all-fields=true

type Absence struct {
	Name      string
	StartDate *timeUtils.Date
	EndDate   *timeUtils.Date
}
