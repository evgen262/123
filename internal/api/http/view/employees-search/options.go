package employees_search

import (
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

type Options struct {
	IsFired bool `json:"isFired"`
	// firedDateRange {from, to} - может передаваться только с isFired
	Range *FiredDateRange `json:"firedDateRange"`
}

type FiredDateRange struct {
	From *timeUtils.Date `json:"from"`
	To   *timeUtils.Date `json:"to"`
}
