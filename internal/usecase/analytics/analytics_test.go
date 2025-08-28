package analytics

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity"
	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
	"git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/auth"
)

func Test_analyticsInteractor_AddMetrics(t *testing.T) {
	type fields struct {
		repository *MockAnalyticsRepository
		logger     *ditzap.MockLogger
	}
	type args struct {
		ctx     context.Context
		headers entityAnalytics.XCFCUserAgentHeader
		body    []byte
	}

	tests := []struct {
		name string
		args args
		want func(f fields, a args) (string, error)
	}{
		{
			name: "success",
			args: args{
				ctx: entity.WithSession(context.Background(), &auth.Session{
					ActivePortal: &auth.ActivePortal{Portal: auth.Portal{URL: "portalURL"}},
				}),
				headers: "DeviceID=3fee48ed-5394-4356-8bbe-21a4ac214681;DeviceType=web",
				body:    []byte(`{"key":"value"}`),
			},
			want: func(f fields, a args) (string, error) {
				headers := entityAnalytics.CFCHeaders{Header: a.headers}
				f.logger.EXPECT().Debug(
					"add metrics data",
					zap.String("call", "analyticsInteractor.AddMetrics"),
					zap.String("body", string(a.body)),
					zap.String("headers", string(a.headers)),
				)
				f.logger.EXPECT().Warn("empty user", zap.String("call", "session.GetUser"))
				f.repository.EXPECT().AddMetrics(a.ctx, headers, a.body).Return("metricID", nil)

				return "metricID", nil
			},
		},
		{
			name: "active portal is nil warn success",
			args: args{
				ctx:     entity.WithSession(context.Background(), &auth.Session{}),
				headers: "DeviceID=3fee48ed-5394-4356-8bbe-21a4ac214681;DeviceType=web",
				body:    []byte(`{"key":"value"}`),
			},
			want: func(f fields, a args) (string, error) {
				headers := entityAnalytics.CFCHeaders{Header: a.headers}
				f.logger.EXPECT().Debug(
					"add metrics data",
					zap.String("call", "analyticsInteractor.AddMetrics"),
					zap.String("body", string(a.body)),
					zap.String("headers", string(a.headers)),
				)
				f.logger.EXPECT().Warn("empty user", zap.String("call", "session.GetUser"))
				f.logger.EXPECT().Warn("empty active portal", zap.String("call", "session.GetActivePortal"))
				f.repository.EXPECT().AddMetrics(a.ctx, headers, a.body).Return("metricID", nil)

				return "metricID", nil
			},
		},
		{
			name: "empty body error",
			args: args{
				ctx:     context.Background(),
				headers: "DeviceID=;DeviceType=web",
				body:    []byte{},
			},
			want: func(f fields, a args) (string, error) {
				f.logger.EXPECT().Error("body is empty", zap.String("call", "analyticsInteractor.AddMetrics"))
				return "", diterrors.NewValidationError(ErrBodyIsEmpty, diterrors.ErrValidationFields{
					Field:   "body",
					Message: ErrBodyIsEmpty.Error(),
				})
			},
		},
		{
			name: "empty headers error",
			args: args{
				ctx:     context.Background(),
				headers: "",
				body:    []byte(`{"key":"value"}`),
			},
			want: func(f fields, a args) (string, error) {
				f.logger.EXPECT().Error("headers are empty", zap.String("call", "analyticsInteractor.AddMetrics"))
				return "", diterrors.NewValidationError(ErrHeadersAreAmpty, diterrors.ErrValidationFields{
					Field:   "headers",
					Message: ErrHeadersAreAmpty.Error(),
				})
			},
		},
		{
			name: "session error",
			args: args{
				ctx:     context.Background(),
				headers: "DeviceID=3fee48ed-5394-4356-8bbe-21a4ac214681;DeviceType=web",
				body:    []byte(`{"key":"value"}`),
			},
			want: func(f fields, a args) (string, error) {
				f.logger.EXPECT().Debug(
					"add metrics data",
					zap.String("call", "analyticsInteractor.AddMetrics"),
					zap.String("body", string(a.body)),
					zap.String("headers", string(a.headers)),
				)

				err := errors.New("session context not found")
				f.logger.EXPECT().Error("invalid session", zap.Error(err),
					zap.String("call", "entity.SessionFromContext"),
				)
				return "", fmt.Errorf("entity.SessionFromContext: %w", ErrUnauthenticated)
			},
		},
		{
			name: "repository error",
			args: args{
				ctx: entity.WithSession(context.Background(), &auth.Session{
					ActivePortal: &auth.ActivePortal{Portal: auth.Portal{URL: "portalURL"}},
				}),
				headers: "DeviceID=3fee48ed-5394-4356-8bbe-21a4ac214681;DeviceType=web",
				body:    []byte(`{"key":"value"}`),
			},
			want: func(f fields, a args) (string, error) {
				headers := entityAnalytics.CFCHeaders{Header: a.headers}
				f.logger.EXPECT().Debug(
					"add metrics data",
					zap.String("call", "analyticsInteractor.AddMetrics"),
					zap.String("body", string(a.body)),
					zap.String("headers", string(a.headers)),
				)
				err := errors.New("error")
				f.logger.EXPECT().Warn("empty user", zap.String("call", "session.GetUser"))
				f.repository.EXPECT().AddMetrics(a.ctx, headers, a.body).Return("", err)
				f.logger.EXPECT().Error("can't add metrics in repository", zap.Error(err),
					zap.String("call", "repository.AddMetrics"),
					zap.String("body", string(a.body)),
					zap.String("headers", string(a.headers)),
				)
				return "", fmt.Errorf("repository.AddMetrics: %w", err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				repository: NewMockAnalyticsRepository(ctrl),
				logger:     ditzap.NewMockLogger(ctrl),
			}

			uc := NewAnalyticsInteractor(f.repository, f.logger)
			want, wantErr := tt.want(f, tt.args)
			got, err := uc.AddMetrics(tt.args.ctx, tt.args.headers, tt.args.body)

			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Equal(t, want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
