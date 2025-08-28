package portals

import (
	"testing"
	"time"

	portalsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/portals/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func TestPortalsMapper_NewPortalsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		portals []*portal.Portal
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portalsv1.AddRequest_Portal
	}{
		{
			name: "correct",
			args: args{
				portals: []*portal.Portal{{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							INN: "770123456789",
						},
					},
					CreatedAt: &testT,
					UpdatedAt: &testT,
					DeletedAt: &testT,
					IsDeleted: false,
				}},
			},
			want: func(a args, f fields) []*portalsv1.AddRequest_Portal {
				return []*portalsv1.AddRequest_Portal{{
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.NewPortalsToPb(tt.args.portals)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_NewPortalToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		portal *portal.Portal
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portalsv1.AddRequest_Portal
	}{
		{
			name: "correct",
			args: args{
				portal: &portal.Portal{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							INN: "770123456789",
						},
					},
					CreatedAt: &testT,
					UpdatedAt: &testT,
					DeletedAt: &testT,
					IsDeleted: false,
				},
			},
			want: func(a args, f fields) *portalsv1.AddRequest_Portal {
				return &portalsv1.AddRequest_Portal{
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.NewPortalToPb(tt.args.portal)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapperPortalsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		portals []*portal.Portal
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portalsv1.Portal
	}{
		{
			name: "correct",
			args: args{
				portals: []*portal.Portal{{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							INN: "770123456789",
						},
					},
					CreatedAt: &testT,
					UpdatedAt: &testT,
					DeletedAt: &testT,
					IsDeleted: false,
				}},
			},
			want: func(a args, f fields) []*portalsv1.Portal {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portals[0].CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portals[0].UpdatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portals[0].DeletedAt).Return(t)
				return []*portalsv1.Portal{{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Organizations: []*portalsv1.Portal_Organization{{
						OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:   "770123456789",
					}},
					Sort:        1,
					CreatedTime: t,
					UpdatedTime: t,
					DeletedTime: t,
					IsDeleted:   false,
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PortalsToPb(tt.args.portals)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_PortalToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		portal *portal.Portal
	}

	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portalsv1.Portal
	}{
		{
			name: "correct",
			args: args{
				portal: &portal.Portal{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Sort:      1,
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							INN: "770123456789",
						},
					},
					CreatedAt: &testT,
					UpdatedAt: &testT,
					DeletedAt: &testT,
					IsDeleted: false,
				},
			},
			want: func(a args, f fields) *portalsv1.Portal {
				t := timestamppb.New(testT)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portal.CreatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portal.UpdatedAt).Return(t)
				f.timeUtils.EXPECT().TimeToTimestamp(a.portal.DeletedAt).Return(t)
				return &portalsv1.Portal{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Organizations: []*portalsv1.Portal_Organization{{
						OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						Inn:   "770123456789",
					}},
					Sort:        1,
					CreatedTime: t,
					UpdatedTime: t,
					DeletedTime: t,
					IsDeleted:   false,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PortalToPb(tt.args.portal)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_PortalOrganizationsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		orgs []*portal.PortalOrganization
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portalsv1.Portal_Organization
	}{
		{
			name: "correct",
			args: args{
				orgs: []*portal.PortalOrganization{{
					Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					INN: "770123456789",
				}},
			},
			want: func(a args, f fields) []*portalsv1.Portal_Organization {
				return []*portalsv1.Portal_Organization{{
					OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					Inn:   "770123456789",
				}}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PortalOrganizationsToPb(tt.args.orgs)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapper_PortalOrganizationToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		org *portal.PortalOrganization
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portalsv1.Portal_Organization
	}{
		{
			name: "correct",
			args: args{
				org: &portal.PortalOrganization{
					Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					INN: "770123456789",
				},
			},
			want: func(a args, f fields) *portalsv1.Portal_Organization {
				return &portalsv1.Portal_Organization{
					OrgId: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
					Inn:   "770123456789",
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewPortalsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PortalOrganizationToPb(tt.args.org)

			assert.Equal(t, want, got)
		})
	}
}

func TestPortalsMapperPortalsToportals(t *testing.T) {
	type fields struct {
		TimeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		response *portalsv1.FilterResponse
	}
	tests := []struct {
		name string
		args args
		want []*portal.Portal
	}{
		{
			name: "nil",
			args: args{response: nil},
			want: []*portal.Portal{},
		},
		{
			name: "empty",
			args: args{response: &portalsv1.FilterResponse{Portals: nil}},
			want: []*portal.Portal{},
		},
		{
			name: "correct",
			args: args{
				response: &portalsv1.FilterResponse{
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
				},
			},
			want: []*portal.Portal{
				{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "https://test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
							INN: "770123456789",
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
					Organizations: []*portal.PortalOrganization{
						{
							Id:  "3c5cbb16-011a-310e-97e2-565400a26506",
							INN: "771234567890",
						},
					},
					IsDeleted: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{TimeUtils: timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewPortalsMapper(f.TimeUtils)
			assert.Equalf(t, tt.want, pm.PortalsToEntity(tt.args.response.GetPortals()), "ToEntities(%v)", tt.args.response)
		})
	}
}

func TestPortalsMapper_PortalToportals(t *testing.T) {
	type fields struct {
		TimeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		dto *portalsv1.Portal
	}
	tests := []struct {
		name string
		args args
		want *portal.Portal
	}{
		{
			name: "nil",
			args: args{
				dto: nil,
			},
			want: nil,
		},
		{
			name: "correct",
			args: args{
				dto: &portalsv1.Portal{
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
			},
			want: &portal.Portal{
				Id:        1,
				FullName:  "Test portal 1",
				ShortName: "test 1",
				Url:       "https://test1.mos.ru",
				LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
				ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
				Organizations: []*portal.PortalOrganization{
					{
						Id:  "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						INN: "770123456789",
					},
				},
				IsDeleted: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{TimeUtils: timeUtils.NewMockTimeUtils(ctrl)}
			pm := NewPortalsMapper(f.TimeUtils)
			assert.Equalf(t, tt.want, pm.PortalToEntity(tt.args.dto), "Toportals(%v)", tt.args.dto)
		})
	}
}
