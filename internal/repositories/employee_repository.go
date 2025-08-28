package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	personv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/person/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/util/converter"
)

const (
	employeePrefix = "emp:"
	employeeTTL    = 2 * 24 * time.Hour
)

type employeeRepository struct {
	basePrefix   string
	source       CacheSource
	employeesAPI employeev1.EmployeesAPIClient
	personsAPI   personv1.PersonAPIClient
	logger       ditzap.Logger
}

func NewEmployeeRepository(
	basePrefix string,
	cacheSource CacheSource,
	employeesAPI employeev1.EmployeesAPIClient,
	personsAPI personv1.PersonAPIClient,
	logger ditzap.Logger,
) *employeeRepository {
	return &employeeRepository{
		basePrefix:   basePrefix,
		source:       cacheSource,
		employeesAPI: employeesAPI,
		personsAPI:   personsAPI,
		logger:       logger,
	}
}
func (er *employeeRepository) getKey(key string) string {
	return er.basePrefix + employeePrefix + key
}

func (er *employeeRepository) Save(ctx context.Context, key string, employees []entity.EmployeeInfo) error {
	if key == "" {
		return errors.New("ключ для сохранения employees пустой")
	}

	data, err := json.Marshal(employees)
	if err != nil {
		return err
	}

	err = er.source.SetEx(ctx, er.getKey(key), data, employeeTTL)
	if err != nil {
		er.logger.Error("не удалось сохранить employees", zap.Error(err))
		return err
	}
	return nil
}

func (er *employeeRepository) Get(ctx context.Context, key string) ([]entity.EmployeeInfo, error) {
	dataStr, err := er.source.Get(ctx, er.getKey(key))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		er.logger.Error("не удалось получить employees", zap.Error(err))
		return nil, err
	}

	data := converter.StringToBytes(dataStr)

	var employees []entity.EmployeeInfo
	if err := json.Unmarshal(data, &employees); err != nil {
		return nil, err
	}

	return employees, nil
}

func (er *employeeRepository) GetPersonIDByEmployeeEmail(ctx context.Context, email string) (uuid.UUID, error) {
	if email == "" {
		return uuid.Nil, diterrors.NewValidationError(fmt.Errorf("employee email is empty"))
	}

	request := &employeev1.CompositeGetRequest{
		Key: &employeev1.CompositeGetRequest_Email{
			Email: email,
		},
		WithPerson: true,
	}

	response, err := er.employeesAPI.CompositeGet(ctx, request)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't get employee by email in employees service: %w", diterrors.GrpcErrorToError(err))
	}

	personID, err := uuid.Parse(response.GetEmployee().GetPerson().GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't parse person id: %w", err)
	}

	return personID, nil
}

func (er *employeeRepository) GetEmployeesInfoByPersonID(ctx context.Context, personID uuid.UUID) ([]entity.EmployeeInfo, error) {
	if personID == uuid.Nil {
		return nil, diterrors.NewValidationError(fmt.Errorf("person id is nil uuid"))
	}

	request := &personv1.CompositeGetRequest{
		Key: &personv1.CompositeGetRequest_Id{
			Id: personID.String(),
		},
		WithEmployees:     true,
		WithOrganizations: true,
	}

	response, err := er.personsAPI.CompositeGet(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("can't get person by id in employees service: %w", diterrors.GrpcErrorToError(err))
	}

	employees := response.GetEmployees()
	employeesInfo := make([]entity.EmployeeInfo, 0, len(employees))
	for _, employee := range employees {
		employeesInfo = append(employeesInfo, entity.EmployeeInfo{
			CloudID: employee.GetPerson().GetCloudId(),
			Inn:     employee.GetOrganization().GetInn(),
			OrgID:   employee.GetOrganization().GetId(),
			FIO:     employee.GetFullName(),
			SNILS:   employee.GetPerson().GetSnils(),
		})
	}

	return employeesInfo, nil
}
