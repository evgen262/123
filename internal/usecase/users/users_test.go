package users

import (
	"context"
	"errors"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	entityEmployee "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/employee"
	entityUser "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/user"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/usecase"
)

func Test_usersUseCase_GetMe(t *testing.T) {
	type fields struct {
		repo   *MockEmployeesRepository
		logger *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
	}

	ctx := context.TODO()
	testID := uuid.NewString()

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*entityUser.UserInfo, error)
	}{
		{
			name: "nil session err",
			args: args{
				ctx: entity.WithSession(ctx, nil),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				_, testErr := entity.SessionFromContext(a.ctx)
				f.logger.EXPECT().Error("can't get session from context", zap.Error(testErr))
				return nil, usecase.ErrGetSessionFromContext
			},
		},
		{
			name: "get employee ext id err",
			args: args{
				ctx: entity.WithSession(ctx, &entityAuth.Session{
					User: &entityAuth.User{},
				}),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				testSession, _ := entity.SessionFromContext(a.ctx)
				f.logger.EXPECT().Error(ErrEmptySessionEmployeeExtID.Error(),
					zap.Object("session", testSession),
				)
				return nil, ErrEmptySessionEmployeeExtID
			},
		},
		{
			name: "get active portal err",
			args: args{
				ctx: entity.WithSession(ctx, &entityAuth.Session{
					User: &entityAuth.User{
						Employee: &entityAuth.Employee{
							ExtID: testID,
						},
					},
				}),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				testSession, _ := entity.SessionFromContext(a.ctx)
				f.logger.EXPECT().Error(ErrZeroSessionActivePortalID.Error(),
					zap.Object("session", testSession),
				)
				return nil, ErrZeroSessionActivePortalID
			},
		},
		{
			name: "get composite employee err",
			args: args{
				ctx: entity.WithSession(ctx, &entityAuth.Session{
					User: &entityAuth.User{
						Employee: &entityAuth.Employee{
							ExtID: testID,
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							ID: 1,
						},
					},
				}),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				testErr := errors.New("test error")
				testSession, _ := entity.SessionFromContext(a.ctx)
				f.repo.EXPECT().GetByExtIDAndPortalID(
					a.ctx,
					testSession.GetUser().GetEmployee().GetExtID(),
					testSession.GetActivePortal().Portal.ID,
				).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get composite employee in repository",
					zap.Object("session", testSession),
					zap.Error(testErr),
				)
				return new(entityUser.UserInfo), nil
			},
		},
		{
			name: "get nil composite employee err",
			args: args{
				ctx: entity.WithSession(ctx, &entityAuth.Session{
					User: &entityAuth.User{
						Employee: &entityAuth.Employee{
							ExtID: testID,
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							ID: 1,
						},
					},
				}),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				testSession, _ := entity.SessionFromContext(a.ctx)
				f.repo.EXPECT().GetByExtIDAndPortalID(
					a.ctx,
					testSession.GetUser().GetEmployee().GetExtID(),
					testSession.GetActivePortal().Portal.ID,
				).Return(nil, nil)
				f.logger.EXPECT().Error("received nil composite employee from repository",
					zap.Object("session", testSession),
				)
				return new(entityUser.UserInfo), nil
			},
		},
		{
			name: "correct",
			args: args{
				ctx: entity.WithSession(ctx, &entityAuth.Session{
					User: &entityAuth.User{
						Employee: &entityAuth.Employee{
							ExtID: testID,
						},
					},
					ActivePortal: &entityAuth.ActivePortal{
						Portal: entityAuth.Portal{
							ID: 1,
						},
					},
				}),
			},
			want: func(a args, f fields) (*entityUser.UserInfo, error) {
				testSession, _ := entity.SessionFromContext(a.ctx)
				testEmployee := &entityEmployee.Employee{}
				f.repo.EXPECT().GetByExtIDAndPortalID(
					a.ctx,
					testSession.GetUser().GetEmployee().GetExtID(),
					testSession.GetActivePortal().Portal.ID,
				).Return(testEmployee, nil)
				return &entityUser.UserInfo{
					User: entityUser.ShortUser{
						LastName:   testEmployee.Person.LastName,
						FirstName:  testEmployee.Person.FirstName,
						MiddleName: testEmployee.Person.MiddleName,
						ImageID:    testEmployee.Person.ImageID,
						Gender:     testEmployee.Person.Gender,
						PortalData: entityUser.PortalData{
							PersonID:   testEmployee.Person.ExtID,
							EmployeeID: testSession.GetUser().GetEmployee().GetExtID(),
						},
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
				repo:   NewMockEmployeesRepository(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			ei := NewUsersInteractor(f.repo, f.logger)
			got, err := ei.GetMe(tt.args.ctx)
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
