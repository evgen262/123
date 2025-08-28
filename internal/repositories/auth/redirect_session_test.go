package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	redirectsessionv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/authfacade/redirectsession/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	gomock "go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	entitySession "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_redirectSessionRepository_CreateSession(t *testing.T) {
	type fields struct {
		client *redirectsessionv1.MockRedirectSessionAPIClient
		mapper *MockRedirectSessionMapper
		logger *ditzap.MockLogger
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
				ctx:      context.Background(),
				userInfo: &entitySession.RedirectSessionUserInfo{},
			},
			want: func(a args, f fields) (string, error) {
				sessionID := "validSessionID"
				f.mapper.EXPECT().UserInfoToCreateRequestUserUnfoPb(a.userInfo).Return(&redirectsessionv1.CreateSessionRequest_UserInfo{})
				f.client.EXPECT().CreateSession(a.ctx, gomock.Any()).Return(&redirectsessionv1.CreateSessionResponse{SessionId: sessionID}, nil)

				return sessionID, nil
			},
		},
		{
			name: "nil userInfo",
			args: args{
				ctx:      context.Background(),
				userInfo: nil,
			},
			want: func(a args, f fields) (string, error) {
				return "", fmt.Errorf("redirectSessionRepository.Create: %w", ErrUserInfoRequired)
			},
		},
		{
			name: "validation error",
			args: args{
				ctx:      context.Background(),
				userInfo: &entitySession.RedirectSessionUserInfo{},
			},
			want: func(a args, f fields) (string, error) {
				err := diterrors.NewValidationError(errors.New("validation error"))
				f.mapper.EXPECT().UserInfoToCreateRequestUserUnfoPb(a.userInfo).Return(&redirectsessionv1.CreateSessionRequest_UserInfo{})
				f.client.EXPECT().CreateSession(a.ctx, gomock.Any()).Return(nil, err)
				f.logger.EXPECT().Warn(
					"invalid session", zap.Error(err),
					zap.String("method", "sr.client.CreateSession(ctx, session)"),
					zap.Any("session", a.userInfo),
				)

				return "", fmt.Errorf("client.CreateSession: invalid session: %w", err)
			},
		},
		{
			name: "other error",
			args: args{
				ctx:      context.Background(),
				userInfo: &entitySession.RedirectSessionUserInfo{
					// Populate with necessary fields
				},
			},
			want: func(a args, f fields) (string, error) {
				err := errors.New("some error")
				f.mapper.EXPECT().UserInfoToCreateRequestUserUnfoPb(a.userInfo).Return(&redirectsessionv1.CreateSessionRequest_UserInfo{})
				f.client.EXPECT().CreateSession(a.ctx, gomock.Any()).Return(nil, err)
				f.logger.EXPECT().Error(
					"can't create session", zap.Error(err),
					zap.String("method", "sr.client.CreateSession(ctx, session)"),
					zap.Any("session", a.userInfo),
				)

				return "", fmt.Errorf("client.CreateSession: can't create session: %w", err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				client: redirectsessionv1.NewMockRedirectSessionAPIClient(ctrl),
				mapper: NewMockRedirectSessionMapper(ctrl),
				logger: ditzap.NewMockLogger(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			r := NewRedirectSessionRepository(f.client, f.mapper, f.logger)
			got, err := r.CreateSession(tt.args.ctx, tt.args.userInfo)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
