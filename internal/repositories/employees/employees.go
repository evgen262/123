package employees

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"go.uber.org/zap"

	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

const TWO_WEEKS = 14 * 24 * time.Hour

type employeesRepository struct {
	employeesAPI           employeev1.EmployeesAPIClient
	employeesServiceMapper MapperEmployees
	tu                     timeUtils.TimeUtils
	logger                 ditzap.Logger
}

func NewEmployeesRepository(
	employeesAPI employeev1.EmployeesAPIClient,
	employeesServiceMapper MapperEmployees,
	tu timeUtils.TimeUtils,
	logger ditzap.Logger,
) *employeesRepository {
	return &employeesRepository{
		employeesAPI:           employeesAPI,
		employeesServiceMapper: employeesServiceMapper,
		tu:                     tu,
		logger:                 logger,
	}
}

func (er *employeesRepository) Get(ctx context.Context, id uuid.UUID) (*entityEmployee.Employee, error) {
	req := employeev1.CompositeGetRequest{
		Key: &employeev1.CompositeGetRequest_Id{
			Id: id.String(),
		},
		WithPerson:          true,
		WithOrganization:    true,
		WithProducts:        true,
		WithSubdivision:     true,
		WithStaffPosition:   true,
		WithManagement:      true,
		WithPosition:        true,
		WithWorkplace:       true,
		WithSubdivisionTree: true,
		WithManagementsTree: true,
		WithEmployeeHistory: true,
		WithAbsences:        true,
	}

	res, err := er.employeesAPI.CompositeGet(ctx, &req)
	if err != nil {
		gErr := diterrors.GrpcErrorToError(err)
		switch {
		case errors.Is(gErr, diterrors.ErrNotFound):
			return nil, ErrEmployeeNotFound
		default:
			return nil, fmt.Errorf("can't get employee by id in employees service: %w", gErr)
		}
	}

	emp := er.employeesServiceMapper.EmployeeToEntity(res.GetEmployee())

	emp.Absence = er.indentityActualAbsence(res.GetEmployee().GetAbsenceDocuments())

	emp.HeadStaffPosition, err = er.getHeadStaffPosition(ctx, id)
	if err != nil {
		var valErr diterrors.ValidationError
		switch {
		case errors.As(err, &valErr):
			er.logger.Debug("employeesRepository.getHeadStaffPosition: validation err", ditzap.UUID("employee_id", id), zap.Error(valErr))
		case errors.Is(err, diterrors.ErrNotFound):
			er.logger.Debug("employeesRepository.getHeadStaffPosition: staffPosition head not found", ditzap.UUID("employee_id", id))
		default:
			er.logger.Error("employeesRepository.getHeadStaffPosition", ditzap.UUID("employee_id", id), zap.Error(err))
		}
	}

	emp.HeadManagement, err = er.getHeadManagement(ctx, id)
	if err != nil {
		var valErr diterrors.ValidationError
		switch {
		case errors.As(err, &valErr):
			er.logger.Debug("employeesRepository.getHeadManagement: validation err", ditzap.UUID("employee_id", id), zap.Error(valErr))
		case errors.Is(err, diterrors.ErrNotFound):
			er.logger.Debug("employeesRepository.getHeadManagement: management head not found", ditzap.UUID("employee_id", id))
		default:
			er.logger.Error("employeesRepository.getHeadManagement", ditzap.UUID("employee_id", id), zap.Error(err))
		}
	}

	return emp, nil
}

