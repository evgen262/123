package config

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
					"HTTP_HOST":           "192.168.1.2",
					"DEV_MODE":            "1234",
					"MAX_CPU":             "0",
					"SUDIR_CLIENT_URL":    "https://sudir-test.mos.ru",
					"SUDIR_CLIENT_ID":     "TestClient",
					"SUDIR_CLIENT_SECRET": "TestPassword",
					"SKS_URL":             "https://predprod-kadry2.mos.ru/hr5_rk",
					"SKS_SUBSCRIBER":      "Test",
					"SKS_USER":            "TestID",
					"SKS_SECRET":          "TestPassword",
					"REDIS_HOST":          "127.0.0.1",
					"REDIS_PORT":          "6379",
					"REDIS_DB":            "0",
					"EMPLOYEES_ENDPOINT":  "http://localhost:8080",
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
					"HTTP_HOST":           "192.168.1.2",
					"HTTP_PORT":           "8081",
					"DEV_MODE":            "true",
					"MAX_CPU":             "0",
					"SUDIR_CLIENT_URL":    "https://sudir-test.mos.ru",
					"SUDIR_CLIENT_ID":     "TestClient",
					"SUDIR_CLIENT_SECRET": "TestPassword",
					"SKS_URL":             "https://predprod-kadry2.mos.ru/hr5_rk",
					"SKS_SUBSCRIBER":      "Test",
					"SKS_USER":            "TestID",
					"SKS_SECRET":          "TestPassword",
					"REDIS_HOST":          "127.0.0.1",
					"REDIS_PORT":          "6379",
					"REDIS_DB":            "0",
					"EMPLOYEES_ENDPOINT":  "http://localhost:8080",
				}
			},
			want: func() (*Config, error) {
				return &Config{
					AppName:         "auth",
					DevMode:         true,
					LogLevel:        "warn",
					ServiceHTTPHost: ":8080",
					GrpcServer: struct {
						Host string `long:"grpc-host" description:"Listen grpc host" env:"GRPC_HOST" default:"0.0.0.0"`
						Port int    `long:"grpc-port" description:"Listen grpc port" env:"GRPC_PORT" required:"true" default:"9999"`
					}{
						Host: "0.0.0.0",
						Port: 9999,
					},
					OAuth: struct {
						URL          string `long:"sudir-url" env:"SUDIR_CLIENT_URL" required:"true"`
						ClientID     string `long:"sudir-id" env:"SUDIR_CLIENT_ID" required:"true"`
						ClientSecret string `long:"sudir-secret" env:"SUDIR_CLIENT_SECRET" required:"true"`
					}{
						URL:          "https://sudir-test.mos.ru",
						ClientID:     "TestClient",
						ClientSecret: "TestPassword",
					},
					SKS: struct {
						URL          string `long:"sks-url" env:"SKS_URL" required:"true"`
						SubscriberID string `long:"sks-subscriber" env:"SKS_SUBSCRIBER" required:"true"`
						UserID       string `long:"sks-user" env:"SKS_USER" required:"true"`
						Secret       string `long:"sks-secret" env:"SKS_SECRET" required:"true"`
					}{
						URL:          "https://predprod-kadry2.mos.ru/hr5_rk",
						SubscriberID: "Test",
						UserID:       "TestID",
						Secret:       "TestPassword",
					},
					Redis: struct {
						Host     string `long:"redis-host" env:"REDIS_HOST" required:"true"`
						Port     int    `long:"redis-port" env:"REDIS_PORT" required:"true"`
						Password string `long:"redis-pass" env:"REDIS_PASS" default:""`
						Username string `long:"redis-user" env:"REDIS_USER" default:""`
						DB       int    `long:"redis-db" env:"REDIS_DB" required:"true"`
						Prefix   string `long:"redis-prefix" env:"REDIS_PREFIX" default:""`
					}{
						Host: "127.0.0.1",
						Port: 6379,
						DB:   0,
					},
					Endpoints: Endpoints{
						EmployeesEndpoint: "http://localhost:8080",
					},
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
