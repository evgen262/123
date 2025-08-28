package employees_search

import (
	searchv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/search/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/shared/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	employees_search "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=employees_search

type EmployeesSearchMapper interface {
	GenderToEntity(gender sharedv1.Gender) entity.Gender
	GendersToEntity(genders []sharedv1.Gender) []entity.Gender
	GendersToPb(genders []entity.Gender) []sharedv1.Gender
	GenderToPb(gender entity.Gender) sharedv1.Gender
	SearchAbsencesToPb(absences []*employees_search.SearchAbsence) []*searchv1.SearchRequest_Filters_Absence
	SearchAbsenceToPb(absence *employees_search.SearchAbsence) *searchv1.SearchRequest_Filters_Absence
	AggregationAbsencesToPb(absences []*employees_search.SearchAbsence) []*searchv1.AggregationsRequest_Absence
	AggregationAbsenceToPb(absence *employees_search.SearchAbsence) *searchv1.AggregationsRequest_Absence
}
