package user

import (
	"context"
	"fmt"
	"net"
	"user_crud/internal/pkg/setting"

	"github.com/tikivn/ultrago/u_graceful"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"
)

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer *HTTPServer
}

func NewApp(httpServer *HTTPServer) App {
	return &app{
		httpServer: httpServer,
	}
}

func (a *app) Start(ctx context.Context) (err error) {

	// Get logger from context
	ctx, logger := u_logger.GetLogger(ctx)

	// Create httpServer listen on setting.HttpPort
	httpLis, err := net.Listen("tcp", setting.HttpPort)
	if err != nil {
		logger.Fatal("failed to listen http port %s: %v", setting.HttpPort, err)
		return err
	}

	eg, childCtx := errgroup.WithContext(ctx)

	// Using u_graceful to blocking (waiting)
	// for function BlockListen to waiting the a.httpServer.Serve return error or
	// the context is done
	eg.Go(func() error {
		return u_graceful.BlockListen(childCtx, func() error {
			logger.Infof("start listening http request on %v", setting.HttpPort)
			if err := a.httpServer.Serve(httpLis); err != nil {
				return fmt.Errorf("failed to serve http: %v", err)
			}
			return nil
		})
	})

	logger.Infof("user app started")

	// Wait for errgroup has done all the job
	err = eg.Wait()

	return
}

func (a *app) Stop(ctx context.Context) (err error) {

	ctx, logger := u_logger.GetLogger(ctx)

	logger.Infof("stop listening http request on %v", setting.HttpPort)
	err = a.httpServer.Shutdown(ctx)

	return
}
