package employees_search

import (
	"context"
	"fmt"
	"strings"

	searchv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/search/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employeessearch/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"

	entityEmployeesSearch "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employees-search"
)

type employeesSearchRepository struct {
	searchAPI searchv1.SearchAPIClient
	mapper    EmployeesSearchMapper
	tu        timeUtils.TimeUtils
	logger    ditzap.Logger
}

func NewEmployeesSearchRepository(
	searchAPI searchv1.SearchAPIClient,
	mapper EmployeesSearchMapper,
	tu timeUtils.TimeUtils,
	logger ditzap.Logger,
) *employeesSearchRepository {
	return &employeesSearchRepository{
		searchAPI: searchAPI,
		mapper:    mapper,
		tu:        tu,
		logger:    logger,
	}
}

func (e employeesSearchRepository) Search(ctx context.Context, request *entityEmployeesSearch.SearchRequest) (*entityEmployeesSearch.SearchResponse, error) {
	if request == nil {
		return nil, diterrors.NewValidationError(ErrIsEmpty)
	}

	req := &searchv1.SearchRequest{
		Query: request.Query,
		Size:  int64(request.Limit),
	}

	if request.GetFilters() != nil {
		req.Filters = &searchv1.SearchRequest_Filters{
			Positions: request.GetFilters().Positions,
			Genders:   e.mapper.GendersToPb(request.GetFilters().Genders),
			Absences:  e.mapper.SearchAbsencesToPb(request.GetFilters().Absences),
		}

		req.Filters.OivIds = make([]int64, 0, len(request.GetFilters().OIVs))
		for _, oiv := range request.GetFilters().OIVs {
			req.Filters.OivIds = append(req.Filters.OivIds, int64(oiv))
		}

		req.Filters.OrganizationsIds = make([]string, 0, len(request.GetFilters().Organizations))
		for _, organization := range request.GetFilters().Organizations {
			req.Filters.OrganizationsIds = append(req.Filters.OrganizationsIds, organization.String())
		}

		req.Filters.ProductsIds = make([]string, 0, len(request.GetFilters().Products))
		for _, product := range request.GetFilters().Products {
			req.Filters.ProductsIds = append(req.Filters.ProductsIds, product.String())
		}

		req.Filters.SubdivisionsIds = make([]string, 0, len(request.GetFilters().Subdivisions))
		for _, subdivision := range request.GetFilters().Subdivisions {
			req.Filters.SubdivisionsIds = append(req.Filters.SubdivisionsIds, subdivision.String())
		}

		if request.GetFilters().GetBirthDay() != nil {
			req.Filters.DayOfBirthdate = *request.GetFilters().GetBirthDay()
		}
	}

	if request.GetOptions() != nil {
		req.Options = &searchv1.SearchRequest_Options{}

		if request.GetOptions().WithFired != nil {
			req.Options.Fired = &searchv1.SearchRequest_Options_Fired{
				From: request.GetOptions().WithFired.From.String(),
				To:   request.GetOptions().WithFired.To.String(),
			}
		}
	}

	if request.GetAfterID() != nil {
		req.Pagination = &sharedv1.Pagination{
			AfterId: *request.GetAfterID(),
		}
	}

	res, err := e.searchAPI.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't search employees in employees-search client: %w", err)
	}

	employees := make([]*entityEmployeesSearch.Employee, 0, len(res.GetPersons()))
	now := e.tu.New().Format("01-02")

	for _, person := range res.GetPersons() {
		var middleName *string
		fullNameArr := make([]string, 0, 3)
		fullNameArr = append(fullNameArr, person.GetLastName(), person.GetFirstName())
		if person.GetMiddleName() != "" {
			mn := person.GetMiddleName()
			middleName = &mn
			fullNameArr = append(fullNameArr, person.GetMiddleName())
		}
		fullName := strings.Join(fullNameArr, " ")

		isBirthday := false
		if person.GetBirthDay() != "" {
			isBirthday = now == person.GetBirthDay()
		}

		for _, employee := range person.GetEmployees() {
			id, _ := uuid.Parse(employee.GetId())
			var (
				pos         *entityEmployeesSearch.Position
				oiv         *entityEmployeesSearch.OIV
				product     *entityEmployeesSearch.Product
				org         *entityEmployeesSearch.Organization
				subdivision *entityEmployeesSearch.Subdivision
				structure   *entityEmployeesSearch.Structure
			)
			if employee.GetPosition() != nil {
				pos = &entityEmployeesSearch.Position{
					Name: employee.GetPosition().GetName(),
				}
			}

			if employee.GetOiv() != nil {
				oiv = &entityEmployeesSearch.OIV{
					Name: employee.GetOiv().GetName(),
				}
			}
			if employee.GetProduct() != nil {
				product = &entityEmployeesSearch.Product{
					Name: employee.GetProduct().GetName(),
				}
			}
			if employee.GetOrganization() != nil {
				org = &entityEmployeesSearch.Organization{
					Name: employee.GetOrganization().GetName(),
				}
			}

			if employee.GetStructure().GetSubdivision() != nil {
				subdivision = &entityEmployeesSearch.Subdivision{
					Name: employee.GetStructure().GetSubdivision().GetName(),
				}
			}

			if pos != nil || oiv != nil || product != nil || org != nil || subdivision != nil {
				structure = &entityEmployeesSearch.Structure{
					Position:     pos,
					Subdivision:  subdivision,
					Organization: org,
					OIV:          oiv,
				}
			}

			absences := make([]*entityEmployeesSearch.Absence, 0, len(employee.GetAbsences()))
			for _, absence := range employee.GetAbsences() {
				if absence == nil {
					continue
				}
				absences = append(absences, &entityEmployeesSearch.Absence{
					Name:      absence.GetName(),
					StartDate: (*timeUtils.Date)(e.tu.TimestampToTime(absence.GetStartTime())),
					EndDate:   (*timeUtils.Date)(e.tu.TimestampToTime(absence.GetEndTime())),
				})
			}

			employees = append(employees, &entityEmployeesSearch.Employee{
				ID:           id,
				FullName:     fullName,
				FirstName:    person.GetFirstName(),
				MiddleName:   middleName,
				LastName:     person.GetLastName(),
				Gender:       e.mapper.GenderToEntity(person.GetGender()),
				ImageID:      person.GetPhoto(),
				Position:     pos,
				OIV:          oiv,
				Product:      product,
				Organization: org,
				Structure:    structure,
				Statuses: &entityEmployeesSearch.Statuses{
					IsFired:    employee.GetIsFired(),
					IsBirthday: isBirthday,
					Absences:   absences,
				},
			})
		}
	}

	total := res.GetPagination().GetTotal()

	return &entityEmployeesSearch.SearchResponse{
		Employees: employees,
		Total:     int(total),
		AfterID:   res.GetPagination().GetAfterId(),
	}, nil
}

