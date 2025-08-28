package employees_search

import (
	viewEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/employees-search"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

type employeesSearchPresenter struct{}

func NewEmployeesSearchPresenter() *employeesSearchPresenter {
	return &employeesSearchPresenter{}
}

func (e employeesSearchPresenter) SearchRequestToEntity(search *viewEmployeesSearch.SearchRequest) *entityEmployeesSearch.SearchParams {
	if search == nil {
		return nil
	}

	params := &entityEmployeesSearch.SearchParams{
		Query:   search.Query,
		Filters: e.FiltersToEntity(search.Filters),
		Options: e.OptionsToEntity(search.Options),
		Limit:   search.Limit,
		AfterID: search.AfterID,
	}

	return params
}

func (e employeesSearchPresenter) FiltersToEntity(filter *viewEmployeesSearch.Filters) *entityEmployeesSearch.FiltersParams {
	if filter == nil {
		return nil
	}

	var genders []entity.Gender
	if filter.Gender != "" {
		genders = e.GendersToEntity([]string{filter.Gender})
	}

	return &entityEmployeesSearch.FiltersParams{
		OIVs:          filter.OIVs,
		Organizations: filter.Organizations,
		Products:      filter.Products,
		Subdivisions:  filter.Subdivisions,
		Positions:     filter.Positions,
		Genders:       genders,
		Statuses:      e.StatusToEntity(filter.Status),
	}
}

func (e employeesSearchPresenter) OptionsToEntity(options *viewEmployeesSearch.Options) *entityEmployeesSearch.OptionsParams {
	if options == nil {
		return nil
	}

	o := &entityEmployeesSearch.OptionsParams{}

	if options.IsFired {
		o.WithFired = &entityEmployeesSearch.FiredDateRange{
			From: options.Range.From,
			To:   options.Range.To,
		}
	}

	return o
}

func (e employeesSearchPresenter) FiltersRequestToEntity(req *viewEmployeesSearch.FiltersRequest) *entityEmployeesSearch.SearchParams {
	if req == nil {
		return nil
	}

	return &entityEmployeesSearch.SearchParams{
		Filters: e.FiltersToEntity(req.Filters),
		Options: e.OptionsToEntity(req.Options),
	}
}

func (e employeesSearchPresenter) GenderToEntity(gender string) entity.Gender {
	switch gender {
	case viewEmployeesSearch.GenderMaleView:
		return entity.GenderMale
	case viewEmployeesSearch.GenderFemaleView:
		return entity.GenderFemale
	default:
		return entity.GenderInvalid
	}
}

func (e employeesSearchPresenter) GendersToEntity(genders []string) []entity.Gender {
	if genders == nil {
		return nil
	}
	arr := make([]entity.Gender, 0, len(genders))
	for _, gender := range genders {
		arr = append(arr, e.GenderToEntity(gender))
	}
	return arr
}

func (e employeesSearchPresenter) StatusToEntity(status string) entityEmployeesSearch.FilterStatuses {

	sts := entityEmployeesSearch.FilterStatuses{}

	switch status {
	case viewEmployeesSearch.StatusTypeVacation:
		sts.IsVacation = true
	case viewEmployeesSearch.StatusTypeMaternityLeave:
		sts.IsMaternityLeave = true
	case viewEmployeesSearch.StatusTypeBirthday:
		sts.IsBirthDay = true
	}

	return sts
}

func (e employeesSearchPresenter) SearchResponseToView(search *entityEmployeesSearch.SearchResponse) *viewEmployeesSearch.SearchResponse {
	if search == nil {
		return nil
	}

	return &viewEmployeesSearch.SearchResponse{
		Employees: e.EmployeesToView(search.Employees),
		Total:     search.Total,
		AfterID:   search.AfterID,
	}
}

func (e employeesSearchPresenter) EmployeeToView(employee *entityEmployeesSearch.Employee) *viewEmployeesSearch.Employee {
	if employee == nil {
		return nil
	}

	return &viewEmployeesSearch.Employee{
		ID:           employee.ID.String(),
		FullName:     employee.FullName,
		FirstName:    employee.FirstName,
		MiddleName:   employee.MiddleName,
		LastName:     employee.LastName,
		Gender:       e.GenderToView(employee.Gender),
		ImageID:      employee.ImageID,
		Position:     employee.Position.GetName(),
		OIV:          e.OIVToView(employee.OIV),
		Product:      e.ProductToView(employee.Product),
		Organization: e.OrganizationToView(employee.Organization),
		Structure:    e.StructureToView(employee.Structure),
		Statuses:     e.StatusesToView(employee.Statuses),
	}
}

func (e *employeesSearchPresenter) StatusesToView(statuses *entityEmployeesSearch.Statuses) *viewEmployeesSearch.Statuses {
	if statuses == nil {
		return nil
	}
	return &viewEmployeesSearch.Statuses{
		IsFired:    statuses.IsFired,
		IsBirthday: statuses.IsBirthday,
		Absences:   statuses.Absences,
	}
}

func (e employeesSearchPresenter) EmployeesToView(employees []*entityEmployeesSearch.Employee) []*viewEmployeesSearch.Employee {
	if employees == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.Employee, 0, len(employees))
	for _, employee := range employees {
		emp := e.EmployeeToView(employee)
		if emp != nil {
			arr = append(arr, emp)
		}
	}
	return arr
}

