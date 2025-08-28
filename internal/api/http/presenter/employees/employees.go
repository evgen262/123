package employees

import (
	"strconv"

	"github.com/google/uuid"

	viewEmployees "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/employees"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

var genderToView = map[entity.Gender]string{
	entity.GenderInvalid: viewEmployees.GenderInvalidView,
	entity.GenderMale:    viewEmployees.GenderMaleView,
	entity.GenderFemale:  viewEmployees.GenderFemaleView,
}

type employeesPresenter struct{}

func NewEmployeesPresenter() *employeesPresenter {
	return &employeesPresenter{}
}

func (p *employeesPresenter) GenderToView(gender entity.Gender) string {
	g, ok := genderToView[gender]
	if !ok {
		return viewEmployees.GenderInvalidView
	}
	return g
}

func (p *employeesPresenter) EmployeeToView(employee *entityEmployee.Employee) *viewEmployees.Employee {
	if employee == nil {
		return nil
	}

	var (
		subUnit     *string
		productName *string
	)

	if employee.StaffPosition.GetSubdivision() != nil {
		subUnit = &employee.StaffPosition.Subdivision.Name
	}

	if employee.MainProduct != nil {
		productName = &employee.MainProduct.FullName
	}

	return &viewEmployees.Employee{
		FullName:         employee.FullName,
		Gender:           p.GenderToView(employee.Person.Gender),
		ImageID:          employee.ImageID.String(),
		Birthday:         employee.Person.Birthday.String(),
		Email:            employee.Email,
		DateOfEmployment: employee.DateOfEmployment,
		Position:         employee.StaffPosition.GetPosition().GetName(),
		Statuses: viewEmployees.Statuses{
			IsFired: employee.IsFired,
			Absence: p.AbsenceToView(employee.Absence),
		},
		OIV: viewEmployees.OIV{
			Name:   employee.Portal.Name,
			IconID: employee.Portal.IconID,
		},
		Product: p.productToView(employee.MainProduct),
		LegalEntity: viewEmployees.LegalEntity{
			Name:   employee.Organization.Name,
			IconID: employee.Organization.IconID,
		},
		Structure: viewEmployees.Structure{
			Position:    employee.StaffPosition.GetPosition().GetName(),
			SubUnit:     subUnit,
			LegalEntity: employee.Organization.Name,
			OIV:         employee.Portal.Name,
			Product:     productName,
		},
		WorkPhone:    employee.Phones.GetWorkNumber(),
		AddPhone:     employee.Phones.GetExtensionNumber(),
		MobilePhone:  employee.Person.Phone,
		Workplace:    p.workplaceToView(employee.Workplace),
		HeadOfOrg:    p.HeadStaffPositionToView(employee.HeadStaffPosition),
		HeadOfManage: p.HeadManagementToView(employee.HeadManagement),
	}
}

func (p *employeesPresenter) HeadStaffPositionToView(h *entityEmployee.HeadStaffPosition) *viewEmployees.OrgHead {
	if h == nil {
		return nil
	}

	e := &viewEmployees.OrgHead{
		ID:         h.EmployeeID.String(),
		FirstName:  h.FirstName,
		MiddleName: h.MiddleName,
		LastName:   h.LastName,
		Gender:     p.GenderToView(h.Gender),
	}

	if h.ImageID != uuid.Nil {
		imageID := h.ImageID.String()
		e.ImageID = &imageID
	}

	return e
}

func (p *employeesPresenter) HeadManagementToView(h *entityEmployee.HeadManagement) *viewEmployees.ManageHead {
	if h == nil {
		return nil
	}

	e := &viewEmployees.ManageHead{
		ID:         h.EmployeeID.String(),
		FirstName:  h.FirstName,
		MiddleName: h.MiddleName,
		LastName:   h.LastName,
		Gender:     p.GenderToView(h.Gender),
	}

	if h.ImageID != uuid.Nil {
		imageID := h.ImageID.String()
		e.ImageID = &imageID
	}

	return e
}

func (p employeesPresenter) AbsenceToView(a *entityEmployee.Absence) *viewEmployees.Absence {
	if a == nil {
		return nil
	}

	return &viewEmployees.Absence{
		From: a.From,
		To:   a.To,
		Type: p.AbsenceTypeToView(a.Type),
	}
}

func (p *employeesPresenter) productToView(product *entityEmployee.Product) *viewEmployees.Product {
	if product == nil {
		return nil
	}

	return &viewEmployees.Product{
		Name:   product.FullName,
		IconID: product.IconID.String(),
	}
}

func (p *employeesPresenter) workplaceToView(workplace *entityEmployee.Workplace) *viewEmployees.Workplace {
	if workplace == nil {
		return nil
	}

	f := strconv.Itoa(workplace.Floor)

	return &viewEmployees.Workplace{
		Address: &workplace.Address,
		Floor:   &f,
		Cabinet: &workplace.CabinetNumber,
	}
}

func (p *employeesPresenter) AbsenceTypeToView(a entityEmployee.AbsenceType) string {
	switch a {
	case entityEmployee.AbsenceTypeDecree:
		return a.String()
	case entityEmployee.AbsenceTypeMedical:
		return a.String()
	case entityEmployee.AbsenceTypeVacation:
		return a.String()
	default:
		return "unknown"
	}
}
