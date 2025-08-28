package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_redirectSessionInteractor_CreateSession(t *testing.T) {
	type fields struct {
		repository *MockRedirectSessionRepository
		authLink   string
	}
	type args struct {
		ctx      context.Context
		userInfo *entitySession.RedirectSessionUserInfo
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userInfo: &entitySession.RedirectSessionUserInfo{
					PortalURL: "portal.example.com",
				},
			},
			want: func(a args, f fields) (string, error) {
				expectedSessionID := "test-session-id"
				expectedRedirectURL := "https://portal.example.com" + f.authLink + "?state=" + expectedSessionID

				f.repository.EXPECT().CreateSession(a.ctx, a.userInfo).Return(expectedSessionID, nil)

				return expectedRedirectURL, nil
			},
		},
		{
			name: "error creating redirect",
			args: args{
				ctx: context.Background(),
				userInfo: &entitySession.RedirectSessionUserInfo{
					PortalURL: "portal.example.com",
				},
			},
			want: func(a args, f fields) (string, error) {
				err := errors.New("some error")
				f.repository.EXPECT().CreateSession(a.ctx, a.userInfo).Return("", err)

				return "", fmt.Errorf("repository.CreateSession: %w", err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repository: NewMockRedirectSessionRepository(ctrl),
				authLink:   "/auth/link",
			}

			want, wantErr := tt.want(tt.args, f)

			rsi := NewRedirectSessionInteractor(f.repository, f.authLink)
			got, err := rsi.CreateSession(tt.args.ctx, tt.args.userInfo)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Equal(t, "", got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
