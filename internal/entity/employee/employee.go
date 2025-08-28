package employee

import (
	"strings"
	"time"

	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
)

//go:generate ditgen -source=employee.go

// EmployeeID идентификатор сотрудника
type EmployeeID uuid.UUID

type OperationType string

const (
	OperationTypeHiring            OperationType = "прием на работу"
	OperationTypeTransfer          OperationType = "перемещение"
	OperationTypeDismissal         OperationType = "увольнение"
	OperationTypeWorkingConditions OperationType = "изменение условий труда"
	OperationTypeNameChange        OperationType = "изменение фио" // ФИО -> фио для нижнего регистра
	OperationTypeMaternityLeave    OperationType = "отпуск по беременности и родам"
	OperationTypeChildcareLeave    OperationType = "отпуск по уходу за ребенком"
	OperationTypeInvalid           OperationType = "невалидный тип операции"
)

func (o *OperationType) String() string {
	if o == nil {
		return string(OperationTypeInvalid)
	}
	return string(*o)
}

func OperationTypeFromString(s string) OperationType {
	switch OperationType(strings.TrimSpace(strings.ToLower(s))) {
	case OperationTypeHiring:
		return OperationTypeHiring
	case OperationTypeTransfer:
		return OperationTypeTransfer
	case OperationTypeDismissal:
		return OperationTypeDismissal
	case OperationTypeWorkingConditions:
		return OperationTypeWorkingConditions
	case OperationTypeNameChange:
		return OperationTypeNameChange
	case OperationTypeMaternityLeave:
		return OperationTypeMaternityLeave
	case OperationTypeChildcareLeave:
		return OperationTypeChildcareLeave
	default:
		return OperationTypeInvalid
	}
}

func (id *EmployeeID) String() string {
	return id.UUID().String()
}

func (id *EmployeeID) UUID() uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return (uuid.UUID)(*id)
}

func (id *EmployeeID) Parse(s string) *EmployeeID {
	if id == nil {
		return nil
	}

	_id, err := uuid.Parse(s)
	if err != nil {
		_id = uuid.Nil
	}
	*id = EmployeeID(_id)

	return id
}

type Employee struct {
	ID                uuid.UUID
	Portal            Portal
	Person            Person
	Number            string
	FullName          string
	EmploymentType    string
	Phones            *Phones
	Email             string
	Rate              float32
	ImageID           uuid.UUID
	MainProduct       *Product
	Products          []*Product
	Organization      Organization
	StaffPosition     *StaffPosition
	Managements       []*Management
	SubdivisionTree   *SubdivisionTree
	ManagementTree    *ManagementTree
	HeadManagement    *HeadManagement
	HeadStaffPosition *HeadStaffPosition
	Absence           *Absence
	Workplace         *Workplace
	DateOfEmployment  *timeUtils.Date
	History           []*History
	IsFired           bool
	CreateTime        *time.Time
	UpdateTime        *time.Time
}

type Phones struct {
	WorkNumber      *string
	ExtensionNumber *string
}

func (p *Phones) GetWorkNumberStr() *string {
	if p == nil {
		return nil
	}
	return p.WorkNumber
}

func (p *Phones) GetExtensionNumberStr() *string {
	if p == nil {
		return nil
	}
	return p.ExtensionNumber
}

type History struct {
	ID             uuid.UUID
	EventType      OperationType
	EventName      string
	EmployeeNumber string
	EmploymentType string
	EventTime      *time.Time
	Rate           float32
	IsReturn       bool
}

type HeadStaffPosition struct {
	PersonID      uuid.UUID
	ExtPersonID   string
	EmployeeID    uuid.UUID
	ExtEmployeeID string
	PositionID    uuid.UUID
	ExtPositionID string
	PositionName  string
	ImageID       uuid.UUID
	FirstName     string
	LastName      string
	MiddleName    string
	Gender        entity.Gender
}

type HeadManagement struct {
	PersonID      uuid.UUID
	ExtPersonID   string
	EmployeeID    uuid.UUID
	ExtEmployeeID string
	RoleID        uuid.UUID
	ExtRoleID     string
	RoleName      string
	ImageID       uuid.UUID
	FirstName     string
	LastName      string
	MiddleName    string
	Gender        entity.Gender
}