func (e employeesSearchRepository) Filters(ctx context.Context, request *entityEmployeesSearch.FiltersRequest) (*entityEmployeesSearch.FiltersResponse, error) {
	if request == nil {
		return nil, diterrors.NewValidationError(ErrIsEmpty)
	}

	req := &searchv1.AggregationsRequest{}

	if request.GetFilters() != nil {
		req.OivIds = make([]int64, 0, len(request.GetFilters().OIVs))
		for _, oiv := range request.GetFilters().OIVs {
			req.OivIds = append(req.OivIds, int64(oiv))
		}

		req.OrganizationsIds = make([]string, 0, len(request.GetFilters().Organizations))
		for _, organization := range request.GetFilters().Organizations {
			req.OrganizationsIds = append(req.OrganizationsIds, organization.String())
		}

		req.ProductsIds = make([]string, 0, len(request.GetFilters().Products))
		for _, product := range request.GetFilters().Products {
			req.ProductsIds = append(req.ProductsIds, product.String())
		}

		req.SubdivisionsIds = make([]string, 0, len(request.GetFilters().Subdivisions))
		for _, subdivision := range request.GetFilters().Subdivisions {
			req.SubdivisionsIds = append(req.SubdivisionsIds, subdivision.String())
		}

		req.Positions = request.GetFilters().Positions
		req.Genders = e.mapper.GendersToPb(request.GetFilters().Genders)

		req.Absences = e.mapper.AggregationAbsencesToPb(request.GetFilters().Absences)
		// TODO: добавтиь в фильтр условие по дню рождения
	}

	req.Parameters = &searchv1.AggregationsRequest_Parameters{
		Absences: &searchv1.AggregationsRequest_Parameters_Absence{
			Names: request.Params.Names,
			From:  request.Params.From.String(),
			To:    request.Params.To.String(),
		},
	}

	if request.GetOptions() != nil {
		req.Options = &searchv1.AggregationsRequest_Options{}
		if request.Options.WithFired != nil {
			req.Options.Fired = &searchv1.AggregationsRequest_Options_Fired{
				From: request.Options.WithFired.From.String(),
				To:   request.Options.WithFired.To.String(),
			}
		}
	}

	res, err := e.searchAPI.Aggregations(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't gen aggregations from employees-search client: %w", err)
	}

	respOIVs := make([]*entityEmployeesSearch.FilterOIV, 0, len(res.GetAggregations().GetOivs()))
	for _, oiv := range res.GetAggregations().GetOivs() {
		id := oiv.GetId()
		respOIVs = append(respOIVs, &entityEmployeesSearch.FilterOIV{
			ID:   int(id),
			Name: oiv.GetName(),
		})
	}

	respOrgs := make([]*entityEmployeesSearch.FilterOrganization, 0, len(res.GetAggregations().GetOrganizations()))
	for _, organization := range res.GetAggregations().GetOrganizations() {
		id, _ := uuid.Parse(organization.GetId())
		respOrgs = append(respOrgs, &entityEmployeesSearch.FilterOrganization{
			ID:   id,
			Name: organization.GetName(),
		})
	}

	respProducts := make([]*entityEmployeesSearch.FilterProduct, 0, len(res.GetAggregations().GetProducts()))
	for _, product := range res.GetAggregations().GetProducts() {
		id, _ := uuid.Parse(product.GetId())
		respProducts = append(respProducts, &entityEmployeesSearch.FilterProduct{
			ID:   id,
			Name: product.GetName(),
		})
	}

	respSubdivisions := make([]*entityEmployeesSearch.FilterSubdivision, 0, len(res.GetAggregations().GetSubdivisions()))
	for _, subdivision := range res.GetAggregations().GetSubdivisions() {
		id, _ := uuid.Parse(subdivision.GetId())
		respSubdivisions = append(respSubdivisions, &entityEmployeesSearch.FilterSubdivision{
			ID:   id,
			Name: subdivision.GetName(),
		})
	}

	respPositions := make([]*entityEmployeesSearch.FilterPosition, 0, len(res.GetAggregations().GetPositions()))
	for _, position := range res.GetAggregations().GetPositions() {
		id, _ := uuid.Parse(position.GetId())
		respPositions = append(respPositions, &entityEmployeesSearch.FilterPosition{
			ID:   id,
			Name: position.GetName(),
		})
	}

	respGenders := make([]*entityEmployeesSearch.FilterGender, 0, len(res.GetAggregations().GetGenders()))
	for _, gender := range res.GetAggregations().GetGenders() {
		respGenders = append(respGenders, &entityEmployeesSearch.FilterGender{
			Gender:     e.mapper.GenderToEntity(gender.GetGender()),
			IsDisabled: gender.GetEmployeesCount() == 0,
		})
	}

	respAbsences := make([]*entityEmployeesSearch.FilterAbsence, 0, len(res.GetAggregations().GetAbsences()))
	for _, absence := range res.GetAggregations().GetAbsences() {
		respAbsences = append(respAbsences, entityEmployeesSearch.MakeFilterAbsence(absence.GetName(), int(absence.GetEmployeesCount())))
	}

	return &entityEmployeesSearch.FiltersResponse{
		OIVs:          respOIVs,
		Organizations: respOrgs,
		Products:      respProducts,
		Subdivisions:  respSubdivisions,
		Positions:     respPositions,
		Genders:       respGenders,
		BirthDayCount: 0, // TODO: добавить кол-во сотрудников у которых сегодня день рождения по реализации в сервисе employees-search
		Absences:      respAbsences,
	}, nil
}
