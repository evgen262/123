package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	organizationv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/organization/v1"
	personv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/person/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

func Test_employeeRepository_getKey(t *testing.T) {
	testBasePrefix := "cache1:"
	type fields struct {
		basePrefix   string
		source       *MockCacheSource
		employeesAPI *employeev1.MockEmployeesAPIClient
		personsAPI   *personv1.MockPersonAPIClient
		logger       *ditzap.MockLogger
	}
	type args struct {
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) string
	}{
		{
			name: "correct",
			args: args{cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"},
			want: func(a args, f fields) string {
				return f.basePrefix + employeePrefix + a.cloudID
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				basePrefix:   testBasePrefix,
				source:       NewMockCacheSource(ctrl),
				employeesAPI: employeev1.NewMockEmployeesAPIClient(ctrl),
				personsAPI:   personv1.NewMockPersonAPIClient(ctrl),
				logger:       ditzap.NewMockLogger(ctrl),
			}
			want := tt.want(tt.args, f)

			er := NewEmployeeRepository(f.basePrefix, f.source, f.employeesAPI, f.personsAPI, f.logger)
			got := er.getKey(tt.args.cloudID)

			assert.Equal(t, want, got)
		})
	}
}

func Test_employeeRepository_Save(t *testing.T) {
	testBasePrefix := "cache1:"
	testCtx := context.TODO()

	type fields struct {
		basePrefix   string
		source       *MockCacheSource
		employeesAPI *employeev1.MockEmployeesAPIClient
		personsAPI   *personv1.MockPersonAPIClient
		logger       *ditzap.MockLogger
	}
	type args struct {
		ctx       context.Context
		key       string
		employees []entity.EmployeeInfo
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "empty key error",
			args: args{
				ctx:       testCtx,
				key:       "",
				employees: []entity.EmployeeInfo{},
			},
			want: func(a args, f fields) error {
				return errors.New("ключ для сохранения employees пустой")
			},
		},
		{
			name: "error",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				employees: []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "77123456789",
					},
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "779876543210",
					},
				},
			},
			want: func(a args, f fields) error {
				testData, _ := json.Marshal(a.employees)
				testErr := errors.New("some save error")
				testKey := f.basePrefix + employeePrefix + a.key

				f.source.EXPECT().SetEx(a.ctx, testKey, testData, employeeTTL).Return(testErr)
				f.logger.EXPECT().Error("не удалось сохранить employees", zap.Error(testErr))
				return testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				employees: []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "77123456789",
					},
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "779876543210",
					},
				},
			},
			want: func(a args, f fields) error {
				testData, _ := json.Marshal(a.employees)
				testKey := f.basePrefix + employeePrefix + a.key

				f.source.EXPECT().SetEx(a.ctx, testKey, testData, employeeTTL).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				basePrefix:   testBasePrefix,
				source:       NewMockCacheSource(ctrl),
				employeesAPI: employeev1.NewMockEmployeesAPIClient(ctrl),
				personsAPI:   personv1.NewMockPersonAPIClient(ctrl),
				logger:       ditzap.NewMockLogger(ctrl),
			}
			er := NewEmployeeRepository(f.basePrefix, f.source, f.employeesAPI, f.personsAPI, f.logger)

			wantErr := tt.want(tt.args, f)
			err := er.Save(tt.args.ctx, tt.args.key, tt.args.employees)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_employeeRepository_Get(t *testing.T) {
	testBasePrefix := "cache1:"
	testCtx := context.TODO()

	type fields struct {
		basePrefix   string
		source       *MockCacheSource
		employeesAPI *employeev1.MockEmployeesAPIClient
		personsAPI   *personv1.MockPersonAPIClient
		logger       *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]entity.EmployeeInfo, error)
	}{
		{
			name: "error",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testKey := f.basePrefix + employeePrefix + a.key
				testErr := errors.New("some save error")

				f.source.EXPECT().Get(a.ctx, testKey).Return("", testErr)
				f.logger.EXPECT().Error("не удалось получить employees", zap.Error(testErr))
				return nil, testErr
			},
		},
		{
			name: "not found",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testKey := f.basePrefix + employeePrefix + a.key

				f.source.EXPECT().Get(a.ctx, testKey).Return("", redis.Nil)
				return nil, ErrNotFound
			},
		},
		{
			name: "err unmarshal",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testKey := f.basePrefix + employeePrefix + a.key
				testStr := "[{\"some key\": \"some value\"}"
				testErr := errors.New("unexpected end of JSON input")

				f.source.EXPECT().Get(a.ctx, testKey).Return(testStr, nil)
				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx: testCtx,
				key: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testKey := f.basePrefix + employeePrefix + a.key
				testStr := "[{\"cloud_id\": \"342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f\", \"inn\": \"779876543210\"},{\"cloud_id\": \"342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f\", \"inn\": \"771234567890\"}]"

				f.source.EXPECT().Get(a.ctx, testKey).Return(testStr, nil)
				return []entity.EmployeeInfo{
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "779876543210",
					},
					{
						CloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:     "771234567890",
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				basePrefix:   testBasePrefix,
				source:       NewMockCacheSource(ctrl),
				employeesAPI: employeev1.NewMockEmployeesAPIClient(ctrl),
				personsAPI:   personv1.NewMockPersonAPIClient(ctrl),
				logger:       ditzap.NewMockLogger(ctrl),
			}
			er := NewEmployeeRepository(f.basePrefix, f.source, f.employeesAPI, f.personsAPI, f.logger)

			want, wantErr := tt.want(tt.args, f)
			got, err := er.Get(tt.args.ctx, tt.args.key)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, want, got)
			}
		})
	}
}

