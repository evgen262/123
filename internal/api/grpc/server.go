package grpc

import (
	"context"
	"fmt"
	"net"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"git.mos.ru/buch-cloud/moscow-team-2.0/build/tree-alive.git"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"goa.design/goa/v3/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	authv1 "git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/gen/infogorod/auth/auth/v1"
	"git.mos.ru/buch-cloud/moscow-team-2.0/pud/auth.git/internal/api/grpc/presenter"
)

type server struct {
	addr   string
	server *grpc.Server

	treeBranch tree.Branch
	logger     ditzap.Logger
}

func NewServer(
	addr string,
	treeBranch tree.Branch,
	logger ditzap.Logger) *server {
	grpcServer := &server{
		addr:       addr,
		treeBranch: treeBranch,
		logger:     logger,
	}

	recoveryHandler := func(p interface{}) (err error) {
		return fmt.Errorf("grpc recovery from panic: %s", p)
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandler),
	}
	interceptor := NewInterceptor()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			middleware.UnaryRequestID(
				middleware.UseXRequestIDMetadataOption(true),
				middleware.XRequestMetadataLimitOption(128),
			),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			interceptor.UnaryLogRequest(),
			grpc_zap.UnaryServerInterceptor(logger.GetLogger(), grpc_zap.WithLevels(grpcServer.grpcCodeToZapLevel), grpc_zap.WithMessageProducer(ditzap.MessageProducer(logger))),
			grpc_validator.UnaryServerInterceptor(),
			interceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			middleware.StreamRequestID(
				middleware.UseXRequestIDMetadataOption(true),
				middleware.XRequestMetadataLimitOption(128),
			),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(logger.GetLogger(), grpc_zap.WithLevels(grpcServer.grpcCodeToZapLevel), grpc_zap.WithMessageProducer(ditzap.MessageProducer(logger))),
			grpc_validator.StreamServerInterceptor(),
			interceptor.Stream(),
		),
	)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s)
	grpcServer.server = s

	return grpcServer
}

func (s *server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc listen error: %w", err)
	}

	go func() {
		<-ctx.Done()
		err := s.Shutdown(ctx)
		if err != nil {
			s.logger.Error("can't shutdown grpc server", zap.Error(err))
			return
		}
	}()

	s.treeBranch.Ready()
	s.logger.Info("grpc server started", zap.String("server-host", s.addr))
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("grpc server shutdown: %w", err)
	}
	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *server) grpcCodeToZapLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zapcore.DebugLevel
	default:
		return grpc_zap.DefaultCodeToLevel(code)
	}
}

func (s *server) RegisterServers(authUsecase AuthInteractor) {
	authPresenter := presenter.NewAuthPresenter()

	// инициализация grpc ручек
	authv1.RegisterAuthAPIServer(s.server, NewAuthServer(authUsecase, authPresenter))

	// Серверная рефлексия
	reflection.Register(s.server)
}
