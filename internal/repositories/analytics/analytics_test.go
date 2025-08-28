package analytics

import (
	"context"
	"errors"
	"fmt"
	"testing"

	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

func Test_analyticsRepository_AddMetrics(t *testing.T) {
	type fields struct {
		client *metricsv1.MockMetricsAPIClient
		mapper *MockMetricsMapper
	}
	type args struct {
		ctx     context.Context
		headers entityAnalytics.CFCHeaders
		body    []byte
	}

	tests := []struct {
		name string
		args args
		want func(a args, f fields) (string, error)
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				headers: entityAnalytics.CFCHeaders{Header: "testHeader"},
				body:    []byte("testBody"),
			},
			want: func(a args, f fields) (string, error) {
				f.mapper.EXPECT().MetricsCFCHeadersToPb(a.headers).Return(&metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)})
				f.client.EXPECT().Add(a.ctx, &metricsv1.AddRequest{CfcHeaders: &metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)}, Body: a.body}).Return(&metricsv1.AddResponse{MetricId: "testMetricID"}, nil)
				return "testMetricID", nil
			},
		},
		{
			name: "invalid argument error",
			args: args{
				ctx:     context.Background(),
				headers: entityAnalytics.CFCHeaders{Header: "testHeader"},
				body:    []byte("testBody"),
			},
			want: func(a args, f fields) (string, error) {
				f.mapper.EXPECT().MetricsCFCHeadersToPb(a.headers).Return(&metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)})
				st := status.New(codes.InvalidArgument, "invalid argument")
				localizedError := diterrors.NewLocalizedError("ru-RU", st.Err())
				f.client.EXPECT().Add(a.ctx, &metricsv1.AddRequest{CfcHeaders: &metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)}, Body: a.body}).Return(nil, st.Err())
				return "", fmt.Errorf("client.Add: %w", diterrors.NewValidationError(localizedError, diterrors.ErrValidationFields{
					Field:   "in",
					Message: localizedError.Error(),
				}))
			},
		},
		{
			name: "default error",
			args: args{
				ctx:     context.Background(),
				headers: entityAnalytics.CFCHeaders{Header: "testHeader"},
				body:    []byte("testBody"),
			},
			want: func(a args, f fields) (string, error) {
				f.mapper.EXPECT().MetricsCFCHeadersToPb(a.headers).Return(&metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)})
				err := errors.New("some error")
				f.client.EXPECT().Add(a.ctx, &metricsv1.AddRequest{CfcHeaders: &metricsv1.AddRequest_CFCHeaders{CfcHeader: string(a.headers.Header)}, Body: a.body}).Return(nil, err)
				return "", fmt.Errorf("client.Add: %w", err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				client: metricsv1.NewMockMetricsAPIClient(ctrl),
				mapper: NewMockMetricsMapper(ctrl),
			}

			want, wantErr := tt.want(tt.args, f)

			r := NewAnalyticsRepository(f.client, f.mapper)
			got, err := r.AddMetrics(tt.args.ctx, tt.args.headers, tt.args.body)

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
