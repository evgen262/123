package portals

import (
	"testing"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"

	organizationsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/organizations/v1"
	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/portals/shared/v1"
	timeUtils "git.mos.ru/buch-cloud/moscow-team-2.0/build/time-utils.git"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_portalMapper_OnceGrbsToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		onceGrbsPb *organizationsv1.Organization_Grbs
	}

	testStartDate := wrapperspb.StringValue{
		Value: "Test StartDate",
	}
	testString := "Test StartDate"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portal.OrganizationGrbs
	}{{
		name: "correct",
		args: args{
			onceGrbsPb: &organizationsv1.Organization_Grbs{
				Name:      "Test Name 1",
				Inn:       "Test Inn 1",
				Ogrn:      "Test Ogrn 1",
				StartDate: &testStartDate,
			},
		},
		want: func(a args, f fields) *portal.OrganizationGrbs {
			return &portal.OrganizationGrbs{
				GrbsId:    "Test Ogrn 1_Test Inn 1",
				Name:      "Test Name 1",
				Inn:       "Test Inn 1",
				Ogrn:      "Test Ogrn 1",
				StartDate: &testString,
			}
		},
	},
		{
			name: "correct with startDate nil",
			args: args{
				onceGrbsPb: &organizationsv1.Organization_Grbs{
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: nil,
				},
			},
			want: func(a args, f fields) *portal.OrganizationGrbs {
				return &portal.OrganizationGrbs{
					GrbsId:    "Test Ogrn 1_Test Inn 1",
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: nil,
				}
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.OnceGrbsToEntity(tt.args.onceGrbsPb)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalMapper_GrbsToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		onceGrbsPb []*organizationsv1.Organization_Grbs
	}

	testStartDate := wrapperspb.StringValue{
		Value: "Test StartDate",
	}
	testString := "Test StartDate"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portal.OrganizationGrbs
	}{{
		name: "correct",
		args: args{
			onceGrbsPb: []*organizationsv1.Organization_Grbs{{
				Name:      "Test Name 1",
				Inn:       "Test Inn 1",
				Ogrn:      "Test Ogrn 1",
				StartDate: &testStartDate,
			}},
		},
		want: func(a args, f fields) []*portal.OrganizationGrbs {
			return []*portal.OrganizationGrbs{
				{
					GrbsId:    "Test Ogrn 1_Test Inn 1",
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: &testString,
				},
			}
		},
	},
		{
			name: "correct with startDate nil",
			args: args{
				onceGrbsPb: []*organizationsv1.Organization_Grbs{{
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: nil,
				}},
			},
			want: func(a args, f fields) []*portal.OrganizationGrbs {
				return []*portal.OrganizationGrbs{
					{
						GrbsId:    "Test Ogrn 1_Test Inn 1",
						Name:      "Test Name 1",
						Inn:       "Test Inn 1",
						Ogrn:      "Test Ogrn 1",
						StartDate: nil,
					},
				}
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.GrbsToEntity(tt.args.onceGrbsPb)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalMapper_OrganizationToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		organizationPb *organizationsv1.Organization
	}

	testStartDate := wrapperspb.StringValue{
		Value: "Test StartDate",
	}
	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testString := "Test StartDate"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) *portal.Organization
	}{{
		name: "correct",
		args: args{
			organizationPb: &organizationsv1.Organization{
				Id:              "Test Id",
				FullName:        "Test FullName",
				ShortName:       "Test ShortName",
				RegCode:         "Test RegCode",
				OrgCode:         "Test OrgCode",
				UnkCode:         "Test UnkCode",
				Ogrn:            "Test Ogrn",
				Inn:             "Test Inn",
				Kpp:             "Test Kpp",
				OrgTypeName:     "Test OrgTypeName",
				OrgTypeCode:     "Test OrgTypeCode",
				UchrezhTypeCode: "Test UchrezhTypeCode",
				Grbs: []*organizationsv1.Organization_Grbs{{
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: &testStartDate,
				}},
				Email:           "Test Email",
				Phone:           "Test Phone",
				AdditionalPhone: "Test AdditionalPhone",
				Site:            "Test Site",
				CreatedTime:     timestamppb.New(testT),
				UpdatedTime:     timestamppb.New(testT),
				LiquidatedAt:    &testStartDate,
				IsLiquidated:    true,
			},
		},
		want: func(a args, f fields) *portal.Organization {
			f.timeUtils.EXPECT().TimestampToTime(a.organizationPb.CreatedTime).Return(&testT)
			f.timeUtils.EXPECT().TimestampToTime(a.organizationPb.UpdatedTime).Return(&testT)
			return &portal.Organization{
				Id:              "Test Id",
				FullName:        "Test FullName",
				ShortName:       "Test ShortName",
				RegCode:         "Test RegCode",
				OrgCode:         "Test OrgCode",
				UNKCode:         "Test UnkCode",
				OGRN:            "Test Ogrn",
				INN:             "Test Inn",
				KPP:             "Test Kpp",
				OrgTypeName:     "Test OrgTypeName",
				OrgTypeCode:     "Test OrgTypeCode",
				UchrezhTypeCode: "Test UchrezhTypeCode",
				Grbs: []*portal.OrganizationGrbs{{
					GrbsId:    "Test Ogrn 1_Test Inn 1",
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: &testString,
				}},
				Email:           "Test Email",
				Phone:           "Test Phone",
				AdditionalPhone: "Test AdditionalPhone",
				Site:            "Test Site",
				CreatedAt:       &testT,
				UpdatedAt:       &testT,
				LiquidatedAt:    &testString,
				IsLiquidated:    true,
			}
		},
	},
		{
			name: "correct liquidation nil",
			args: args{
				organizationPb: &organizationsv1.Organization{
					Id:              "Test Id",
					FullName:        "Test FullName",
					ShortName:       "Test ShortName",
					RegCode:         "Test RegCode",
					OrgCode:         "Test OrgCode",
					UnkCode:         "Test UnkCode",
					Ogrn:            "Test Ogrn",
					Inn:             "Test Inn",
					Kpp:             "Test Kpp",
					OrgTypeName:     "Test OrgTypeName",
					OrgTypeCode:     "Test OrgTypeCode",
					UchrezhTypeCode: "Test UchrezhTypeCode",
					Grbs: []*organizationsv1.Organization_Grbs{{
						Name:      "Test Name 1",
						Inn:       "Test Inn 1",
						Ogrn:      "Test Ogrn 1",
						StartDate: &testStartDate,
					}},
					Email:           "Test Email",
					Phone:           "Test Phone",
					AdditionalPhone: "Test AdditionalPhone",
					Site:            "Test Site",
					CreatedTime:     timestamppb.New(testT),
					UpdatedTime:     timestamppb.New(testT),
				},
			},
			want: func(a args, f fields) *portal.Organization {
				f.timeUtils.EXPECT().TimestampToTime(a.organizationPb.CreatedTime).Return(&testT)
				f.timeUtils.EXPECT().TimestampToTime(a.organizationPb.UpdatedTime).Return(&testT)
				return &portal.Organization{
					Id:              "Test Id",
					FullName:        "Test FullName",
					ShortName:       "Test ShortName",
					RegCode:         "Test RegCode",
					OrgCode:         "Test OrgCode",
					UNKCode:         "Test UnkCode",
					OGRN:            "Test Ogrn",
					INN:             "Test Inn",
					KPP:             "Test Kpp",
					OrgTypeName:     "Test OrgTypeName",
					OrgTypeCode:     "Test OrgTypeCode",
					UchrezhTypeCode: "Test UchrezhTypeCode",
					Grbs: []*portal.OrganizationGrbs{{
						GrbsId:    "Test Ogrn 1_Test Inn 1",
						Name:      "Test Name 1",
						Inn:       "Test Inn 1",
						Ogrn:      "Test Ogrn 1",
						StartDate: &testString,
					}},
					Email:           "Test Email",
					Phone:           "Test Phone",
					AdditionalPhone: "Test AdditionalPhone",
					Site:            "Test Site",
					CreatedAt:       &testT,
					UpdatedAt:       &testT,
				}
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.OrganizationToEntity(tt.args.organizationPb)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalMapper_OrganizationsToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		organizationPb []*organizationsv1.Organization
	}

	testStartDate := wrapperspb.StringValue{
		Value: "Test StartDate",
	}
	testT := time.Date(2023, 10, 31, 13, 00, 00, 00, time.UTC)
	testString := "Test StartDate"

	tests := []struct {
		name string
		args args
		want func(a args, f fields) []*portal.Organization
	}{{
		name: "correct",
		args: args{
			organizationPb: []*organizationsv1.Organization{{
				Id:              "Test Id",
				FullName:        "Test FullName",
				ShortName:       "Test ShortName",
				RegCode:         "Test RegCode",
				OrgCode:         "Test OrgCode",
				UnkCode:         "Test UnkCode",
				Ogrn:            "Test Ogrn",
				Inn:             "Test Inn",
				Kpp:             "Test Kpp",
				OrgTypeName:     "Test OrgTypeName",
				OrgTypeCode:     "Test OrgTypeCode",
				UchrezhTypeCode: "Test UchrezhTypeCode",
				Grbs: []*organizationsv1.Organization_Grbs{{
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: &testStartDate,
				}},
				Email:           "Test Email",
				Phone:           "Test Phone",
				AdditionalPhone: "Test AdditionalPhone",
				Site:            "Test Site",
				CreatedTime:     timestamppb.New(testT),
				UpdatedTime:     timestamppb.New(testT),
			}},
		},
		want: func(a args, f fields) []*portal.Organization {
			f.timeUtils.EXPECT().TimestampToTime(a.organizationPb[0].CreatedTime).Return(&testT)
			f.timeUtils.EXPECT().TimestampToTime(a.organizationPb[0].UpdatedTime).Return(&testT)
			return []*portal.Organization{{
				Id:              "Test Id",
				FullName:        "Test FullName",
				ShortName:       "Test ShortName",
				RegCode:         "Test RegCode",
				OrgCode:         "Test OrgCode",
				UNKCode:         "Test UnkCode",
				OGRN:            "Test Ogrn",
				INN:             "Test Inn",
				KPP:             "Test Kpp",
				OrgTypeName:     "Test OrgTypeName",
				OrgTypeCode:     "Test OrgTypeCode",
				UchrezhTypeCode: "Test UchrezhTypeCode",
				Grbs: []*portal.OrganizationGrbs{{
					GrbsId:    "Test Ogrn 1_Test Inn 1",
					Name:      "Test Name 1",
					Inn:       "Test Inn 1",
					Ogrn:      "Test Ogrn 1",
					StartDate: &testString,
				}},
				Email:           "Test Email",
				Phone:           "Test Phone",
				AdditionalPhone: "Test AdditionalPhone",
				Site:            "Test Site",
				CreatedAt:       &testT,
				UpdatedAt:       &testT,
			}}
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.OrganizationsToEntity(tt.args.organizationPb)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalOrganizationsMapper_PaginationToEntity(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		paginationPb *sharedv1.PaginationResponse
	}
	testTime := time.Now()
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *entity.StringPagination
	}{
		{
			name: "correct",
			args: args{
				paginationPb: &sharedv1.PaginationResponse{
					Total:           1,
					Limit:           2,
					LastId:          "3",
					LastCreatedTime: timestamppb.New(testTime),
				},
			},
			want: func(a args, f fields) *entity.StringPagination {
				f.timeUtils.EXPECT().TimestampToTime(a.paginationPb.LastCreatedTime).Return(&testTime)
				total := 1
				limit := 2
				return &entity.StringPagination{
					Total:           &total,
					Limit:           &limit,
					LastId:          portal.OrganizationId("3"),
					LastCreatedTime: &testTime,
				}
			},
		},
		{
			name: "pagination nil",
			args: args{
				paginationPb: nil,
			},
			want: func(a args, f fields) *entity.StringPagination {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PaginationToEntity(tt.args.paginationPb)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalOrganizationsMapper_PaginationToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		pagination *entity.StringPagination
	}
	testInt := 1
	testTime := time.Now()
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *sharedv1.PaginationRequest
	}{
		{
			name: "correct",
			args: args{pagination: &entity.StringPagination{
				Limit:           &testInt,
				Total:           &testInt,
				LastId:          portal.OrganizationId("3"),
				LastCreatedTime: &testTime,
			}},
			want: func(a args, f fields) *sharedv1.PaginationRequest {
				f.timeUtils.EXPECT().TimeToTimestamp(a.pagination.GetLastCreatedTime()).Return(timestamppb.New(testTime))
				return &sharedv1.PaginationRequest{
					Limit:           uint32(*a.pagination.GetLimit()),
					LastId:          a.pagination.LastId.String(),
					LastCreatedTime: timestamppb.New(testTime),
				}
			},
		},
		{
			name: "pagination nil",
			args: args{pagination: nil},
			want: func(a args, f fields) *sharedv1.PaginationRequest {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.PaginationToPb(tt.args.pagination)

			assert.Equal(t, want, got)
		})
	}
}

func Test_portalOrganizationsMapper_OptionsToPb(t *testing.T) {
	type fields struct {
		timeUtils *timeUtils.MockTimeUtils
	}
	type args struct {
		options portal.OrganizationsFilterOptions
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) *organizationsv1.OrganizationFilterOptions
	}{
		{
			name: "true",
			args: args{options: portal.OrganizationsFilterOptions{
				WithLiquidated: true,
			}},
			want: func(a args, f fields) *organizationsv1.OrganizationFilterOptions {
				return &organizationsv1.OrganizationFilterOptions{WithLiquidated: true}
			},
		},
		{
			name: "false",
			args: args{options: portal.OrganizationsFilterOptions{}},
			want: func(a args, f fields) *organizationsv1.OrganizationFilterOptions {
				return &organizationsv1.OrganizationFilterOptions{WithLiquidated: false}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{timeUtils.NewMockTimeUtils(ctrl)}
			pfm := NewOrganizationsMapper(f.timeUtils)
			want := tt.want(tt.args, f)
			got := pfm.OptionsToPb(tt.args.options)

			assert.Equal(t, want, got)
		})
	}
}
