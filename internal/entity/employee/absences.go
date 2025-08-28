package employee

import (
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
)

//go:generate ditgen -source=absences.go

type Absence struct {
	From *timeUtils.Date
	To   *timeUtils.Date
	Name string
	Type AbsenceType
}

type AbsenceType int // Тип причины отсутствия сотрудника

const (
	AbsenceTypeUnknown      AbsenceType = iota
	AbsenceTypeDecree                   // Декрет
	AbsenceTypeDecreeWork               // Работа в декрете
	AbsenceTypeVacation                 // Отпуск
	AbsenceTypeBusinessTrip             // Командировка
	AbsenceTypeMedical                  // Больничный
)

// Возвращает по текущему соответствующе строковое представление типа причины
func (r AbsenceType) String() string {
	a := [...]string{
		"",
		"maternityLeave",
		"", // В данный момент в контракте отсутствует
		"vacation",
		"", // В данный момент в контракте отсутствует
		"sickLeave",
	}

	if int(r) >= len(a) {
		return ""
	}

	return a[r]
}

func (a *Absence) GetPriority() int {
	if a == nil {
		return 999
	}
	switch a.Type {
	case AbsenceTypeDecree:
		return 1
	case AbsenceTypeMedical:
		return 2
	case AbsenceTypeVacation:
		return 3
	default:
		return 999
	}
}
