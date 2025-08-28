package app

import "context"

//go:generate mockgen -source=interfaces.go -destination=./app_mock.go -package=app

type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
