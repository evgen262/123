package employees

import (
	"context"
	"errors"
	"fmt"
	"testing"

	employeev1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/employees/employee/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
)

func Test_employeesRepository_Get(t *testing.T) {
	type fields struct {
		client *employeev1.MockEmployeesAPIClient
		mapper *MockMapperEmployees
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	ctx := context.TODO()
	testEmployeeID := uuid.New()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityEmployee.Employee, error)
	}{
		{
			name: "get composite employee err",
			args: args{
				ctx: ctx,
				id:  testEmployeeID,
			},
			want: func(a args, f fields) (*entityEmployee.Employee, error) {
				testRequest := &employeev1.CompositeGetRequest{
					Key: &employeev1.CompositeGetRequest_Id{
						Id: a.id.String(),
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
				f.client.EXPECT().CompositeGet(a.ctx, testRequest).Return(nil, testErr)
				return nil, fmt.Errorf("can't get employee by id in employees service: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				client: employeev1.NewMockEmployeesAPIClient(ctrl),
				mapper: NewMockMapperEmployees(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			er := NewEmployeesRepository(f.client, f.mapper, f.tu, f.logger)
			got, err := er.Get(tt.args.ctx, tt.args.id)
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

func Test_employeesRepository_GetByExtIDAndPortalID(t *testing.T) {
	type fields struct {
		client *employeev1.MockEmployeesAPIClient
		mapper *MockMapperEmployees
		tu     *timeUtils.MockTimeUtils
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx      context.Context
		extID    string
		portalID int
	}

	ctx := context.TODO()
	testEmployeeExtID := uuid.NewString()
	testErr := errors.New("test error")

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityEmployee.Employee, error)
	}{
		{
			name: "get composite employee err",
			args: args{
				ctx:      ctx,
				extID:    testEmployeeExtID,
				portalID: 0,
			},
			want: func(a args, f fields) (*entityEmployee.Employee, error) {
				testRequest := &employeev1.CompositeGetRequest{
					Key: &employeev1.CompositeGetRequest_ExtId{
						ExtId: &employeev1.ExtIDWithPortalID{
							Id:       a.extID,
							PortalId: int32(a.portalID),
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
				f.client.EXPECT().CompositeGet(a.ctx, testRequest).Return(nil, testErr)
				return nil, fmt.Errorf("can't get employee by ext_id in employees service: %w", diterrors.GrpcErrorToError(testErr))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				client: employeev1.NewMockEmployeesAPIClient(ctrl),
				mapper: NewMockMapperEmployees(ctrl),
				tu:     timeUtils.NewMockTimeUtils(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			er := NewEmployeesRepository(f.client, f.mapper, f.tu, f.logger)
			got, err := er.GetByExtIDAndPortalID(tt.args.ctx, tt.args.extID, tt.args.portalID)
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
