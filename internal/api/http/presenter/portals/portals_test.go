package portals

import (
	"testing"

	"github.com/stretchr/testify/assert"

	viewPortals "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/portals"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/portal"
)

func TestPortalsPresenter_ToWebView(t *testing.T) {
	type args struct {
		portal *portal.Portal
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.WebPortal
	}{
		{
			name: "correct",
			args: args{
				portal: &portal.Portal{
					Id:        1,
					FullName:  "Test portal 1",
					ShortName: "test 1",
					Url:       "test1.mos.ru",
					LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
					ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
					Organizations: []*portal.PortalOrganization{
						{
							Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
						},
					},
					IsDeleted: false,
				},
			},
			want: &viewPortals.WebPortal{
				Name:  "Test portal 1",
				URL:   "test1.mos.ru",
				Image: "https://test1.mos.ru/path/to/logo.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPortalsPresenter()
			assert.Equalf(t, tt.want, p.ToWebView(tt.args.portal), "ToView(%v)", tt.args.portal)
		})
	}
}

func TestPortalsPresenter_ToWebViews(t *testing.T) {
	type args struct {
		portals []*portal.Portal
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.WebPortal
	}{
		{
			name: "correct",
			args: args{
				portals: []*portal.Portal{
					{
						Id:        1,
						FullName:  "Test portal 1",
						ShortName: "test 1",
						Url:       "test1.mos.ru",
						LogoUrl:   "https://test1.mos.ru/path/to/logo.jpg",
						ChatUrl:   "https://test1.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{
							{Id: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"},
						},
						IsDeleted: false,
					},
					{
						Id:        2,
						FullName:  "Test portal 2",
						ShortName: "test 2",
						Url:       "test2.mos.ru",
						LogoUrl:   "https://test2.mos.ru/path/to/logo.jpg",
						ChatUrl:   "https://test2.mos.ru/path/to/chat/utl/",
						Organizations: []*portal.PortalOrganization{
							{Id: "3c5cbb16-011a-310e-97e2-565400a26506"},
						},
						IsDeleted: false,
					},
				},
			},
			want: []*viewPortals.WebPortal{
				{
					Name:  "Test portal 1",
					URL:   "test1.mos.ru",
					Image: "https://test1.mos.ru/path/to/logo.jpg",
				},
				{
					Name:  "Test portal 2",
					URL:   "test2.mos.ru",
					Image: "https://test2.mos.ru/path/to/logo.jpg",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equalf(t, tt.want, p.ToWebViews(tt.args.portals), "ToViews(%v)", tt.args.portals)
		})
	}
}

func TestPortalsPresenter_organizationsToView(t *testing.T) {
	type args struct {
		organizations []*portal.PortalOrganization
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.OrganizationInfo
	}{
		{
			name: "correct",
			args: args{
				organizations: []*portal.PortalOrganization{
					{Id: "test", INN: "test"},
				},
			},
			want: []*viewPortals.OrganizationInfo{
				{ID: "test", INN: "test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.organizationsToView(tt.args.organizations))
		})
	}
}

func TestPortalsPresenter_ToView(t *testing.T) {
	type args struct {
		portal *portal.Portal
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.Portal
	}{
		{
			name: "correct",
			args: args{
				portal: &portal.Portal{Id: portal.PortalID(123)},
			},
			want: &viewPortals.Portal{
				Id:            123,
				Organizations: []*viewPortals.OrganizationInfo{},
				Active:        true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.ToView(tt.args.portal))
		})
	}
}

func TestPortalsPresenter_ToViews(t *testing.T) {
	type args struct {
		portals []*portal.Portal
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.Portal
	}{
		{
			name: "correct",
			args: args{
				portals: []*portal.Portal{
					{Id: portal.PortalID(123)},
				},
			},
			want: []*viewPortals.Portal{
				{
					Id:            123,
					Organizations: []*viewPortals.OrganizationInfo{},
					Active:        true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.ToViews(tt.args.portals))
		})
	}
}

func TestPortalsPresenter_ToShortView(t *testing.T) {
	type args struct {
		portal *portal.Portal
	}
	tests := []struct {
		name string
		args args
		want *viewPortals.PortalInfo
	}{
		{
			name: "correct",
			args: args{
				portal: &portal.Portal{FullName: "name"},
			},
			want: &viewPortals.PortalInfo{
				FullName:      "name",
				Organizations: []*viewPortals.OrganizationInfo{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.ToShortView(tt.args.portal))
		})
	}
}

func TestPortalsPresenter_ToShortViews(t *testing.T) {
	type args struct {
		portals []*portal.Portal
	}
	tests := []struct {
		name string
		args args
		want []*viewPortals.PortalInfo
	}{
		{
			name: "correct",
			args: args{
				portals: []*portal.Portal{
					{FullName: "name"},
				},
			},
			want: []*viewPortals.PortalInfo{
				{
					FullName:      "name",
					Organizations: []*viewPortals.OrganizationInfo{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.ToShortViews(tt.args.portals))
		})
	}
}

func TestPortalsPresenter_FilterOptionsToEntity(t *testing.T) {
	type args struct {
		options viewPortals.PortalsFilterOptions
	}
	tests := []struct {
		name string
		args args
		want portal.PortalsFilterOptions
	}{
		{
			name: "correct",
			args: args{
				options: viewPortals.PortalsFilterOptions{
					PortalIDs: []int{1, 2},
					OrgIDs:    []string{"id1", "id2"},
					INNs:      []string{"1", "2"},
				},
			},
			want: portal.PortalsFilterOptions{
				PortalIDs: portal.PortalIDs{1, 2},
				OrgIDs:    portal.OrganizationIDs{"id1", "id2"},
				INNs:      portal.OrganizationINNs{"1", "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := portalsPresenter{}
			assert.Equal(t, tt.want, p.FilterOptionsToEntity(tt.args.options))
		})
	}
}
