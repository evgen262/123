package portalsv2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	portalsv2View "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portalsv2"
	portalv2Entity "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portalv2"
)

func TestPortalsPresenter_PortalsWithCountToView(t *testing.T) {
	var (
		imageID1           = "image-id-1"
		imageID2           = "image-id-2"
		managerImageID1    = "manager-image-id-1"
		managerMiddleName1 = "Ivanovich"
	)

	portalWithCounts1 := &portalv2Entity.PortalWithCounts{
		Portal: &portalv2Entity.Portal{
			ID:        101,
			ShortName: "ShortName1",
			Name:      "Full Name 1",
			ImageID:   &imageID1,
			Manager: &portalv2Entity.PortalManager{
				FirstName:  "Petr",
				LastName:   "Petrov",
				MiddleName: &managerMiddleName1,
				ImageID:    &managerImageID1,
				Prosition:  "Director",
			},
		},
		EmployeesCount: 1500,
		OrgsCount:      15,
	}

	portalWithCounts2 := &portalv2Entity.PortalWithCounts{
		Portal: &portalv2Entity.Portal{
			ID:        102,
			ShortName: "ShortName2",
			Name:      "Full Name 2",
			ImageID:   &imageID2,
			Manager: &portalv2Entity.PortalManager{
				FirstName:  "Anna",
				LastName:   "Sidorova",
				MiddleName: nil,
				ImageID:    nil,
				Prosition:  "Manager",
			},
		},
		EmployeesCount: 500,
		OrgsCount:      5,
	}

	expectedPortal1 := &portalsv2View.Portal{
		ID:          101,
		Name:        "ShortName1",
		IconID:      imageID1,
		Description: "Full Name 1",
		Count: portalsv2View.Count{
			Employees: 1500,
			Podved:    15,
		},
		SctructureType: []string{"staffpositions", "management"},
		Head: &portalsv2View.Head{
			FirstName:   "Petr",
			LastName:    "Petrov",
			MiddleName:  managerMiddleName1,
			ImageID:     managerImageID1,
			Description: "Director",
		},
	}

	expectedPortal2 := &portalsv2View.Portal{
		ID:          102,
		Name:        "ShortName2",
		IconID:      imageID2,
		Description: "Full Name 2",
		Count: portalsv2View.Count{
			Employees: 500,
			Podved:    5,
		},
		SctructureType: []string{"staffpositions", "management"},
		Head: &portalsv2View.Head{
			FirstName:   "Anna",
			LastName:    "Sidorova",
			MiddleName:  "",
			ImageID:     "",
			Description: "Manager",
		},
	}

	type args struct {
		portalsWithCounts []*portalv2Entity.PortalWithCounts
	}
	tests := []struct {
		name string
		args args
		want []*portalsv2View.Portal
	}{
		{
			name: "empty slice",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{},
			},
			want: []*portalsv2View.Portal{},
		},
		{
			name: "slice with one valid item",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{portalWithCounts1},
			},
			want: []*portalsv2View.Portal{expectedPortal1},
		},
		{
			name: "slice with multiple valid items",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{portalWithCounts1, portalWithCounts2},
			},
			want: []*portalsv2View.Portal{expectedPortal1, expectedPortal2},
		},
		{
			name: "slice with nil item",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{nil, portalWithCounts1},
			},
			want: []*portalsv2View.Portal{expectedPortal1},
		},
		{
			name: "slice with item having nil Portal",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{{Portal: nil}, portalWithCounts1},
			},
			want: []*portalsv2View.Portal{expectedPortal1},
		},
		{
			name: "slice with mix of valid, nil, and nil-portal items",
			args: args{
				portalsWithCounts: []*portalv2Entity.PortalWithCounts{portalWithCounts1, nil, portalWithCounts2, {Portal: nil}},
			},
			want: []*portalsv2View.Portal{expectedPortal1, expectedPortal2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPortalsPresenter()

			actual := p.PortalsWithCountToView(tt.args.portalsWithCounts)

			assert.Equalf(t, tt.want, actual, "PortalsWithCountToView(%v)", tt.args.portalsWithCounts)
		})
	}
}

func TestPortalsPresenter_PortalsFilterToEntity(t *testing.T) {
	type args struct {
		filter *portalsv2View.PortalsFilterRequest
	}
	tests := []struct {
		name string
		args args
		want *portalv2Entity.FilterPortalsFilters
	}{
		{
			name: "correct with IDs",
			args: args{
				filter: &portalsv2View.PortalsFilterRequest{
					PortalIDs: []int{1, 5, 10},
				},
			},
			want: &portalv2Entity.FilterPortalsFilters{
				IDs: []int{1, 5, 10},
			},
		},
		{
			name: "empty IDs slice",
			args: args{
				filter: &portalsv2View.PortalsFilterRequest{
					PortalIDs: []int{},
				},
			},
			want: &portalv2Entity.FilterPortalsFilters{
				IDs: []int{},
			},
		},
		{
			name: "nil IDs slice",
			args: args{
				filter: &portalsv2View.PortalsFilterRequest{
					PortalIDs: nil,
				},
			},
			want: &portalv2Entity.FilterPortalsFilters{
				IDs: []int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPortalsPresenter()
			assert.Equalf(t, tt.want, p.PortalsFilterToEntity(tt.args.filter), "PortalsFilterToEntity(%v)", tt.args.filter)
		})
	}
}
