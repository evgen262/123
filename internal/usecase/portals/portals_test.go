package portals

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestPortalsUseCase_Filter(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		logger     *ditzap.MockLogger
		repository *MockPortalRepository
	}
	type args struct {
		ctx  context.Context
		opts portal.PortalsFilterOptions
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Portal, error)
	}{
		{
			name: "error",
			args: args{
				ctx: ctx,
				opts: portal.PortalsFilterOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := errors.New("some service error")
				f.repository.EXPECT().Filter(a.ctx, a.opts).Return(nil, testErr)
				f.logger.EXPECT().Error("can't filter portal", zap.Error(testErr))
				return nil, fmt.Errorf("can't filter portal: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				opts: portal.PortalsFilterOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portals := []*portal.Portal{
					{
						Id:            1,
						FullName:      "Test portal 1",
						ShortName:     "test 1",
						Url:           "https://test1.mos.ru",
						LogoUrl:       "https://test1.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test1.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"}},
						IsDeleted:     false,
					},
					{
						Id:            2,
						FullName:      "Test portal 2",
						ShortName:     "test 2",
						Url:           "https://test2.mos.ru",
						LogoUrl:       "https://test2.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test2.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "3c5cbb16-011a-310e-97e2-565400a26506"}},
						IsDeleted:     false,
					},
				}
				f.repository.EXPECT().Filter(a.ctx, a.opts).Return(portals, nil)
				return portals, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repository: NewMockPortalRepository(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			p := NewPortalsUseCase(f.repository, f.logger)
			got, err := p.Filter(tt.args.ctx, tt.args.opts)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestPortalsUseCase_GetByEmployees(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		logger     *ditzap.MockLogger
		repository *MockPortalRepository
	}
	type args struct {
		ctx       context.Context
		employees []portal.EmployeeInfo
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Portal, error)
	}{
		{
			name: "validation error",
			args: args{
				ctx: ctx,
				employees: []portal.EmployeeInfo{
					{
						Inn:   "770123456789",
						OrgID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:   "Тестов Тест Тестович",
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				validationErr := diterrors.NewValidationError(errors.New("some validation error"))

				opts := portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"770123456789"},
					OnlyLinked: true,
				}

				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, validationErr)
				return nil, validationErr
			},
		},
		{
			name: "not found error",
			args: args{
				ctx: ctx,
				employees: []portal.EmployeeInfo{
					{
						Inn:   "770123456789",
						OrgID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:   "Тестов Тест Тестович",
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				opts := portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"770123456789"},
					OnlyLinked: true,
				}

				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, repositories.ErrNotFound)
				return nil, nil
			},
		},
		{
			name: "localized error",
			args: args{
				ctx: ctx,
				employees: []portal.EmployeeInfo{
					{
						Inn:   "770123456789",
						OrgID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:   "Тестов Тест Тестович",
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, errors.New("some error"))

				opts := portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"770123456789"},
					OnlyLinked: true,
				}

				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, testErr)
				f.logger.EXPECT().Error("can't filter portal by employees", zap.Error(testErr.Unwrap()))
				return nil, fmt.Errorf("can't filter portal by employees: %w", testErr)
			},
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				employees: []portal.EmployeeInfo{
					{
						Inn:   "770123456789",
						OrgID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:   "Тестов Тест Тестович",
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := errors.New("some service error")

				opts := portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"770123456789"},
					OnlyLinked: true,
				}

				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, testErr)
				f.logger.EXPECT().Error("can't filter portal by employees", zap.Error(testErr))
				return nil, fmt.Errorf("can't filter portal by employees: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				employees: []portal.EmployeeInfo{
					{
						Inn:   "770123456789",
						OrgID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						FIO:   "Тестов Тест Тестович",
					},
					{
						Inn:   "779876543210",
						OrgID: "dc858901-57c4-44a7-991d-b0cbc9108678",
						FIO:   "Тестов Тест Тестович",
					},
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsArr := []*portal.Portal{
					{
						Id:            1,
						FullName:      "Test portal 1",
						ShortName:     "test 1",
						Url:           "https://test1.mos.ru",
						LogoUrl:       "https://test1.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test1.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"}},
						IsDeleted:     false,
					},
					{
						Id:            2,
						FullName:      "Test portal 2",
						ShortName:     "test 2",
						Url:           "https://test2.mos.ru",
						LogoUrl:       "https://test2.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test2.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "3c5cbb16-011a-310e-97e2-565400a26506"}},
						IsDeleted:     false,
					},
				}
				opts := portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"770123456789", "779876543210"},
					OnlyLinked: true,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(portalsArr, nil)
				return portalsArr, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repository: NewMockPortalRepository(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			p := NewPortalsUseCase(f.repository, f.logger)
			got, err := p.GetByEmployees(tt.args.ctx, tt.args.employees)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

//nolint:funlen
func TestPortalsUseCase_GetAll(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		logger     *ditzap.MockLogger
		repository *MockPortalRepository
	}
	type args struct {
		ctx  context.Context
		opts portal.GetAllOptions
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Portal, error)
	}{
		{
			name: "localized error",
			args: args{
				ctx: ctx,
				opts: portal.GetAllOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := errors.New("some service error")
				localizedError := diterrors.NewLocalizedError(diterrors.LocalizeLocale, testErr)
				f.repository.EXPECT().Filter(a.ctx, portal.PortalsFilterOptions{
					WithDeleted: a.opts.WithDeleted,
					OnlyLinked:  a.opts.OnlyLinked,
				}).Return(nil, localizedError)
				f.logger.EXPECT().Error("can't get all portal", zap.Error(localizedError.Unwrap()))
				return nil, fmt.Errorf("can't get all portal: %w", localizedError)
			},
		},
		{
			name: "not found error",
			args: args{
				ctx: ctx,
				opts: portal.GetAllOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				f.repository.EXPECT().Filter(a.ctx, portal.PortalsFilterOptions{
					WithDeleted: a.opts.WithDeleted,
					OnlyLinked:  a.opts.OnlyLinked,
				}).Return(nil, repositories.ErrNotFound)
				return nil, nil
			},
		},
		{
			name: "validation error",
			args: args{
				ctx: ctx,
				opts: portal.GetAllOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := errors.New("some service error")
				validationErr := diterrors.NewValidationError(testErr)
				f.repository.EXPECT().Filter(a.ctx, portal.PortalsFilterOptions{
					WithDeleted: a.opts.WithDeleted,
					OnlyLinked:  a.opts.OnlyLinked,
				}).Return(nil, validationErr)
				return nil, validationErr
			},
		},
		{
			name: "default error",
			args: args{
				ctx: ctx,
				opts: portal.GetAllOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				testErr := errors.New("some service error")
				f.repository.EXPECT().Filter(a.ctx, portal.PortalsFilterOptions{
					WithDeleted: a.opts.WithDeleted,
					OnlyLinked:  a.opts.OnlyLinked,
				}).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get all portal", zap.Error(testErr))
				return nil, fmt.Errorf("can't get all portal: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				opts: portal.GetAllOptions{
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsArr := []*portal.Portal{
					{
						Id:            1,
						FullName:      "Test portal 1",
						ShortName:     "test 1",
						Url:           "https://test1.mos.ru",
						LogoUrl:       "https://test1.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test1.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"}},
						IsDeleted:     false,
					},
					{
						Id:            2,
						FullName:      "Test portal 2",
						ShortName:     "test 2",
						Url:           "https://test2.mos.ru",
						LogoUrl:       "https://test2.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test2.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "3c5cbb16-011a-310e-97e2-565400a26506"}},
						IsDeleted:     false,
					},
				}
				f.repository.EXPECT().Filter(a.ctx, portal.PortalsFilterOptions{
					WithDeleted: a.opts.WithDeleted,
					OnlyLinked:  a.opts.OnlyLinked,
				}).Return(portalsArr, nil)
				return portalsArr, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repository: NewMockPortalRepository(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			p := NewPortalsUseCase(f.repository, f.logger)
			got, err := p.GetAll(tt.args.ctx, tt.args.opts)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}

//nolint:funlen
func TestPortalsUseCase_Get(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		logger     *ditzap.MockLogger
		repository *MockPortalRepository
	}
	type args struct {
		ctx         context.Context
		id          int
		withDeleted bool
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Portal, error)
	}{
		{
			name: "correct",
			args: args{
				ctx:         ctx,
				id:          1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				portalsArr := []*portal.Portal{
					{
						Id:            1,
						FullName:      "Test portal 1",
						ShortName:     "test 1",
						Url:           "https://test1.mos.ru",
						LogoUrl:       "https://test1.mos.ru/path/to/logo.jpg",
						ChatUrl:       "https://test1.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{{Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"}},
						IsDeleted:     false,
					},
				}
				opts := portal.PortalsFilterOptions{
					PortalIDs:   portal.PortalIDs{portal.PortalID(a.id)},
					WithDeleted: a.withDeleted,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(portalsArr, nil)
				return portalsArr[0], nil
			},
		},
		{
			name: "not found err",
			args: args{
				ctx:         ctx,
				id:          1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				opts := portal.PortalsFilterOptions{
					PortalIDs:   portal.PortalIDs{portal.PortalID(a.id)},
					WithDeleted: a.withDeleted,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, repositories.ErrNotFound)
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "localized err",
			args: args{
				ctx:         ctx,
				id:          1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				testErr := errors.New("some service error")
				localizedErr := diterrors.NewLocalizedError(diterrors.LocalizeLocale, testErr)
				opts := portal.PortalsFilterOptions{
					PortalIDs:   portal.PortalIDs{portal.PortalID(a.id)},
					WithDeleted: a.withDeleted,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, localizedErr)
				f.logger.EXPECT().Error("can't get portal", zap.Error(localizedErr.Unwrap()))
				return nil, fmt.Errorf("can't get portal: %w", localizedErr)
			},
		},
		{
			name: "validation err",
			args: args{
				ctx:         ctx,
				id:          1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				testErr := errors.New("some service error")
				validationErr := diterrors.NewValidationError(testErr)
				opts := portal.PortalsFilterOptions{
					PortalIDs:   portal.PortalIDs{portal.PortalID(a.id)},
					WithDeleted: a.withDeleted,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, validationErr)
				return nil, validationErr
			},
		},
		{
			name: "filter err",
			args: args{
				ctx:         ctx,
				id:          1,
				withDeleted: false,
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				testErr := errors.New("some service error")
				opts := portal.PortalsFilterOptions{
					PortalIDs:   portal.PortalIDs{portal.PortalID(a.id)},
					WithDeleted: a.withDeleted,
				}
				f.repository.EXPECT().Filter(a.ctx, opts).Return(nil, testErr)
				f.logger.EXPECT().Error("can't get portal", zap.Error(testErr))
				return nil, fmt.Errorf("can't get portal: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repository: NewMockPortalRepository(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			p := NewPortalsUseCase(f.repository, f.logger)
			got, err := p.Get(tt.args.ctx, tt.args.id, tt.args.withDeleted)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