func (er *employeesRepository) GetByExtIDAndPortalID(ctx context.Context, extID string, portalID int) (*entityEmployee.Employee, error) {
	req := employeev1.CompositeGetRequest{
		Key: &employeev1.CompositeGetRequest_ExtId{
			ExtId: &employeev1.ExtIDWithPortalID{
				Id:       extID,
				PortalId: int32(portalID),
			},
		},
		WithPerson:          true,
		WithOrganization:    true,
		WithProducts:        true,
		WithSubdivision:     true,
		WithStaffPosition:   true,
		WithManagement:      true,
		WithPosition:        true,
		WithWorkplace:       true,
		WithSubdivisionTree: true,
		WithManagementsTree: true,
		WithEmployeeHistory: true,
	}

	res, err := er.employeesAPI.CompositeGet(ctx, &req)
	if err != nil {
		gErr := diterrors.GrpcErrorToError(err)
		switch {
		case errors.Is(gErr, diterrors.ErrNotFound):
			return nil, ErrEmployeeNotFound
		default:
			return nil, fmt.Errorf("can't get employee by ext_id in employees service: %w", gErr)
		}
	}

	emp := er.employeesServiceMapper.EmployeeToEntity(res.GetEmployee())

	emp.Absence = er.indentityActualAbsence(res.GetEmployee().GetAbsenceDocuments())

	emp.HeadStaffPosition, err = er.getHeadStaffPosition(ctx, emp.ID)
	if err != nil {
		var valErr diterrors.ValidationError
		switch {
		case errors.As(err, &valErr):
			er.logger.Debug("employeesRepository.getHeadStaffPosition: validation err", zap.String("ext_employee_id", extID), zap.Error(valErr))
		case errors.Is(err, diterrors.ErrNotFound):
			er.logger.Debug("employeesRepository.getHeadStaffPosition: staffPosition head not found", zap.String("ext_employee_id", extID), zap.Int("portal_id", portalID))
		default:
			er.logger.Error("employeesRepository.getHeadManagement", zap.String("ext_employee_id", extID), zap.Int("portal_id", portalID), zap.Error(err))
		}
	}

	emp.HeadManagement, err = er.getHeadManagement(ctx, emp.ID)
	if err != nil {
		var valErr diterrors.ValidationError
		switch {
		case errors.As(err, &valErr):
			er.logger.Debug("employeesRepository.getHeadManagement: validation err", zap.String("employee_id", extID), zap.Error(valErr))
		case errors.Is(err, diterrors.ErrNotFound):
			er.logger.Debug("employeesRepository.getHeadManagement: staffPosition head not found", zap.String("ext_employee_id", extID), zap.Int("portal_id", portalID))
		default:
			er.logger.Error("employeesRepository.getHeadManagement", zap.String("ext_employee_id", extID), zap.Int("portal_id", portalID), zap.Error(err))
		}
	}

	return emp, nil
}

func (er *employeesRepository) getHeadStaffPosition(ctx context.Context, id uuid.UUID) (*entityEmployee.HeadStaffPosition, error) {
	h, err := er.employeesAPI.GetStaffPositionHead(ctx, &employeev1.GetStaffPositionHeadRequest{Key: &employeev1.GetStaffPositionHeadRequest_Id{Id: id.String()}})
	if err != nil {
		gErr := diterrors.GrpcErrorToError(err)
		switch {
		case errors.Is(gErr, diterrors.ErrNotFound):
			return nil, diterrors.ErrNotFound
		}
		return nil, fmt.Errorf("can't get employees by id in employees service: %w", gErr)
	}

	p := &entityEmployee.HeadStaffPosition{
		ExtPersonID:   h.GetExtPersonId(),
		ExtEmployeeID: h.GetExtEmployeeId(),
		ExtPositionID: h.GetExtPositionId(),
		PositionName:  h.GetPositionName(),
		FirstName:     h.GetFirstName(),
		LastName:      h.GetLastName(),
		MiddleName:    h.GetMiddleName(),
		Gender:        er.employeesServiceMapper.GenderToEntity(h.GetGender()),
	}

	pID, err := uuid.Parse(h.GetPersonId())
	if err != nil {
		p.PersonID = uuid.Nil
	}
	p.PersonID = pID

	eID, err := uuid.Parse(h.GetEmployeeId())
	if err != nil {
		return nil, diterrors.NewValidationError(fmt.Errorf("can't parse employee id in employees service: %w", err))
	}
	p.EmployeeID = eID

	if h.GetImageId() != "" {
		iID, err := uuid.Parse(h.GetImageId())
		if err != nil {
			return nil, diterrors.NewValidationError(fmt.Errorf("can't parse image id in employees service: %w", err))
		}
		p.ImageID = iID
	}

	return p, nil
}

