package portalsv2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	viewPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	entityPortalsV2 "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

func TestComplexesPresenter_ComplexesToView(t *testing.T) {
	type args struct {
		complexes []*entityPortalsV2.Complex
	}
	tests := []struct {
		name string
		args args
		want []*viewPortalsV2.Complex
	}{
		{
			name: "empty input slice",
			args: args{
				complexes: []*entityPortalsV2.Complex{},
			},
			want: []*viewPortalsV2.Complex{},
		},
		{
			name: "single complex without responsible and portals",
			args: args{
				complexes: []*entityPortalsV2.Complex{
					{
						ID:           1,
						Sort:         10,
						ComplexGroup: 0,
						Responsible:  nil,
						Portals:      []*entityPortalsV2.ComplexPortal{},
					},
				},
			},
			want: []*viewPortalsV2.Complex{
				{
					ID:              1,
					Sort:            10,
					Group:           0,
					FirstName:       "",
					LastName:        "",
					MiddleName:      "",
					IconID:          "",
					HeadDescription: "",
					PortalIDs:       []int{},
				},
			},
		},
		{
			name: "single complex with responsible (optional fields nil) and portals",
			args: args{
				complexes: []*entityPortalsV2.Complex{
					{
						ID:           2,
						Sort:         20,
						ComplexGroup: 0,
						Responsible: &entityPortalsV2.ComplexResponsible{
							FirstName:   "Ivan",
							LastName:    "Ivanov",
							MiddleName:  nil,
							ImageID:     nil,
							Description: "Head of Department",
						},
						Portals: []*entityPortalsV2.ComplexPortal{
							{ID: 101},
							{ID: 102},
						},
					},
				},
			},
			want: []*viewPortalsV2.Complex{
				{
					ID:              2,
					Sort:            20,
					Group:           0,
					FirstName:       "Ivan",
					LastName:        "Ivanov",
					MiddleName:      "",
					IconID:          "",
					HeadDescription: "Head of Department",
					PortalIDs:       []int{101, 102},
				},
			},
		},
		{
			name: "single complex with responsible (optional fields non-nil) and portals",
			args: args{
				complexes: []*entityPortalsV2.Complex{
					{
						ID:           3,
						Sort:         30,
						ComplexGroup: 0,
						Responsible: &entityPortalsV2.ComplexResponsible{
							FirstName:   "Petr",
							LastName:    "Petrov",
							MiddleName:  testPtr("Petrovich"),
							ImageID:     testPtr("photo-id-123"),
							Description: "Chief Specialist",
						},
						Portals: []*entityPortalsV2.ComplexPortal{
							{ID: 201},
						},
					},
				},
			},
			want: []*viewPortalsV2.Complex{
				{
					ID:              3,
					Sort:            30,
					Group:           0,
					FirstName:       "Petr",
					LastName:        "Petrov",
					MiddleName:      "Petrovich",
					IconID:          "photo-id-123",
					HeadDescription: "Chief Specialist",
					PortalIDs:       []int{201},
				},
			},
		},
		{
			name: "multiple complexes with mixed data",
			args: args{
				complexes: []*entityPortalsV2.Complex{
					{
						ID:           4,
						Sort:         40,
						ComplexGroup: 0,
						Responsible:  nil,
						Portals: []*entityPortalsV2.ComplexPortal{
							{ID: 301},
							{ID: 302},
							{ID: 303},
						},
					},
					{
						ID:           5,
						Sort:         50,
						ComplexGroup: 0,
						Responsible: &entityPortalsV2.ComplexResponsible{
							FirstName:   "Anna",
							LastName:    "Sidorova",
							MiddleName:  testPtr("Annovna"),
							ImageID:     nil,
							Description: "Manager",
						},
						Portals: []*entityPortalsV2.ComplexPortal{},
					},
				},
			},
			want: []*viewPortalsV2.Complex{
				{
					ID:              4,
					Sort:            40,
					Group:           0,
					FirstName:       "",
					LastName:        "",
					MiddleName:      "",
					IconID:          "",
					HeadDescription: "",
					PortalIDs:       []int{301, 302, 303},
				},
				{
					ID:              5,
					Sort:            50,
					Group:           0,
					FirstName:       "Anna",
					LastName:        "Sidorova",
					MiddleName:      "Annovna",
					IconID:          "",
					HeadDescription: "Manager",
					PortalIDs:       []int{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewComplexesPresenter()
			assert.Equalf(t, tt.want, p.ComplexesToView(tt.args.complexes), "ComplexesToView(%v)", tt.args.complexes)
		})
	}
}

func testPtr(s string) *string {
	return &s
}
