package auth

import (
	"testing"

	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"
	"github.com/stretchr/testify/assert"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_sessionMapper_SessionToEntity(t *testing.T) {
	tests := []struct {
		name    string
		session *redirectsessionv1.UserInfo
		want    *entitySession.RedirectSessionUserInfo
	}{
		{
			name: "success",
			session: &redirectsessionv1.UserInfo{
				SessionId: "test-id",
				Email:     "test@example.com",
				Snils:     "123-456-789 00",
				PortalUrl: "http://portal.url",
				TargetUrl: "http://target.url",
				Useragent: "test-agent",
				Ip:        "127.0.0.1",
			},
			want: &entitySession.RedirectSessionUserInfo{
				SessionID: "test-id",
				Email:     "test@example.com",
				SNILS:     "123-456-789 00",
				PortalURL: "http://portal.url",
				TargetURL: "http://target.url",
				UserAgent: "test-agent",
				IP:        "127.0.0.1",
			},
		},
		{
			name:    "nil session",
			session: nil,
			want:    nil,
		},
	}

	m := NewRedirectSessionMapper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.UserInfoToEntity(tt.session)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sessionMapper_SessionToPb(t *testing.T) {
	tests := []struct {
		name    string
		session *entitySession.RedirectSessionUserInfo
		want    *redirectsessionv1.CreateSessionRequest_UserInfo
	}{
		{
			name: "success",
			session: &entitySession.RedirectSessionUserInfo{
				Email:     "test@example.com",
				SNILS:     "123-456-789 00",
				PortalURL: "http://portal.url",
				TargetURL: "http://target.url",
				UserAgent: "test-agent",
				IP:        "127.0.0.1",
			},
			want: &redirectsessionv1.CreateSessionRequest_UserInfo{
				Email:     "test@example.com",
				Snils:     "123-456-789 00",
				PortalUrl: "http://portal.url",
				TargetUrl: "http://target.url",
				Useragent: "test-agent",
				Ip:        "127.0.0.1",
			},
		},
		{
			name:    "nil session",
			session: nil,
			want:    nil,
		},
	}

	m := NewRedirectSessionMapper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.UserInfoToCreateRequestUserUnfoPb(tt.session)
			assert.Equal(t, tt.want, got)
		})
	}
}