func (e employeesSearchPresenter) OIVToView(oiv *entityEmployeesSearch.OIV) *viewEmployeesSearch.OIV {
	if oiv == nil {
		return nil
	}
	return &viewEmployeesSearch.OIV{
		Name: oiv.Name,
	}
}

func (e employeesSearchPresenter) ProductToView(product *entityEmployeesSearch.Product) *viewEmployeesSearch.Product {
	if product == nil {
		return nil
	}
	return &viewEmployeesSearch.Product{
		Name: product.Name,
	}
}

func (e employeesSearchPresenter) OrganizationToView(organization *entityEmployeesSearch.Organization) *viewEmployeesSearch.Organization {
	if organization == nil {
		return nil
	}
	return &viewEmployeesSearch.Organization{
		Name: organization.Name,
	}
}

func (e employeesSearchPresenter) StructureToView(structure *entityEmployeesSearch.Structure) *viewEmployeesSearch.Structure {
	if structure == nil {
		return nil
	}
	return &viewEmployeesSearch.Structure{
		Position:     structure.Position.GetName(),
		Subdivision:  structure.Subdivision.GetName(),
		Organization: structure.Organization.GetName(),
		OIV:          structure.OIV.GetName(),
	}
}

func (e employeesSearchPresenter) FilterGenderToView(g *entityEmployeesSearch.FilterGender) *viewEmployeesSearch.FilterGender {
	if g == nil {
		return nil
	}
	return &viewEmployeesSearch.FilterGender{
		Name: e.GenderToNameView(g.Gender),
		ID:   e.GenderToView(g.Gender),
	}
}

func (e employeesSearchPresenter) FilterGendersToView(genders []*entityEmployeesSearch.FilterGender) []*viewEmployeesSearch.FilterGender {
	if genders == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterGender, 0, len(genders))
	for _, gender := range genders {
		g := e.FilterGenderToView(gender)
		if g != nil {
			arr = append(arr, g)
		}
	}
	return arr
}

func (e employeesSearchPresenter) FiltersResponseToView(filter *entityEmployeesSearch.FiltersResponse) *viewEmployeesSearch.FiltersResponse {
	if filter == nil {
		return nil
	}

	statuses := make([]string, 0, len(filter.Absences)+1)
	for _, a := range filter.Absences {
		statuses = append(statuses, a.Name)
	}
	if filter.BirthDayCount > 0 {
		statuses = append(statuses, viewEmployeesSearch.StatusTypeBirthday)
	}

	return &viewEmployeesSearch.FiltersResponse{
		Filters: viewEmployeesSearch.FiltersResult{
			OIVs:          e.FilterOIVsToView(filter.OIVs),
			Organizations: e.FilterOrganizationsToView(filter.Organizations),
			Products:      e.FilterProductsToView(filter.Products),
			Subdivisions:  e.FilterSubdivisionsToView(filter.Subdivisions),
			Positions:     e.FilterPositionsToView(filter.Positions),
			Genders:       e.FilterGendersToView(filter.Genders),
			Statuses:      statuses,
		},
	}
}

func (e employeesSearchPresenter) FilterOIVToView(oiv *entityEmployeesSearch.FilterOIV) *viewEmployeesSearch.FilterOIV {
	if oiv == nil {
		return nil
	}

	return &viewEmployeesSearch.FilterOIV{
		ID:   oiv.ID,
		Name: oiv.Name,
	}
}

