package config

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
)

// Config Конфигурация приложения
type Config struct {
	// LogLevel уровень логирования
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn, info,debug" env:"LOG_LEVEL" default:"warn"`

	// AppName наименование сервиса
	AppName string `long:"appname" env:"APP_NAME" required:"true" default:"auth"`
	// Environment окружение
	Environment string `env:"ENVIRONMENT" description:"App environment" default:"develop"`
	// DevMode режим отладки
	DevMode bool `long:"dev-mode" env:"DEV_MODE" description:"Developer mode"`
	MaxCpu  int  `long:"max-cpu" env:"MAX_CPU" description:"Max cpu usage (GOMAXPROC)" default:"0"`
	// SentryDSN подключение к sentry
	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN"`

	// ServiceHTTPHost хост-адрес для входящих HTTP-подключений healthz, readyz
	ServiceHTTPHost string `long:"service-host" env:"SERVICE_HOST" description:"Host for HTTP-server (ex. 127.0.0.1:8080)" default:":8080"`

	// GrpcServer настройки gRPC-сервера
	GrpcServer struct {
		// Host Хост для Grpc-сервера
		Host string `long:"grpc-host" description:"Listen grpc host" env:"GRPC_HOST" default:"0.0.0.0"`
		// Port Порт для Grpc-сервера
		Port int `long:"grpc-port" description:"Listen grpc port" env:"GRPC_PORT" required:"true" default:"9999"`
	}
	// OAuth параметры клиента авторизации СУДИР
	OAuth struct {
		URL          string `long:"sudir-url" env:"SUDIR_CLIENT_URL" required:"true"`
		ClientID     string `long:"sudir-id" env:"SUDIR_CLIENT_ID" required:"true"`
		ClientSecret string `long:"sudir-secret" env:"SUDIR_CLIENT_SECRET" required:"true"`
	}
	// SKS параметры клиента системы кадров
	SKS struct {
		URL          string `long:"sks-url" env:"SKS_URL" required:"true"`
		SubscriberID string `long:"sks-subscriber" env:"SKS_SUBSCRIBER" required:"true"`
		UserID       string `long:"sks-user" env:"SKS_USER" required:"true"`
		Secret       string `long:"sks-secret" env:"SKS_SECRET" required:"true"`
	}
	Redis struct {
		Host     string `long:"redis-host" env:"REDIS_HOST" required:"true"`
		Port     int    `long:"redis-port" env:"REDIS_PORT" required:"true"`
		Password string `long:"redis-pass" env:"REDIS_PASS" default:""`
		Username string `long:"redis-user" env:"REDIS_USER" default:""`
		DB       int    `long:"redis-db" env:"REDIS_DB" required:"true"`
		Prefix   string `long:"redis-prefix" env:"REDIS_PREFIX" default:""`
	}

	// Endpoints энд-поинты клиентских сервисов
	Endpoints Endpoints
}

type Endpoints struct {
	EmployeesEndpoint string `long:"employees-endpoint" description:"Employees gRpc endpoint address" env:"EMPLOYEES_ENDPOINT" required:"true"`
}

// NewConfig ...
func NewConfig() (*Config, error) {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(log.Writer())
		return nil, fmt.Errorf("config parse failed: %v", err)
	}

	return &cfg, nil
}
