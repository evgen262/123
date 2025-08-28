package employees_search

import (
	searchv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/search/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/shared/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	employees_search "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

type employeesSearchMapper struct {
}

func NewEmployeesSearchMapper() *employeesSearchMapper {
	return &employeesSearchMapper{}
}

func (e employeesSearchMapper) GenderToEntity(gender sharedv1.Gender) entity.Gender {
	switch gender {
	case sharedv1.Gender_GENDER_MALE:
		return entity.GenderMale
	case sharedv1.Gender_GENDER_FEMALE:
		return entity.GenderFemale
	default:
		return entity.GenderInvalid
	}
}

func (e employeesSearchMapper) GenderToPb(gender entity.Gender) sharedv1.Gender {
	switch gender {
	case entity.GenderMale:
		return sharedv1.Gender_GENDER_MALE
	case entity.GenderFemale:
		return sharedv1.Gender_GENDER_FEMALE
	default:
		return sharedv1.Gender_GENDER_INVALID
	}
}

func (e employeesSearchMapper) GendersToEntity(genders []sharedv1.Gender) []entity.Gender {
	if genders == nil {
		return nil
	}
	arr := make([]entity.Gender, 0, len(genders))
	for _, g := range genders {
		arr = append(arr, e.GenderToEntity(g))
	}
	return arr
}

func (e employeesSearchMapper) GendersToPb(genders []entity.Gender) []sharedv1.Gender {
	if genders == nil {
		return nil
	}
	arr := make([]sharedv1.Gender, 0, len(genders))
	for _, g := range genders {
		arr = append(arr, e.GenderToPb(g))
	}
	return arr
}

func (e employeesSearchMapper) SearchAbsencesToPb(absences []*employees_search.SearchAbsence) []*searchv1.SearchRequest_Filters_Absence {
	if absences == nil {
		return nil
	}
	absPb := make([]*searchv1.SearchRequest_Filters_Absence, 0, len(absences))
	for _, a := range absences {
		if aPb := e.SearchAbsenceToPb(a); aPb != nil {
			absPb = append(absPb, aPb)
		}
	}
	return absPb
}
func (e employeesSearchMapper) SearchAbsenceToPb(absence *employees_search.SearchAbsence) *searchv1.SearchRequest_Filters_Absence {
	if absence == nil {
		return nil
	}

	return &searchv1.SearchRequest_Filters_Absence{
		Name: absence.Name,
		From: absence.From.String(),
		To:   absence.To.String(),
	}
}

func (e employeesSearchMapper) AggregationAbsencesToPb(absences []*employees_search.SearchAbsence) []*searchv1.AggregationsRequest_Absence {
	if absences == nil {
		return nil
	}
	absPb := make([]*searchv1.AggregationsRequest_Absence, 0, len(absences))
	for _, a := range absences {
		if aPb := e.AggregationAbsenceToPb(a); aPb != nil {
			absPb = append(absPb, aPb)
		}
	}
	return absPb
}
func (e employeesSearchMapper) AggregationAbsenceToPb(absence *employees_search.SearchAbsence) *searchv1.AggregationsRequest_Absence {
	if absence == nil {
		return nil
	}

	return &searchv1.AggregationsRequest_Absence{
		Name: absence.Name,
		From: absence.From.String(),
		To:   absence.To.String(),
	}
}
