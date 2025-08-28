package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/entity"
)

func Test_stateRepository_New(t *testing.T) {
	ctx := context.TODO()
	testErr := errors.New("some source error")
	type fields struct {
		basePrefix string
		source     *MockCacheSource
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		options *entity.StateOptions
	}
	tests := []struct {
		name string
		args args
		want func(a args, f *fields) (*entity.State, error)
	}{
		{
			name: "source err",
			args: args{ctx: ctx},
			want: func(a args, f *fields) (*entity.State, error) {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)
				uuid.SetClockSequence(0)

				testState := &entity.State{
					ID: "31313131-3131-4131-b131-313131313131",
				}
				testKey := f.basePrefix + statePrefix + testState.ID

				data, err := json.Marshal(testState)
				assert.NoError(t, err)

				f.source.EXPECT().SetEx(a.ctx, testKey, data, stateTTL).Return(testErr)
				f.logger.EXPECT().Error("не удалось сохранить state",
					zap.Error(testErr),
				)

				return nil, testErr
			},
		},
		{
			name: "correct",
			args: args{
				ctx: ctx,
				options: &entity.StateOptions{
					ClientID:     "TestClientID",
					ClientSecret: "TestClientSecret",
					CodeVerifier: "Some generated code verifier",
					CallbackURL:  "http://test.com",
				},
			},
			want: func(a args, f *fields) (*entity.State, error) {
				reader := bytes.NewReader([]byte("1111111111111111"))
				uuid.SetRand(reader)
				uuid.SetClockSequence(0)

				testState := &entity.State{
					ID:           "31313131-3131-4131-b131-313131313131",
					ClientID:     "TestClientID",
					ClientSecret: "TestClientSecret",
					CodeVerifier: "Some generated code verifier",
					CallbackURL:  "http://test.com",
				}
				testKey := f.basePrefix + statePrefix + testState.ID

				data, err := json.Marshal(testState)
				assert.NoError(t, err)

				f.source.EXPECT().SetEx(a.ctx, testKey, data, stateTTL).Return(nil)

				return testState, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				basePrefix: "test:",
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}
			sr := NewStateRepository(f.basePrefix, f.source, f.logger)

			want, wantErr := tt.want(tt.args, f)
			got, err := sr.New(tt.args.ctx, tt.args.options)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Empty(t, got)
			} else {
				assert.Equal(t, want, got)
				assert.NoError(t, err)
			}

		})
	}
}

func Test_stateRepository_IsValid(t *testing.T) {
	ctx := context.TODO()
	testErr := errors.New("some test error")
	type fields struct {
		basePrefix string
		source     *MockCacheSource
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx   context.Context
		state string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f *fields) error
	}{
		{
			name: "parse error",
			args: args{
				ctx:   ctx,
				state: "wrong_state",
			},
			want: func(a args, f *fields) error {
				return ErrInvalidState
			},
		},
		{
			name: "source error",
			args: args{
				ctx:   ctx,
				state: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f *fields) error {
				testKey := f.basePrefix + statePrefix + a.state
				f.source.EXPECT().Exists(a.ctx, testKey).Return(false, testErr)
				f.logger.EXPECT().Error("не удалось получить state",
					zap.Error(testErr),
				)
				return testErr
			},
		},
		{
			name: "state not found",
			args: args{
				ctx:   ctx,
				state: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f *fields) error {
				testKey := f.basePrefix + statePrefix + a.state
				f.source.EXPECT().Exists(a.ctx, testKey).Return(false, nil)
				return ErrStateNotFound
			},
		},
		{
			name: "correct",
			args: args{
				ctx:   ctx,
				state: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f *fields) error {
				testKey := f.basePrefix + statePrefix + a.state
				f.source.EXPECT().Exists(a.ctx, testKey).Return(true, nil)
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				basePrefix: "test:",
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}
			sr := NewStateRepository(f.basePrefix, f.source, f.logger)

			want := tt.want(tt.args, f)
			got := sr.IsExists(tt.args.ctx, tt.args.state)

			if want != nil {
				assert.EqualError(t, got, want.Error())
			} else {
				assert.NoError(t, got)
			}
		})
	}
}

func Test_stateRepository_Delete(t *testing.T) {
	ctx := context.TODO()
	testErr := errors.New("some test error")
	type fields struct {
		basePrefix string
		source     *MockCacheSource
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx   context.Context
		state string
	}
	tests := []struct {
		name string
		args args
		want func(a args, f *fields)
	}{
		{
			name: "source err",
			args: args{
				ctx:   ctx,
				state: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f *fields) {
				testKey := f.basePrefix + statePrefix + a.state
				f.source.EXPECT().Delete(a.ctx, testKey).Return(testErr)
				f.logger.EXPECT().Error("не удалось удалить state",
					zap.Error(testErr),
				)
			},
		},
		{
			name: "correct",
			args: args{
				ctx:   ctx,
				state: "342a2c0c-d9ef-4cd6-b328-b67d9baf6a7f",
			},
			want: func(a args, f *fields) {
				testKey := f.basePrefix + statePrefix + a.state
				f.source.EXPECT().Delete(a.ctx, testKey).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := &fields{
				basePrefix: "test:",
				source:     NewMockCacheSource(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}
			sr := NewStateRepository(f.basePrefix, f.source, f.logger)

			tt.want(tt.args, f)
			sr.Delete(tt.args.ctx, tt.args.state)
		})
	}
}
