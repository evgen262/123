package app

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/cmd/auth/config"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/source/cache"
)

func Test_app_GracefulShutdown(t *testing.T) {
	type fields struct {
		config        *config.Config
		serviceServer *MockServer
		redisClient   *cache.MockRedis
		logger        *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.TODO()
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "http server shutdown error",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) error {
				c := context.WithoutCancel(a.ctx)
				srvErr := errors.New("http shutdown error")
				redisErr := errors.New("redis close error")

				f.serviceServer.EXPECT().Shutdown(c).Return(srvErr)
				f.redisClient.EXPECT().Close().Return(redisErr)
				return errors.Join(
					fmt.Errorf("can't shutdown service server: %w", srvErr),
					fmt.Errorf("can't close redis connection: %w", redisErr),
				)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) error {
				c := context.WithoutCancel(a.ctx)
				f.serviceServer.EXPECT().Shutdown(c).Return(nil)
				f.redisClient.EXPECT().Close().Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				serviceServer: NewMockServer(ctrl),
				redisClient:   cache.NewMockRedis(ctrl),
				logger:        ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			a := NewApp(f.config, nil, f.logger)
			a.serviceHTTPServer = f.serviceServer
			a.redisClient = f.redisClient

			err := a.GracefulShutdown(tt.args.ctx)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
