package portal

import (
	"context"
	"errors"
	"fmt"
	"testing"

	organizationsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/organizations/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/shared/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_portalOrganizationsRepository_LinkOrganizationsToPortal(t *testing.T) {
	type fields struct {
		client *organizationsv1.MockOrganizationsAPIClient
		mapper *MockOrganizationsMapper
	}
	type args struct {
		ctx      context.Context
		portalId portal.PortalID
		ids      portal.OrganizationIDs
	}
	ctx := context.TODO()
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "correct",
			args: args{
				ctx:      ctx,
				portalId: 1,
				ids:      portal.OrganizationIDs{},
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Link(a.ctx, &organizationsv1.LinkRequest{
					PortalId: int32(a.portalId),
					OrgIds:   a.ids.ToStringSlice(),
				}).Return(nil, nil)
				return nil
			},
		},
		{
			name: "link to portal service error",
			args: args{
				ctx:      ctx,
				portalId: 1,
				ids:      portal.OrganizationIDs{},
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().Link(a.ctx, &organizationsv1.LinkRequest{
					PortalId: int32(a.portalId),
					OrgIds:   a.ids.ToStringSlice(),
				}).Return(nil, testErr)
				return fmt.Errorf("can't link to portal #%d: %w", 1, testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: organizationsv1.NewMockOrganizationsAPIClient(ctrl),
				mapper: NewMockOrganizationsMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewOrganizationsRepository(f.client, f.mapper)
			err := repo.LinkOrganizationsToPortal(tt.args.ctx, tt.args.portalId, tt.args.ids)
			assert.Equal(t, wantErr, err)
		})
	}
}

func Test_portalOrganizationsRepository_UnlinkOrganizations(t *testing.T) {
	type fields struct {
		client *organizationsv1.MockOrganizationsAPIClient
		mapper *MockOrganizationsMapper
	}
	type args struct {
		ctx context.Context
		ids portal.OrganizationIDs
	}
	ctx := context.TODO()
	testErr := errors.New("error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "correct",
			args: args{
				ctx: ctx,
				ids: portal.OrganizationIDs{},
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().
					Unlink(a.ctx, &organizationsv1.UnlinkRequest{OrgIds: a.ids.ToStringSlice()}).
					Return(nil, nil)
				return nil
			},
		},
		{
			name: "unlink organizations service error",
			args: args{
				ctx: ctx,
				ids: portal.OrganizationIDs{},
			},
			want: func(a args, f fields) error {
				f.client.EXPECT().
					Unlink(a.ctx, &organizationsv1.UnlinkRequest{OrgIds: a.ids.ToStringSlice()}).
					Return(nil, testErr)
				return fmt.Errorf("can't unlink organizations: %w", testErr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: organizationsv1.NewMockOrganizationsAPIClient(ctrl),
				mapper: NewMockOrganizationsMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewOrganizationsRepository(f.client, f.mapper)
			err := repo.UnlinkOrganizations(tt.args.ctx, tt.args.ids)
			assert.Equal(t, wantErr, err)
		})
	}
}

