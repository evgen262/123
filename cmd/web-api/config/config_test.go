package config

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		prepare func() map[string]string
		want    func() (*Config, error)
	}{
		{
			name: "parse error",
			prepare: func() map[string]string {
				return map[string]string{
					"POSTGRES_DBNAME":           "public",
					"DEV_MODE":                  "1234",
					"SYSAPIKEY":                 "exampleKey",
					"AUTH_FACADE_ENDPOINT":      "exampleEndpoint",
					"ANALYTICS_ENDPOINT":        "exampleEndpoint2",
					"NOTIFICATIONS_ENDPOINT":    "exampleEndpoint3",
					"S3_BUCKET":                 "exampleBucket",
					"S3_ENDPOINT":               "exampleS3Endpoint",
					"S3_ACCESS_KEY_ID":          "exampleKeyId",
					"S3_SECRET_ACCESS_KEY":      "exampleAccKey",
					"S3_USE_SSL":                "true",
					"UPLOAD_PATH":               "examplePath",
					"PORTALS_ENDPOINT":          "exampleEndpoint4",
					"SURVEYS_ENDPOINT":          "exampleEndpoint5",
					"PROXY_FACADE_ENDPOINT":     "exampleEndpoint6",
					"FILES_ENDPOINT":            "exampleEndpoint7",
					"WEB_AUTH_URL":              "http://localhost/auth",
					"HTTP_EXTERNAL_HOST":        "localhostTest",
					"WEB_AUTH_REDIRECT_URI":     "test",
					"EMPLOYEES_SEARCH_ENDPOINT": "exampleEndpoint6",
					"EMPLOYEES_ENDPOINT":        "exampleEndpoint7",
					"PORTALS_FACADE_ENDPOINT":   "exampleEndpoint9",
					"NEWS_ENDPOINT":             "exampleEndpoint10",
					"PORTALSV2_ENDPOINT":        "exampleEndpoint10",
					"BANNERS_ENDPOINT":          "exampleEndpoint11",
					"NEWS_FACADE_ENDPOINT":      "exampleEndpoint12",
				}
			},
			want: func() (*Config, error) {
				prefix := "--"
				if runtime.GOOS == "windows" {
					prefix = "/"
				}
				return nil, fmt.Errorf(`config parse failed: invalid argument for flag `+"`"+`%s%s' (`+
					`expected bool): strconv.ParseBool: parsing "%s": invalid syntax`, prefix, "dev-mode", "1234")
			},
		},
		{
			name: "correct",
			prepare: func() map[string]string {
				return map[string]string{
					"POSTGRES_DBNAME":           "public",
					"DEV_MODE":                  "true",
					"SYSAPIKEY":                 "exampleKey",
					"AUTH_FACADE_ENDPOINT":      "exampleEndpoint",
					"ANALYTICS_ENDPOINT":        "exampleEndpoint2",
					"NOTIFICATIONS_ENDPOINT":    "exampleEndpoint3",
					"HTTP_EXTERNAL_HOST":        "localhostTest",
					"HTTP_SCHEMA":               "https",
					"UPLOAD_PATH":               "examplePath",
					"PORTALS_ENDPOINT":          "exampleEndpoint4",
					"SURVEYS_ENDPOINT":          "exampleEndpoint5",
					"PROXY_FACADE_ENDPOINT":     "exampleEndpoint6",
					"FILES_ENDPOINT":            "exampleEndpoint7",
					"WEB_AUTH_URL":              "http://localhost/auth",
					"WEB_AUTH_REDIRECT_URI":     "test",
					"EMPLOYEES_SEARCH_ENDPOINT": "exampleEndpoint7",
					"EMPLOYEES_ENDPOINT":        "exampleEndpoint8",
					"PORTALS_FACADE_ENDPOINT":   "exampleEndpoint9",
					"NEWS_ENDPOINT":             "exampleEndpoint10",
					"PORTALSV2_ENDPOINT":        "exampleEndpoint10",
					"BANNERS_ENDPOINT":          "exampleEndpoint12",
					"NEWS_FACADE_ENDPOINT":      "exampleEndpoint12"}
			},
			want: func() (*Config, error) {
				return &Config{
					AppName:         "web-api",
					DevMode:         true,
					LogLevel:        "warn",
					ApiKey:          "exampleApiKey",
					WebAuthURL:      "http://localhost/auth",
					ServiceHTTPHost: ":8080",
					HttpServer: &HttpServerConfig{
						Host:         "0.0.0.0",
						Port:         80,
						ExternalHost: "localhostTest",
						Schema:       "https",
						AllowOrigins: "*",
					},
					Path: struct {
						UploadPath string `long:"upload-path" description:"upload path" env:"UPLOAD_PATH" required:"true"`
					}{
						UploadPath: "examplePath",
					},
					Endpoints: &Endpoints{
						PortalsEndpoint:         "exampleEndpoint4",
						SurveysEndpoint:         "exampleEndpoint5",
						AuthFacadeEndpoint:      "exampleEndpoint",
						ProxyFacadeEndpoint:     "exampleEndpoint6",
						EmployeesSearchEndpoint: "exampleEndpoint7",
						FilesEndpoint:           "exampleEndpoint7",
						EmployeesEndpoint:       "exampleEndpoint8",
						AnalyticsEndpoint:       "exampleEndpoint2",
						PortalsFacadeEndpoint:   "exampleEndpoint9",
						NewsEndpoint:            "exampleEndpoint10",
						PortalsV2Endpoint:       "exampleEndpoint10",
						BannersEndpoint:         "exampleEndpoint12",
						NewsFacadeEndpoint:      "exampleEndpoint12",
					},
					TTL: &TTL{
						AccessToken:  time.Duration(7776000000000000),
						RefreshToken: time.Duration(7776000000000000),
					},
					WebAuthRedirectURI: "test",
				}, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envs := tt.prepare()
			for k, v := range envs {
				t.Setenv(k, v)
			}

			got, err := NewConfig()
			want, wantErr := tt.want()
			if wantErr != nil {
				assert.EqualError(t, err, wantErr.Error())
				assert.Nil(t, want)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}
}
