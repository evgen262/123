package app

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	DialTimeout = 10 * time.Second
)

func (a *app) dial(_ context.Context, target string, required bool, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var (
		cc  *grpc.ClientConn
		err error
	)
	cc, err = grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't create client [%s]: %w", target, err)
	}
	defer func() {
		if err != nil {
			_ = cc.Close()
		}
	}()
	grpc.WithDefaultCallOptions()
	time.AfterFunc(DialTimeout, func() {
		defer func() {
			if e := recover(); e != nil {
				a.logger.Error("dial after-func", zap.Error(fmt.Errorf("%s", e)))
			}
		}()
		if cc.GetState() != connectivity.Ready && cc.GetState() != connectivity.Idle {
			a.logger.Error("can't connect to " + target)
			if required {
				a.stop()
			}
		}
	})

	return cc, nil
}
