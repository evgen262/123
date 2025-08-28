package files

import (
	"testing"

	sharedv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/fileservice/shared/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_visitorMapper_SessionToVisitorPb(t *testing.T) {
	tests := []struct {
		name    string
		session *auth.Session
		want    *sharedv1.Visitor
	}{
		{
			name: "Authenticated user",
			session: &auth.Session{
				UserAuthType: auth.UserAuthTypeAuth,
				User: &auth.User{
					ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
				},
			},
			want: &sharedv1.Visitor{
				Visitor: &sharedv1.Visitor_User_{
					User: &sharedv1.Visitor_User{
						Id: "22222222-2222-2222-2222-222222222222",
					},
				},
			},
		},
		{
			name: "Old Authenticated user",
			session: &auth.Session{
				UserAuthType: auth.UserAuthTypeOldAuth,
				User: &auth.User{
					ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"),
				},
			},
			want: &sharedv1.Visitor{
				Visitor: &sharedv1.Visitor_User_{
					User: &sharedv1.Visitor_User{
						Id: "33333333-3333-3333-3333-333333333333",
					},
				},
			},
		},
		{
			name: "Unknown user type",
			session: &auth.Session{
				UserAuthType: auth.UserAuthType(-1),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := NewVisitorMapper()
			got := mapper.SessionToVisitorPb(tt.session)
			assert.Equal(t, tt.want, got)
		})
	}
}
