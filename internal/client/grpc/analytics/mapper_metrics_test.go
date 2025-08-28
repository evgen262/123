package analytics

import (
	"testing"

	metricsv1 "git.mos.ru/buch-cloud/moscow-team-2.0/infrastructure/protolib.git/gen/infogorod/analytics/metrics/v1"
	"github.com/stretchr/testify/assert"

	entityAnalytics "git.mos.ru/buch-cloud/moscow-team-2.0/backend/web-api.git/internal/entity/analytics"
)

func Test_mapperAnalytics_MetricsCFCHeadersToPb(t *testing.T) {
	tests := []struct {
		name    string
		request entityAnalytics.CFCHeaders
		want    *metricsv1.AddRequest_CFCHeaders
	}{
		{
			name: "valid data",
			request: entityAnalytics.CFCHeaders{
				Header:        "test_header",
				Portal:        "test_portal",
				Authorization: "test_auth",
			},
			want: &metricsv1.AddRequest_CFCHeaders{
				CfcHeader:        "test_header",
				CfcPortal:        "test_portal",
				CfcAuthorization: "test_auth",
			},
		},
		{
			name: "empty data",
			request: entityAnalytics.CFCHeaders{
				Header:        "",
				Portal:        "",
				Authorization: "",
			},
			want: &metricsv1.AddRequest_CFCHeaders{
				CfcHeader:        "",
				CfcPortal:        "",
				CfcAuthorization: "",
			},
		},
	}

	m := NewMetricsMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := m.MetricsCFCHeadersToPb(tt.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
