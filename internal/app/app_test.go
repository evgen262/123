package app

import (
	"context"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/cmd/web-api/config"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/api/http"
)

func Test_app_GracefulShutdown(t *testing.T) {
	type fields struct {
		config     *config.Config
		sqlmock    sqlmock.Sqlmock
		httpServer *http.MockServer
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.TODO()
	testErr := fmt.Errorf("testErr")
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
				ctx := context.WithoutCancel(a.ctx)
				f.httpServer.EXPECT().Shutdown(ctx).Return(testErr)
				f.sqlmock.ExpectClose()
				return fmt.Errorf("can't shutdown http-server: %w", testErr)
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
			},
			want: func(a args, f fields) error {
				ctx := context.WithoutCancel(a.ctx)
				f.httpServer.EXPECT().Shutdown(ctx).Return(nil)
				f.sqlmock.ExpectClose().WillReturnError(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			f := fields{
				sqlmock:    mock,
				httpServer: http.NewMockServer(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}
			wantErr := tt.want(tt.args, f)
			a := NewApp(f.config, nil, f.logger)
			a.httpServer = f.httpServer

			err = a.GracefulShutdown(tt.args.ctx)
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