func (er employeesRepository) getHeadManagement(ctx context.Context, id uuid.UUID) (*entityEmployee.HeadManagement, error) {
	h, err := er.employeesAPI.GetManagementHead(ctx, &employeev1.GetManagementHeadRequest{Key: &employeev1.GetManagementHeadRequest_Id{Id: id.String()}})
	if err != nil {
		gErr := diterrors.GrpcErrorToError(err)
		switch {
		case errors.Is(gErr, diterrors.ErrNotFound):
			return nil, diterrors.ErrNotFound
		}
		return nil, fmt.Errorf("can't get employees by id in employees service: %w", gErr)
	}

	p := &entityEmployee.HeadManagement{
		ExtPersonID:   h.GetExtPersonId(),
		ExtEmployeeID: h.GetExtEmployeeId(),
		ExtRoleID:     h.GetExtRoleId(),
		RoleName:      h.GetRoleName(),
		FirstName:     h.GetFirstName(),
		LastName:      h.GetLastName(),
		MiddleName:    h.GetMiddleName(),
		Gender:        er.employeesServiceMapper.GenderToEntity(h.GetGender()),
	}

	pID, err := uuid.Parse(h.GetPersonId())
	if err != nil {
		return nil, diterrors.NewValidationError(fmt.Errorf("can't parse person id in employees service: %w", err))
	}
	p.PersonID = pID

	eID, err := uuid.Parse(h.GetEmployeeId())
	if err != nil {
		return nil, diterrors.NewValidationError(fmt.Errorf("can't parse employee id in employees service: %w", err))
	}
	p.EmployeeID = eID

	if h.GetImageId() != "" {
		iID, err := uuid.Parse(h.GetImageId())
		if err != nil {
			return nil, diterrors.NewValidationError(fmt.Errorf("can't parse image id in employees service: %w", err))
		}
		p.ImageID = iID
	}

	return p, nil
}

func (er employeesRepository) indentityActualAbsence(documents []*employeev1.AbsenceDocument) *entityEmployee.Absence {
	if len(documents) == 0 {
		return nil
	}
	var absencesMap = make(map[int64]*entityEmployee.Absence, len(documents)) // UNIX.start_date -> absenceWithPriority

	// Сбор и предварительная фильтрация отсутствий
	for _, doc := range documents {

		for _, abs := range doc.Absences {

			entityAbs := er.absenceToEntity(abs)

			// Проверка флага видимости актуальности периода
			if abs.IsVisible && er.isAbsenceRelevant(entityAbs.From.AsTime(), entityAbs.To.AsTime(), TWO_WEEKS) {
				if entityAbs == nil || entityAbs.From == nil {
					continue
				}
				// В текущий момент отдаем только два вида отсутствий (Отпуск и Декрет)
				if entityAbs.Type != entityEmployee.AbsenceTypeVacation &&
					entityAbs.Type != entityEmployee.AbsenceTypeDecree {
					er.logger.Debug("employeeRepository.indentityActualAbsence: absence with wrong type", zap.Any("absence", entityAbs))
					continue
				}

				mapAbsence, ok := absencesMap[entityAbs.From.AsTime().Unix()]
				if ok && entityAbs.GetPriority() > mapAbsence.GetPriority() {
					continue
				}
				absencesMap[entityAbs.From.AsTime().Unix()] = entityAbs
			}
		}
	}

	if len(absencesMap) == 0 {
		return nil
	}

	var firstAbsence *entityEmployee.Absence
	for k, v := range absencesMap {
		if v.From.IsZero() {
			er.logger.Warn("employeeRepository.indentityActualAbsence: absence with zero start date", zap.Any("absence", v))
			continue
		}
		if firstAbsence == nil {
			firstAbsence = v
			continue
		}
		if k < firstAbsence.From.AsTime().Unix() {
			firstAbsence = v
		}
	}
	return firstAbsence
}

func (er employeesRepository) isAbsenceRelevant(start, end time.Time, duration time.Duration) bool {
	now := er.tu.New()
	threshold := now.Add(duration)

	if end.Before(*now) {
		return false // Уже закончилось
	}
	if start.After(threshold) {
		return false // Если начало после заданного периода
	}

	return true
}

// TODO: отрефачить в маппер, когда будут енамы в контракте
func (er *employeesRepository) absenceToEntity(absence *employeev1.Absence) *entityEmployee.Absence {
	if absence == nil {
		return nil
	}

	a := &entityEmployee.Absence{
		Name: absence.GetReasonName(),
		Type: er.identityType(absence.GetReasonName()),
	}

	if absence.StartDate != "" {
		d, err := timeUtils.ParseDate(absence.StartDate)
		if err == nil {
			a.From = &d
		}
	}

	if absence.EndDate != "" {
		d, err := timeUtils.ParseDate(absence.EndDate)
		if err == nil {
			a.To = &d
		}
	}

	return a
}

func (er *employeesRepository) identityType(reasonName string) entityEmployee.AbsenceType {
	reasonLower := strings.ToLower(reasonName)
	switch reasonLower {
	case "декрет":
		return entityEmployee.AbsenceTypeDecree
	case "работа в декрете":
		return entityEmployee.AbsenceTypeDecreeWork
	case "больничный":
		return entityEmployee.AbsenceTypeBusinessTrip
	case "отпуск":
		return entityEmployee.AbsenceTypeVacation
	default:
		return entityEmployee.AbsenceTypeUnknown
	}
}
