package grpc

import (
	"context"

	"git.mos.ru/buch-cloud/moscow-team-2.0/build/ditzap.git"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type interceptor struct {
}

func NewInterceptor() *interceptor {
	return &interceptor{}
}

func (i *interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			return handler(metadata.NewOutgoingContext(ctx, md), req)
		}
		return handler(ctx, req)
	}
}

func (i *interceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			if err := ss.SetHeader(md); err != nil {
				return err
			}
		}
		return handler(srv, ss)
	}
}

func (i *interceptor) UnaryLogRequest() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		return handler(context.WithValue(ctx, ditzap.ReqKey, req), req)
	}
}
