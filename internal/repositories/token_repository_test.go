package repositories

import (
	"context"
	"errors"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_tokenRepository_getKey(t *testing.T) {
	testBasePrefix := "cache1:"

	type fields struct {
		source     *MockCacheSource
		logger     *ditzap.MockLogger
		basePrefix string
	}
	type args struct {
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) string
	}{
		{
			name: "correct",
			args: args{cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f"},
			want: func(a args, f fields) string {
				return f.basePrefix + tokenPrefix + a.cloudID
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
				basePrefix: testBasePrefix,
			}
			want := tt.want(tt.args, f)

			or := NewTokenRepository(f.basePrefix, f.source, f.logger)
			got := or.getKey(tt.args.cloudID)

			assert.Equal(t, want, got)
		})
	}
}

func Test_tokenRepository_Save(t *testing.T) {
	testBasePrefix := "cache1:"
	testCtx := context.TODO()

	type fields struct {
		source     *MockCacheSource
		logger     *ditzap.MockLogger
		basePrefix string
	}
	type args struct {
		ctx     context.Context
		cloudID string
		token   string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) error
	}{
		{
			name: "error",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				token:   "WNjZilQWNRjK9yFMg3YGmSAp-KTJFYySKT64TEiKGuBX_kRaTRImpztssfRetWT81sl98llHXwpCl9C1HyHALg",
			},
			want: func(a args, f fields) error {
				testErr := errors.New("some save error")
				testKey := f.basePrefix + tokenPrefix + a.cloudID

				f.source.EXPECT().SetEx(a.ctx, testKey, a.token, tokenTtl).Return(testErr)
				f.logger.EXPECT().Error("не удалось сохранить refresh токен",
					zap.Error(testErr),
				)
				return testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
				token:   "WNjZilQWNRjK9yFMg3YGmSAp-KTJFYySKT64TEiKGuBX_kRaTRImpztssfRetWT81sl98llHXwpCl9C1HyHALg",
			},
			want: func(a args, f fields) error {
				testKey := f.basePrefix + tokenPrefix + a.cloudID

				f.source.EXPECT().SetEx(a.ctx, testKey, a.token, tokenTtl).Return(nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
				basePrefix: testBasePrefix,
			}
			or := NewTokenRepository(f.basePrefix, f.source, f.logger)

			wantErr := tt.want(tt.args, f)
			err := or.Save(tt.args.ctx, tt.args.cloudID, tt.args.token)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_tokenRepository_Get(t *testing.T) {
	testBasePrefix := "cache1:"
	testCtx := context.TODO()

	type fields struct {
		source     *MockCacheSource
		logger     *ditzap.MockLogger
		basePrefix string
	}
	type args struct {
		ctx     context.Context
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "err",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) (string, error) {
				testKey := f.basePrefix + tokenPrefix + a.cloudID
				testErr := errors.New("some get error")

				f.source.EXPECT().Get(a.ctx, testKey).Return("", testErr)
				f.logger.EXPECT().Error("не удалось получить refresh токен",
					zap.Error(testErr),
				)
				return "", testErr
			},
		},
		{
			name: "not found",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) (string, error) {
				testKey := f.basePrefix + tokenPrefix + a.cloudID

				f.source.EXPECT().Get(a.ctx, testKey).Return("", redis.Nil)
				return "", ErrNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) (string, error) {
				testKey := f.basePrefix + tokenPrefix + a.cloudID
				testToken := "WNjZilQWNRjK9yFMg3YGmSAp-KTJFYySKT64TEiKGuBX_kRaTRImpztssfRetWT81sl98llHXwpCl9C1HyHALg"

				f.source.EXPECT().Get(a.ctx, testKey).Return(testToken, nil)
				return testToken, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
				basePrefix: testBasePrefix,
			}
			or := NewTokenRepository(f.basePrefix, f.source, f.logger)

			want, wantErr := tt.want(tt.args, f)
			got, err := or.Get(tt.args.ctx, tt.args.cloudID)

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

func Test_tokenRepository_Delete(t *testing.T) {
	testBasePrefix := "cache1:"
	testCtx := context.TODO()

	type fields struct {
		source     *MockCacheSource
		logger     *ditzap.MockLogger
		basePrefix string
	}
	type args struct {
		ctx     context.Context
		cloudID string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f fields)
	}{
		{
			name: "err",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) {
				testKey := f.basePrefix + tokenPrefix + a.cloudID
				testErr := errors.New("some get error")

				f.source.EXPECT().Delete(a.ctx, testKey).Return(testErr)
				f.logger.EXPECT().Error("не удалось удалить refresh токен",
					zap.Error(testErr),
				)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:     testCtx,
				cloudID: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f fields) {
				testKey := f.basePrefix + tokenPrefix + a.cloudID
				f.source.EXPECT().Delete(a.ctx, testKey).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
				basePrefix: testBasePrefix,
			}
			or := NewTokenRepository(f.basePrefix, f.source, f.logger)

			tt.want(tt.args, f)
			or.Delete(tt.args.ctx, tt.args.cloudID)
		})
	}
}
