package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/tree-alive.git"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

//go:generate mockgen -source=server.go -destination=./service_mock.go -package=service

const RequestTimeout = 3600 * time.Second

const (
	ContextKeyServiceAddr = "service-addr"
)

type HTTPSrv interface {
	Shutdown() error
	ShutdownWithContext(ctx context.Context) error
	ListenAndServe(addr string) error
}

type AppInfo struct {
	Name      string `json:"name"`
	Instance  string `json:"instance"`
	BuildTime string `json:"buildTime"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
}

type server struct {
	branch tree.Branch
	server HTTPSrv
}

func NewServer(treeBranch tree.Branch, appInfo *AppInfo) *server {
	return &server{
		branch: treeBranch,
		server: &fasthttp.Server{
			Handler: fasthttp.TimeoutWithCodeHandler(
				func(ctx *fasthttp.RequestCtx) {
					switch string(ctx.RequestURI()) {
					case "/healthz":
						if treeBranch.Tree().IsAlive() {
							ctx.SetStatusCode(fasthttp.StatusOK)
							ctx.SetBody([]byte("OK"))
						} else {
							ctx.SetStatusCode(fasthttp.StatusBadGateway)
						}
					case "/readyz":
						if treeBranch.Tree().IsReady() {
							ctx.SetStatusCode(fasthttp.StatusOK)
							ctx.SetBody([]byte("OK"))
						} else {
							ctx.SetStatusCode(fasthttp.StatusBadGateway)
						}
					case "/metrics":
						metricsHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
						metricsHandler(ctx)
					case "/info":
						body, err := json.Marshal(appInfo)
						if err != nil {
							ctx.SetStatusCode(fasthttp.StatusInternalServerError)
							return
						}
						ctx.SetStatusCode(fasthttp.StatusOK)
						ctx.SetBody(body)
					default:
						ctx.SetStatusCode(fasthttp.StatusNotImplemented)
					}
				},
				RequestTimeout,
				"request is timeout",
				fasthttp.StatusRequestTimeout,
			),
		},
	}
}

func (s *server) Run(ctx context.Context) error {
	addr, ok := ctx.Value(ContextKeyServiceAddr).(string)
	if !ok {
		return fmt.Errorf("service addr not a string: %+v", addr)
	}
	s.branch.Ready()
	if err := s.server.ListenAndServe(addr); err != nil {
		return fmt.Errorf("can't start service http server: %w", err)
	}
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	defer s.branch.Die()
	return s.server.ShutdownWithContext(ctx)
}