func Test_portalOrganizationsRepository_Filter(t *testing.T) {
	ctx := context.TODO()
	type fields struct {
		client *organizationsv1.MockOrganizationsAPIClient
		mapper *MockOrganizationsMapper
	}
	type args struct {
		ctx        context.Context
		filters    portal.OrganizationsFilters
		pagination *entity.StringPagination
		options    portal.OrganizationsFilterOptions
	}
	testInt := 10
	testErr := fmt.Errorf("test error")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.OrganizationsWithPagination, error)
	}{
		{
			name: "invalid argument",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Ids: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Ids_{
						Ids: &organizationsv1.FilterRequest_Ids{
							Ids: a.filters.Ids,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				testErr := status.Error(codes.InvalidArgument, "bad code parameter")
				msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, testErr)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(nil, testErr)
				return nil, diterrors.NewValidationError(msg)
			},
		},
		{
			name: "not found",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Ids: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Ids_{
						Ids: &organizationsv1.FilterRequest_Ids{
							Ids: a.filters.Ids,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				testErr := status.Error(codes.NotFound, "not found")
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(nil, testErr)
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "internal error",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Ids: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Ids_{
						Ids: &organizationsv1.FilterRequest_Ids{
							Ids: a.filters.Ids,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				msg := diterrors.NewLocalizedError(diterrors.LocalizeLocale, testErr)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(nil, testErr)
				return nil, fmt.Errorf("can't filter organizations: %s", msg)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Ids: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
				options: portal.OrganizationsFilterOptions{
					WithLiquidated: true,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{WithLiquidated: true}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Ids_{
						Ids: &organizationsv1.FilterRequest_Ids{
							Ids: a.filters.Ids,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				response := &organizationsv1.FilterResponse{
					Organizations: []*organizationsv1.Organization{
						{
							Id: "1",
						},
						{
							Id: "2",
						},
						{
							Id: "3",
						},
					},
					Pagination: &sharedv1.PaginationResponse{
						Limit: 10,
					},
				}
				var organizations []*portal.Organization
				pagination := &entity.StringPagination{}
				f.mapper.EXPECT().PaginationToEntity(response.Pagination).Return(pagination)
				f.mapper.EXPECT().OrganizationsToEntity(response.Organizations).Return(organizations)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(response, nil)
				return &portal.OrganizationsWithPagination{
					Pagination:    pagination,
					Organizations: organizations,
				}, nil
			},
		},
		{
			name: "correct inns",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Inns: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Inns{
						Inns: &organizationsv1.FilterRequest_Inn{
							Inns: a.filters.Inns,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				response := &organizationsv1.FilterResponse{
					Organizations: []*organizationsv1.Organization{},
					Pagination: &sharedv1.PaginationResponse{
						Limit: 10,
					},
				}
				var organizations []*portal.Organization
				pagination := &entity.StringPagination{}
				f.mapper.EXPECT().PaginationToEntity(response.Pagination).Return(pagination)
				f.mapper.EXPECT().OrganizationsToEntity(response.Organizations).Return(organizations)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(response, nil)
				return &portal.OrganizationsWithPagination{
					Pagination:    pagination,
					Organizations: organizations,
				}, nil
			},
		},
		{
			name: "correct ogrns",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Ogrns: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Ogrns{
						Ogrns: &organizationsv1.FilterRequest_Ogrn{
							Ogrns: a.filters.Ogrns,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				response := &organizationsv1.FilterResponse{
					Organizations: []*organizationsv1.Organization{},
					Pagination: &sharedv1.PaginationResponse{
						Limit: 10,
					},
				}
				var organizations []*portal.Organization
				pagination := &entity.StringPagination{}
				f.mapper.EXPECT().PaginationToEntity(response.Pagination).Return(pagination)
				f.mapper.EXPECT().OrganizationsToEntity(response.Organizations).Return(organizations)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(response, nil)
				return &portal.OrganizationsWithPagination{
					Pagination:    pagination,
					Organizations: organizations,
				}, nil
			},
		},
		{
			name: "correct names",
			args: args{
				ctx: ctx,
				filters: portal.OrganizationsFilters{
					Names: []string{"1", "2", "3"},
				},
				pagination: &entity.StringPagination{
					Limit: &testInt,
				},
			},
			want: func(a args, f fields) (*portal.OrganizationsWithPagination, error) {
				limit := *a.pagination.GetLimit()
				paginationPb := &sharedv1.PaginationRequest{
					Limit: uint32(limit),
				}
				optionsPB := &organizationsv1.OrganizationFilterOptions{}
				filterPb := &organizationsv1.FilterRequest{
					Filter: &organizationsv1.FilterRequest_Names_{
						Names: &organizationsv1.FilterRequest_Names{
							Names: a.filters.Names,
						},
					},
					Pagination: paginationPb,
					Options:    optionsPB,
				}

				f.mapper.EXPECT().PaginationToPb(a.pagination).Return(paginationPb)
				f.mapper.EXPECT().OptionsToPb(a.options).Return(optionsPB)
				response := &organizationsv1.FilterResponse{
					Organizations: []*organizationsv1.Organization{},
					Pagination: &sharedv1.PaginationResponse{
						Limit: 10,
					},
				}
				var organizations []*portal.Organization
				pagination := &entity.StringPagination{}
				f.mapper.EXPECT().PaginationToEntity(response.Pagination).Return(pagination)
				f.mapper.EXPECT().OrganizationsToEntity(response.Organizations).Return(organizations)
				f.client.EXPECT().Filter(a.ctx, filterPb).Return(response, nil)
				return &portal.OrganizationsWithPagination{
					Pagination:    pagination,
					Organizations: organizations,
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				client: organizationsv1.NewMockOrganizationsAPIClient(ctrl),
				mapper: NewMockOrganizationsMapper(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			repo := NewOrganizationsRepository(f.client, f.mapper)
			got, err := repo.Filter(tt.args.ctx, tt.args.filters, tt.args.pagination, tt.args.options)

			if wantErr != nil {
				assert.Empty(t, got)
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}
		})
	}
}