func Test_employeeRepository_GetPersonIDByEmployeeEmail(t *testing.T) {
	type fields struct {
		basePrefix   string
		source       *MockCacheSource
		employeesAPI *employeev1.MockEmployeesAPIClient
		personsAPI   *personv1.MockPersonAPIClient
		logger       *ditzap.MockLogger
	}
	type args struct {
		ctx   context.Context
		email string
	}

	testBasePrefix := "cache1:"
	ctx := context.TODO()
	testErr := errors.New("test error")
	testPersonID := uuid.New()

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (uuid.UUID, error)
	}{
		{
			name: "empty email",
			args: args{
				ctx:   ctx,
				email: "",
			},
			want: func(a args, f fields) (uuid.UUID, error) {
				return uuid.Nil, diterrors.NewValidationError(fmt.Errorf("employee email is empty"))
			},
		},
		{
			name: "composite get err",
			args: args{
				ctx:   ctx,
				email: "testEmail",
			},
			want: func(a args, f fields) (uuid.UUID, error) {
				testRequest := &employeev1.CompositeGetRequest{
					Key: &employeev1.CompositeGetRequest_Email{
						Email: a.email,
					},
					WithPerson: true,
				}
				f.employeesAPI.EXPECT().CompositeGet(a.ctx, testRequest).Return(nil, testErr)
				return uuid.Nil, fmt.Errorf("can't get employee by email in employees service: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
		{
			name: "parse person id err",
			args: args{
				ctx:   ctx,
				email: "testEmail",
			},
			want: func(a args, f fields) (uuid.UUID, error) {
				testRequest := &employeev1.CompositeGetRequest{
					Key: &employeev1.CompositeGetRequest_Email{
						Email: a.email,
					},
					WithPerson: true,
				}
				testResponse := &employeev1.CompositeGetResponse{
					Employee: &employeev1.CompositeEmployee{
						Person: &employeev1.CompositeEmployee_Person{
							Id: "invalid",
						},
					},
				}
				f.employeesAPI.EXPECT().CompositeGet(a.ctx, testRequest).Return(testResponse, nil)
				testParseErr := fmt.Errorf(fmt.Sprintf("invalid UUID length: %d",
					len(testResponse.GetEmployee().GetPerson().GetId())))
				return uuid.Nil, fmt.Errorf("can't parse person id: %w", testParseErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:   ctx,
				email: "testEmail",
			},
			want: func(a args, f fields) (uuid.UUID, error) {
				testRequest := &employeev1.CompositeGetRequest{
					Key: &employeev1.CompositeGetRequest_Email{
						Email: a.email,
					},
					WithPerson: true,
				}
				testResponse := &employeev1.CompositeGetResponse{
					Employee: &employeev1.CompositeEmployee{
						Person: &employeev1.CompositeEmployee_Person{
							Id: testPersonID.String(),
						},
					},
				}
				f.employeesAPI.EXPECT().CompositeGet(a.ctx, testRequest).Return(testResponse, nil)
				return testPersonID, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				basePrefix:   testBasePrefix,
				source:       NewMockCacheSource(ctrl),
				employeesAPI: employeev1.NewMockEmployeesAPIClient(ctrl),
				personsAPI:   personv1.NewMockPersonAPIClient(ctrl),
				logger:       ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			er := NewEmployeeRepository(f.basePrefix, f.source, f.employeesAPI, f.personsAPI, f.logger)
			got, err := er.GetPersonIDByEmployeeEmail(tt.args.ctx, tt.args.email)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Equal(t, uuid.Nil, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func Test_employeeRepository_GetEmployeesInfoByPersonID(t *testing.T) {
	type fields struct {
		basePrefix   string
		source       *MockCacheSource
		employeesAPI *employeev1.MockEmployeesAPIClient
		personsAPI   *personv1.MockPersonAPIClient
		logger       *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		personID uuid.UUID
	}

	testBasePrefix := "cache1:"
	ctx := context.TODO()
	testErr := errors.New("test error")
	testPersonID := uuid.New()

	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]entity.EmployeeInfo, error)
	}{
		{
			name: "nil uuid person id",
			args: args{
				ctx:      ctx,
				personID: uuid.Nil,
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				return nil, diterrors.NewValidationError(fmt.Errorf("person id is nil uuid"))
			},
		},
		{
			name: "composite get err",
			args: args{
				ctx:      ctx,
				personID: testPersonID,
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testRequest := &personv1.CompositeGetRequest{
					Key: &personv1.CompositeGetRequest_Id{
						Id: a.personID.String(),
					},
					WithEmployees:     true,
					WithOrganizations: true,
				}
				f.personsAPI.EXPECT().CompositeGet(a.ctx, testRequest).Return(nil, testErr)
				return nil, fmt.Errorf("can't get person by id in employees service: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
		{
			name: "correct",
			args: args{
				ctx:      ctx,
				personID: testPersonID,
			},
			want: func(a args, f fields) ([]entity.EmployeeInfo, error) {
				testRequest := &personv1.CompositeGetRequest{
					Key: &personv1.CompositeGetRequest_Id{
						Id: a.personID.String(),
					},
					WithEmployees:     true,
					WithOrganizations: true,
				}
				testResponse := &personv1.CompositeGetResponse{
					Employees: []*employeev1.CompositeEmployee{
						{
							Person: &employeev1.CompositeEmployee_Person{
								CloudId: "testCloudID",
								Snils:   "testSNILS",
							},
							Organization: &organizationv1.Organization{
								Id:  "testOrganizationID",
								Inn: "testINN",
							},
							FullName: "testFullName",
						},
					},
				}
				f.personsAPI.EXPECT().CompositeGet(a.ctx, testRequest).Return(testResponse, nil)
				return []entity.EmployeeInfo{
					{
						CloudID: "testCloudID",
						Inn:     "testINN",
						OrgID:   "testOrganizationID",
						FIO:     "testFullName",
						SNILS:   "testSNILS",
					},
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				basePrefix:   testBasePrefix,
				source:       NewMockCacheSource(ctrl),
				employeesAPI: employeev1.NewMockEmployeesAPIClient(ctrl),
				personsAPI:   personv1.NewMockPersonAPIClient(ctrl),
				logger:       ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			er := NewEmployeeRepository(f.basePrefix, f.source, f.employeesAPI, f.personsAPI, f.logger)
			got, err := er.GetEmployeesInfoByPersonID(tt.args.ctx, tt.args.personID)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
