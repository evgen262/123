package auth

import (
	"testing"

	viewAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http/view/auth"
	entityAuth "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
	"github.com/stretchr/testify/assert"
)

func Test_authPresenter_AuthToView(t *testing.T) {
	type args struct {
		authInfo *entityAuth.Auth
	}

	tests := []struct {
		name string
		args args
		want *viewAuth.AuthResponse
	}{
		{
			name: "correct with active portal",
			args: args{
				authInfo: &entityAuth.Auth{
					PortalSession: "testPortalSID",
					Portals: []*entityAuth.Portal{
						{ID: 1, Name: "testName1", IsSelected: false},
						{ID: 2, Name: "testName2", IsSelected: true},
					},
				},
			},
			want: &viewAuth.AuthResponse{
				PortalSession: "testPortalSID",
				Portals: []*viewAuth.Portal{
					{ID: 1, Name: "testName1", IsActive: false},
					{ID: 2, Name: "testName2", IsActive: true},
				},
			},
		},
		{
			name: "correct without active portal",
			args: args{
				authInfo: &entityAuth.Auth{
					PortalSession: "testPortalSID",
					Portals: []*entityAuth.Portal{
						{ID: 1, Name: "testName1", IsSelected: false},
						{ID: 2, Name: "testName2", IsSelected: false},
					},
				},
			},
			want: &viewAuth.AuthResponse{
				PortalSession: "testPortalSID",
				Portals: []*viewAuth.Portal{
					{ID: 1, Name: "testName1", IsActive: true},
					{ID: 2, Name: "testName2", IsActive: false},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := NewAuthPresenter()
			got := ap.AuthToView(tt.args.authInfo)

			assert.Equal(t, tt.want, got)
		})
	}
}
