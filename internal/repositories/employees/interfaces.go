package employees

import (
	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=employees

type MapperEmployees interface {
	EmployeeToEntity(employee *employeev1.CompositeEmployee) *entityEmployee.Employee
	GenderToEntity(gender employeev1.GenderType) entity.Gender
}
