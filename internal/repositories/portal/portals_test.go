package portal

import (
	"context"
	"fmt"
	"testing"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/repositories"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPortalsRepository_Add(t *testing.T) {
	type fields struct {
		portalsClient *portalsv1.MockPortalsAPIClient
		mapper        *MockPortalsMapper
	}
	type args struct {
		ctx     context.Context
		portals []*portal.Portal
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Portal, error)
	}{
		{
			name: "portalsAPIClient error",
			args: args{
				ctx:     ctx,
				portals: []*portal.Portal{{FullName: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsPb := []*portalsv1.AddRequest_Portal{{FullName: "test"}}
				f.mapper.EXPECT().NewPortalsToPb(a.portals).Return(portalsPb)
				f.portalsClient.EXPECT().
					Add(a.ctx, &portalsv1.AddRequest{Portals: portalsPb}).
					Return(nil, testErr)
				return nil, fmt.Errorf("can't add portal: %w", testErr)
			},
		},
		{
			name: "portalsAPIClient invalid argument error",
			args: args{
				ctx:     ctx,
				portals: []*portal.Portal{{FullName: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsPb := []*portalsv1.AddRequest_Portal{{FullName: "test"}}
				f.mapper.EXPECT().NewPortalsToPb(a.portals).Return(portalsPb)
				f.portalsClient.EXPECT().
					Add(a.ctx, &portalsv1.AddRequest{Portals: portalsPb}).
					Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "portalsAPIClient not found error",
			args: args{
				ctx:     ctx,
				portals: []*portal.Portal{{FullName: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsPb := []*portalsv1.AddRequest_Portal{{FullName: "test"}}
				f.mapper.EXPECT().NewPortalsToPb(a.portals).Return(portalsPb)
				f.portalsClient.EXPECT().
					Add(a.ctx, &portalsv1.AddRequest{Portals: portalsPb}).
					Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     ctx,
				portals: []*portal.Portal{{FullName: "test"}},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				portalsPb := []*portalsv1.AddRequest_Portal{{FullName: "test"}}
				NewPortalsPb := []*portalsv1.Portal{{Id: 1, FullName: "test"}}
				NewPortals := []*portal.Portal{{Id: 1, FullName: "test"}}
				f.mapper.EXPECT().NewPortalsToPb(a.portals).Return(portalsPb)
				f.portalsClient.EXPECT().
					Add(a.ctx, &portalsv1.AddRequest{Portals: portalsPb}).
					Return(&portalsv1.AddResponse{Portals: NewPortalsPb}, nil)
				f.mapper.EXPECT().PortalsToEntity(NewPortalsPb).Return(NewPortals)
				return NewPortals, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			portalsClient := portalsv1.NewMockPortalsAPIClient(ctrl)
			f := fields{
				portalsClient: portalsClient,
				mapper:        NewMockPortalsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewPortalsRepository(f.portalsClient, f.mapper)
			got, err := repo.Add(tt.args.ctx, tt.args.portals)
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

func TestPortalsRepository_Delete(t *testing.T) {
	type fields struct {
		portalsClient *portalsv1.MockPortalsAPIClient
		mapper        *MockPortalsMapper
	}
	type args struct {
		ctx context.Context
		id  portal.PortalID
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "portalsAPIClient error",
			args: args{
				ctx: ctx,
				id:  1,
			},
			want: func(a args, f fields) error {
				f.portalsClient.EXPECT().
					Delete(a.ctx, &portalsv1.DeleteRequest{PortalId: int32(a.id)}).
					Return(nil, testErr)
				return fmt.Errorf("can't delete portal with id[%d]: %w", a.id, testErr)
			},
		},
		{
			name: "portalsAPIClient invalid argument error",
			args: args{
				ctx: ctx,
				id:  1,
			},
			want: func(a args, f fields) error {
				f.portalsClient.EXPECT().
					Delete(a.ctx, &portalsv1.DeleteRequest{PortalId: int32(a.id)}).
					Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "portalsAPIClient not found error",
			args: args{
				ctx: ctx,
				id:  1,
			},
			want: func(a args, f fields) error {
				f.portalsClient.EXPECT().
					Delete(a.ctx, &portalsv1.DeleteRequest{PortalId: int32(a.id)}).
					Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return repositories.ErrNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				id:  1,
			},
			want: func(a args, f fields) error {
				f.portalsClient.EXPECT().
					Delete(a.ctx, &portalsv1.DeleteRequest{PortalId: int32(a.id)}).
					Return(&portalsv1.DeleteResponse{}, nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			portalsClient := portalsv1.NewMockPortalsAPIClient(ctrl)
			f := fields{
				portalsClient: portalsClient,
				mapper:        NewMockPortalsMapper(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			repo := NewPortalsRepository(f.portalsClient, f.mapper)
			err := repo.Delete(tt.args.ctx, tt.args.id)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPortalsRepository_Update(t *testing.T) {
	type fields struct {
		portalsClient *portalsv1.MockPortalsAPIClient
		mapper        *MockPortalsMapper
	}
	type args struct {
		ctx    context.Context
		portal *portal.Portal
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (*portal.Portal, error)
	}{
		{
			name: "portalsAPIClient error",
			args: args{
				ctx:    ctx,
				portal: &portal.Portal{Id: 1, FullName: "test"},
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				portalPb := &portalsv1.Portal{Id: 1, FullName: "test"}
				f.mapper.EXPECT().PortalToPb(a.portal).Return(portalPb)
				f.portalsClient.EXPECT().
					Update(a.ctx, &portalsv1.UpdateRequest{Portal: portalPb}).
					Return(nil, testErr)
				return nil, fmt.Errorf("can't update portal: %w", testErr)
			},
		},
		{
			name: "portalsAPIClient invalid argument error",
			args: args{
				ctx:    ctx,
				portal: &portal.Portal{Id: 1, FullName: "test"},
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				portalPb := &portalsv1.Portal{Id: 1, FullName: "test"}
				f.mapper.EXPECT().PortalToPb(a.portal).Return(portalPb)
				f.portalsClient.EXPECT().
					Update(a.ctx, &portalsv1.UpdateRequest{Portal: portalPb}).
					Return(nil, status.Error(codes.InvalidArgument, testErr.Error()))
				return nil, diterrors.NewValidationError(testErr)
			},
		},
		{
			name: "portalsAPIClient not found error",
			args: args{
				ctx:    ctx,
				portal: &portal.Portal{Id: 1, FullName: "test"},
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				portalPb := &portalsv1.Portal{Id: 1, FullName: "test"}
				f.mapper.EXPECT().PortalToPb(a.portal).Return(portalPb)
				f.portalsClient.EXPECT().
					Update(a.ctx, &portalsv1.UpdateRequest{Portal: portalPb}).
					Return(nil, status.Error(codes.NotFound, testErr.Error()))
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx:    ctx,
				portal: &portal.Portal{Id: 1, FullName: "test"},
			},
			want: func(a args, f fields) (*portal.Portal, error) {
				portalPb := &portalsv1.Portal{Id: 1, FullName: "test"}
				portal := &portal.Portal{Id: 1, FullName: "test"}
				f.mapper.EXPECT().PortalToPb(a.portal).Return(portalPb)
				f.portalsClient.EXPECT().
					Update(a.ctx, &portalsv1.UpdateRequest{Portal: portalPb}).
					Return(&portalsv1.UpdateResponse{Portal: portalPb}, nil)
				f.mapper.EXPECT().PortalToEntity(portalPb).Return(portal)
				return portal, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			portalsClient := portalsv1.NewMockPortalsAPIClient(ctrl)
			f := fields{
				portalsClient: portalsClient,
				mapper:        NewMockPortalsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewPortalsRepository(f.portalsClient, f.mapper)
			got, err := repo.Update(tt.args.ctx, tt.args.portal)
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

func TestPortalsRepository_Filter(t *testing.T) {
	type fields struct {
		portalsClient *portalsv1.MockPortalsAPIClient
		mapper        *MockPortalsMapper
	}
	type args struct {
		ctx     context.Context
		options portal.PortalsFilterOptions
	}
	ctx := context.TODO()
	tests := []struct {
		name string
		args args
		want func(a args, f fields) ([]*portal.Portal, error)
	}{
		{
			name: "err bad parameter",
			args: args{
				ctx: ctx,
				options: portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"771234567890abcdefg, 770987654321"},
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				req := &portalsv1.FilterRequest{
					Filters: &portalsv1.FilterRequest_Filters{
						Inns: a.options.INNs.ToStringSlice(),
					},
					Options: &portalsv1.FilterRequest_Options{
						OnlyLinked: a.options.OnlyLinked,
					},
				}
				testErr := status.Error(codes.InvalidArgument, "bad inn parameter")
				f.portalsClient.EXPECT().Filter(a.ctx, req).Return(nil, testErr)
				return nil, diterrors.NewValidationError(fmt.Errorf("bad inn parameter"))
			},
		},
		{
			name: "err not found",
			args: args{
				ctx: ctx,
				options: portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"771234567890, 770987654321"},
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				req := &portalsv1.FilterRequest{
					Filters: &portalsv1.FilterRequest_Filters{
						Inns: a.options.INNs.ToStringSlice(),
					},
					Options: &portalsv1.FilterRequest_Options{
						OnlyLinked: a.options.OnlyLinked,
					},
				}
				testErr := status.Error(codes.NotFound, "not found")
				f.portalsClient.EXPECT().Filter(a.ctx, req).Return(nil, testErr)
				return nil, repositories.ErrNotFound
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				options: portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"771234567890abcdefg, 770987654321"},
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				req := &portalsv1.FilterRequest{
					Filters: &portalsv1.FilterRequest_Filters{
						Inns: a.options.INNs.ToStringSlice(),
					},
					Options: &portalsv1.FilterRequest_Options{
						OnlyLinked: a.options.OnlyLinked,
					},
				}
				testErr := fmt.Errorf("some service error")
				f.portalsClient.EXPECT().Filter(a.ctx, req).Return(nil, testErr)
				return nil, fmt.Errorf("can't filter portal: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				options: portal.PortalsFilterOptions{
					INNs:       portal.OrganizationINNs{"771234567890abcdefg, 770987654321"},
					OnlyLinked: true,
				},
			},
			want: func(a args, f fields) ([]*portal.Portal, error) {
				req := &portalsv1.FilterRequest{
					Filters: &portalsv1.FilterRequest_Filters{
						Inns: a.options.INNs.ToStringSlice(),
					},
					Options: &portalsv1.FilterRequest_Options{
						OnlyLinked: a.options.OnlyLinked,
					},
				}
				resp := &portalsv1.FilterResponse{
					Portals: []*portalsv1.Portal{
						{
							Id:        1,
							FullName:  "Test portal 1",
							ShortName: "test 1",
							Url:       "https://test1.mos.ru",
							LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
							ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
							Organizations: []*portalsv1.Portal_Organization{
								{
									OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
									Inn:   "770123456789",
								},
							},
							IsDeleted: false,
						},
						{
							Id:        2,
							FullName:  "Test portal 2",
							ShortName: "test 2",
							Url:       "https://test2.mos.ru",
							LogoUrl:   "https://test2.mos.ru/path/to/logo.jpg",
							ChatUrl:   "https://test2.mos.ru/path/to/chat/utl/",
							Organizations: []*portalsv1.Portal_Organization{
								{
									OrgId: "3c5cbb16-011a-310e-97e2-565400a26506",
									Inn:   "771234567890",
								},
							},
							IsDeleted: false,
						},
					},
				}

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

				f.portalsClient.EXPECT().Filter(a.ctx, req).Return(resp, nil)
				f.mapper.EXPECT().PortalsToEntity(resp.GetPortals()).Return(portals)
				return portals, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			portalsClient := portalsv1.NewMockPortalsAPIClient(ctrl)
			f := fields{
				portalsClient: portalsClient,
				mapper:        NewMockPortalsMapper(ctrl),
			}
			want, wantErr := tt.want(tt.args, f)
			repo := NewPortalsRepository(f.portalsClient, f.mapper)
			got, err := repo.Filter(tt.args.ctx, tt.args.options)

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
