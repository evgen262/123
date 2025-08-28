package config

import (
	"fmt"
	"log"
	"time"

	"github.com/jessevdk/go-flags"
)

// Config Конфигурация приложения
type Config struct {
	// LogLevel уровень логирования
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn, info,debug" env:"LOG_LEVEL" default:"warn"`

	// AppName наименование сервиса
	AppName string `long:"appname" env:"APP_NAME" default:"web-api"`
	// Environment окружение
	Environment string `env:"ENVIRONMENT" description:"App environment (develop, prod)" default:"develop"`
	// DevMode режим отладки
	DevMode bool `long:"dev-mode" env:"DEV_MODE" description:"Developer mode"`
	MaxCpu  int  `long:"max-cpu" env:"MAX_CPU" description:"Max cpu usage (GOMAXPROC)" default:"0"`
	// ApiKey апи-ключ для доступа к служебным API
	ApiKey    string `long:"api-key" env:"APIKEY" default:"exampleApiKey"`
	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN"`

	Path struct {
		UploadPath string `long:"upload-path" description:"upload path" env:"UPLOAD_PATH" required:"true"`
	}

	WebAuthURL string `long:"auth-redirect-url" description:"WEB portal url for auth with scheme" env:"WEB_AUTH_URL" required:"true"`

	// ServiceHTTPHost хост-адрес для входящих HTTP-подключений healthz, readyz, info
	ServiceHTTPHost string `long:"service-host" env:"SERVICE_HOST" description:"Host for service HTTP-server (ex. 127.0.0.1:8080)" default:":8080"`
	// HttpServer настройки Http-сервера
	HttpServer *HttpServerConfig

	// Endpoints энд-поинты клиентских сервисов
	Endpoints *Endpoints

	// TTL настройки времи действия
	TTL *TTL

	WebAuthRedirectURI string `long:"web-auth-redirect-uri" description:"Old portal handler uri for short redirect session" env:"WEB_AUTH_REDIRECT_URI" required:"true"`

	AccessListFile string `long:"access-list-file" env:"ACCESS_LIST_FILE"`
}

type HttpServerConfig struct {
	// Host Хост для Http-сервера
	Host string `long:"http-host" description:"Listen http host" env:"HTTP_HOST" default:"0.0.0.0"`
	// Port Порт для Http-сервера
	Port         int    `long:"http-port" description:"Listen http port" env:"HTTP_PORT" default:"80"`
	ExternalHost string `long:"http-external-host" description:"External host for Http-server" env:"HTTP_EXTERNAL_HOST" required:"true"`
	Schema       string `long:"http-schema" description:"Http schema" env:"HTTP_SCHEMA" default:"http"`
	AllowOrigins string `long:"allow-origins" description:"Allow origins for CORS (separate by commas)" env:"ALLOW_ORIGINS" default:"*"`
}

type Endpoints struct {
	PortalsEndpoint         string `long:"portal-endpoint" description:"Portals gRpc endpoint address" env:"PORTALS_ENDPOINT" required:"true"`
	PortalsV2Endpoint       string `long:"portal-v2-endpoint" description:"Portals gRpc endpoint address" env:"PORTALSV2_ENDPOINT" required:"true"`
	SurveysEndpoint         string `long:"survey-endpoint" description:"Surveys gRpc endpoint address" env:"SURVEYS_ENDPOINT" required:"true"`
	AuthFacadeEndpoint      string `long:"auth-facade-endpoint" description:"Auth-facade gRpc endpoint address" env:"AUTH_FACADE_ENDPOINT" required:"true"`
	ProxyFacadeEndpoint     string `long:"proxy-facade-endpoint" description:"Proxy-facade gRpc endpoint address" env:"PROXY_FACADE_ENDPOINT" required:"true"`
	EmployeesSearchEndpoint string `long:"employees-search-endpoint" description:"Employees search gRpc endpoint address" env:"EMPLOYEES_SEARCH_ENDPOINT" required:"true"`
	FilesEndpoint           string `long:"files-endpoint" description:"Files gRpc endpoint address" env:"FILES_ENDPOINT" required:"true"`
	EmployeesEndpoint       string `long:"employees-endpoint" description:"Employees gRpc endpoint address" env:"EMPLOYEES_ENDPOINT" required:"true"`
	AnalyticsEndpoint       string `long:"analytics-endpoint" description:"Analytics gRpc endpoint address" env:"ANALYTICS_ENDPOINT" required:"true"`
	PortalsFacadeEndpoint   string `long:"portals-facade-endpoint" description:"Portals-facade gRpc endpoint address" env:"PORTALS_FACADE_ENDPOINT" required:"true"`
	NewsEndpoint            string `long:"news-endpoint" description:"News gRpc endpoint address" env:"NEWS_ENDPOINT" required:"true"`
	BannersEndpoint         string `long:"banners-endpoint" description:"Banners gRPC endpoint address" env:"BANNERS_ENDPOINT" required:"true"`
	NewsFacadeEndpoint      string `long:"news-facade-endpoint" description:"News-facade gRpc endpoint address" env:"NEWS_FACADE_ENDPOINT" required:"true"`
}

type TTL struct {
	AccessToken  time.Duration `long:"access-token-ttl" description:"Access token expiration time" env:"ACCESS_TOKEN_TTL" default:"2160h"`
	RefreshToken time.Duration `long:"refresh-token-ttl" description:"Refresh token expiration time" env:"REFRESH_TOKEN_TTL" default:"2160h"`
}

// NewConfig ...
func NewConfig() (*Config, error) {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(log.Writer())
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	return &cfg, nil
}