func (e employeesSearchPresenter) FilterOIVsToView(oivs []*entityEmployeesSearch.FilterOIV) []*viewEmployeesSearch.FilterOIV {
	if oivs == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterOIV, 0, len(oivs))
	for _, oiv := range oivs {
		o := e.FilterOIVToView(oiv)
		if o != nil {
			arr = append(arr, o)
		}
	}
	return arr
}

func (e employeesSearchPresenter) FilterOrganizationToView(organization *entityEmployeesSearch.FilterOrganization) *viewEmployeesSearch.FilterOrganization {
	if organization == nil {
		return nil
	}

	return &viewEmployeesSearch.FilterOrganization{
		ID:   organization.ID,
		Name: organization.Name,
	}
}

func (e employeesSearchPresenter) FilterOrganizationsToView(organizations []*entityEmployeesSearch.FilterOrganization) []*viewEmployeesSearch.FilterOrganization {
	if organizations == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterOrganization, 0, len(organizations))
	for _, organization := range organizations {
		o := e.FilterOrganizationToView(organization)
		if o != nil {
			arr = append(arr, o)
		}
	}
	return arr
}

func (e employeesSearchPresenter) FilterProductToView(product *entityEmployeesSearch.FilterProduct) *viewEmployeesSearch.FilterProduct {
	if product == nil {
		return nil
	}

	return &viewEmployeesSearch.FilterProduct{
		ID:   product.ID,
		Name: product.Name,
	}
}

func (e employeesSearchPresenter) FilterProductsToView(products []*entityEmployeesSearch.FilterProduct) []*viewEmployeesSearch.FilterProduct {
	if products == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterProduct, 0, len(products))
	for _, product := range products {
		p := e.FilterProductToView(product)
		if p != nil {
			arr = append(arr, p)
		}
	}
	return arr
}

func (e employeesSearchPresenter) FilterSubdivisionToView(subdivision *entityEmployeesSearch.FilterSubdivision) *viewEmployeesSearch.FilterSubdivision {
	if subdivision == nil {
		return nil
	}

	return &viewEmployeesSearch.FilterSubdivision{
		ID:   subdivision.ID,
		Name: subdivision.Name,
	}
}

func (e employeesSearchPresenter) FilterSubdivisionsToView(subdivisions []*entityEmployeesSearch.FilterSubdivision) []*viewEmployeesSearch.FilterSubdivision {
	if subdivisions == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterSubdivision, 0, len(subdivisions))
	for _, subdivision := range subdivisions {
		s := e.FilterSubdivisionToView(subdivision)
		if s != nil {
			arr = append(arr, s)
		}
	}
	return arr
}

func (e employeesSearchPresenter) FilterPositionToView(position *entityEmployeesSearch.FilterPosition) *viewEmployeesSearch.FilterPosition {
	if position == nil {
		return nil
	}

	return &viewEmployeesSearch.FilterPosition{
		// ID:   position.ID,
		Name: position.Name,
	}
}

func (e employeesSearchPresenter) FilterPositionsToView(positions []*entityEmployeesSearch.FilterPosition) []*viewEmployeesSearch.FilterPosition {
	if positions == nil {
		return nil
	}
	arr := make([]*viewEmployeesSearch.FilterPosition, 0, len(positions))
	for _, position := range positions {
		p := e.FilterPositionToView(position)
		if p != nil {
			arr = append(arr, p)
		}
	}
	return arr
}

func (e employeesSearchPresenter) GenderToView(gender entity.Gender) string {
	switch gender {
	case entity.GenderMale:
		return viewEmployeesSearch.GenderMaleView
	case entity.GenderFemale:
		return viewEmployeesSearch.GenderFemaleView
	default:
		return viewEmployeesSearch.GenderInvalidView
	}
}

func (e employeesSearchPresenter) GenderToNameView(gender entity.Gender) string {
	switch gender {
	case entity.GenderMale:
		return viewEmployeesSearch.GenderMaleStringName
	case entity.GenderFemale:
		return viewEmployeesSearch.GenderFemaleStringName
	default:
		return viewEmployeesSearch.GenderInvalidStringName
	}
}
