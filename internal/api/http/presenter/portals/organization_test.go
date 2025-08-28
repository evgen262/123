package portals

import (
	"testing"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view"
	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"

	"github.com/stretchr/testify/assert"
)

func Test_organizationPresenter_OrganizationIdsToEntity(t *testing.T) {
	type args struct {
		ids viewPortals.OrganizationIds
	}
	tests := []struct {
		name string
		args args
		want portal.OrganizationIDs
	}{
		{
			name: "correct",
			args: args{
				ids: viewPortals.OrganizationIds{"1", "2"},
			},
			want: portal.OrganizationIDs{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOrganizationPresenter()
			assert.Equalf(t, tt.want, op.OrganizationIdsToEntity(tt.args.ids), "OrganizationIdsToEntity(%v)", tt.args.ids)
		})
	}
}

func Test_organizationPresenter_FiltersToEntity(t *testing.T) {
	type args struct {
		filters viewPortals.OrganizationsFilters
	}
	tests := []struct {
		name string
		args args
		want portal.OrganizationsFilters
	}{
		{
			name: "correct",
			args: args{
				filters: viewPortals.OrganizationsFilters{
					Ids:   []string{"1", "2"},
					Names: []string{"3", "4"},
					Ogrns: []string{"5", "6"},
					Inns:  []string{"7", "8"},
				},
			},
			want: portal.OrganizationsFilters{
				Ids:   []string{"1", "2"},
				Names: []string{"3", "4"},
				Ogrns: []string{"5", "6"},
				Inns:  []string{"7", "8"},
			},
		},
		{
			name: "correct withLiquidated true",
			args: args{
				filters: viewPortals.OrganizationsFilters{
					Ids:   []string{"1", "2"},
					Names: []string{"3", "4"},
					Ogrns: []string{"5", "6"},
					Inns:  []string{"7", "8"},
				},
			},
			want: portal.OrganizationsFilters{
				Ids:   []string{"1", "2"},
				Names: []string{"3", "4"},
				Ogrns: []string{"5", "6"},
				Inns:  []string{"7", "8"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := organizationPresenter{}
			assert.Equalf(t, tt.want, op.FiltersToEntity(tt.args.filters), "FiltersToEntity(%v)", tt.args.filters)
		})
	}
}

func Test_organizationPresenter_OrganizationsToView(t *testing.T) {
	type args struct {
		orgs []*portal.Organization
	}
	testTime := time.Now()
	testString := "test"
	tests := []struct {
		name string
		args args
		want []*viewPortals.Organization
	}{
		{
			name: "correct",
			args: args{
				orgs: []*portal.Organization{
					{
						Id:              "1",
						FullName:        "2",
						ShortName:       "3",
						RegCode:         "4",
						OrgCode:         "5",
						UNKCode:         "6",
						OGRN:            "7",
						INN:             "8",
						KPP:             "9",
						OrgTypeName:     "10",
						OrgTypeCode:     "11",
						UchrezhTypeCode: "12",
						Grbs: []*portal.OrganizationGrbs{
							{
								GrbsId:    "13",
								Name:      "14",
								Inn:       "15",
								Ogrn:      "16",
								StartDate: &testString,
							},
						},
						Email:           "18",
						Phone:           "19",
						AdditionalPhone: "20",
						Site:            "21",
						CreatedAt:       &testTime,
						UpdatedAt:       &testTime,
						LiquidatedAt:    &testString,
						IsLiquidated:    true,
					},
				},
			},
			want: []*viewPortals.Organization{
				{
					Id:              "1",
					FullName:        "2",
					ShortName:       "3",
					RegCode:         "4",
					OrgCode:         "5",
					UNKCode:         "6",
					OGRN:            "7",
					INN:             "8",
					KPP:             "9",
					OrgTypeName:     "10",
					OrgTypeCode:     "11",
					UchrezhTypeCode: "12",
					Grbs: []*viewPortals.OrganizationGrbs{
						{
							GrbsId:    "13",
							Name:      "14",
							Inn:       "15",
							Ogrn:      "16",
							StartDate: &testString,
						},
					},
					Email:           "18",
					Phone:           "19",
					AdditionalPhone: "20",
					Site:            "21",
					CreatedAt:       &testTime,
					UpdatedAt:       &testTime,
					LiquidatedAt:    &testString,
					IsLiquidated:    true,
				},
			},
		},
		{
			name: "correct grbs start date and liquidation date nil",
			args: args{
				orgs: []*portal.Organization{
					{
						Id:              "1",
						FullName:        "2",
						ShortName:       "3",
						RegCode:         "4",
						OrgCode:         "5",
						UNKCode:         "6",
						OGRN:            "7",
						INN:             "8",
						KPP:             "9",
						OrgTypeName:     "10",
						OrgTypeCode:     "11",
						UchrezhTypeCode: "12",
						Grbs: []*portal.OrganizationGrbs{
							{
								GrbsId:    "13",
								Name:      "14",
								Inn:       "15",
								Ogrn:      "16",
								StartDate: nil,
							},
						},
						Email:           "18",
						Phone:           "19",
						AdditionalPhone: "20",
						Site:            "21",
						CreatedAt:       &testTime,
						UpdatedAt:       &testTime,
						LiquidatedAt:    nil,
					},
				},
			},
			want: []*viewPortals.Organization{
				{
					Id:              "1",
					FullName:        "2",
					ShortName:       "3",
					RegCode:         "4",
					OrgCode:         "5",
					UNKCode:         "6",
					OGRN:            "7",
					INN:             "8",
					KPP:             "9",
					OrgTypeName:     "10",
					OrgTypeCode:     "11",
					UchrezhTypeCode: "12",
					Grbs: []*viewPortals.OrganizationGrbs{
						{
							GrbsId:    "13",
							Name:      "14",
							Inn:       "15",
							Ogrn:      "16",
							StartDate: nil,
						},
					},
					Email:           "18",
					Phone:           "19",
					AdditionalPhone: "20",
					Site:            "21",
					CreatedAt:       &testTime,
					UpdatedAt:       &testTime,
					LiquidatedAt:    nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := organizationPresenter{}
			assert.Equalf(t, tt.want, op.OrganizationsToView(tt.args.orgs), "OrganizationsToView(%v)", tt.args.orgs)
		})
	}
}

func Test_organizationPresenter_PaginationToView(t *testing.T) {
	type args struct {
		pagination *entity.StringPagination
	}
	testInt := 1
	testTime := time.Now()
	tests := []struct {
		name string
		args args
		want *view.StringPagination
	}{
		{
			name: "correct",
			args: args{
				pagination: &entity.StringPagination{
					Total:           &testInt,
					Limit:           &testInt,
					LastId:          portal.OrganizationId("3"),
					LastCreatedTime: &testTime,
				},
			},
			want: &view.StringPagination{
				Total:           &testInt,
				Limit:           &testInt,
				LastId:          "3",
				LastCreatedTime: &testTime,
			},
		},
		{
			name: "last id nil",
			args: args{
				pagination: &entity.StringPagination{
					LastId:          nil,
					LastCreatedTime: &testTime,
				},
			},
			want: nil,
		},
		{
			name: "pagination nil",
			args: args{
				pagination: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := organizationPresenter{}
			assert.Equalf(t, tt.want, op.PaginationToView(tt.args.pagination), "PaginationToView(%v)", tt.args.pagination)
		})
	}
}

func Test_organizationPresenter_PaginationToEntity(t *testing.T) {
	type args struct {
		pagination *view.StringPagination
	}
	testInt := 1
	testTime := time.Now()
	tests := []struct {
		name string
		args args
		want *entity.StringPagination
	}{
		{
			name: "correct",
			args: args{
				pagination: &view.StringPagination{
					Limit:           &testInt,
					Total:           &testInt,
					LastId:          "3",
					LastCreatedTime: &testTime,
				},
			},
			want: &entity.StringPagination{
				Limit:           &testInt,
				Total:           &testInt,
				LastId:          portal.OrganizationId("3"),
				LastCreatedTime: &testTime,
			},
		},
		{
			name: "pagination nil",
			args: args{
				pagination: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := organizationPresenter{}
			assert.Equalf(t, tt.want, op.PaginationToEntity(tt.args.pagination), "PaginationToEntity(%v)", tt.args.pagination)
		})
	}
}

func Test_organizationPresenter_OptionsToEntity(t *testing.T) {
	type args struct {
		options *viewPortals.OrganizationsFilterOptions
	}
	tests := []struct {
		name string
		args args
		want portal.OrganizationsFilterOptions
	}{
		{
			name: "true",
			args: args{
				options: &viewPortals.OrganizationsFilterOptions{
					WithLiquidated: true,
				},
			},
			want: portal.OrganizationsFilterOptions{WithLiquidated: true},
		},
		{
			name: "nil input",
			args: args{
				options: nil,
			},
			want: portal.OrganizationsFilterOptions{WithLiquidated: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := organizationPresenter{}
			assert.Equalf(t, tt.want, op.OptionsToEntity(tt.args.options), "PaginationToEntity(%v)", tt.args.options)
		})
	}
}
