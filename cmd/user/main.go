package main

import (
	"context"

	"github.com/tikivn/ultrago/u_graceful"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, logger := u_logger.GetLogger(u_graceful.NewCtx())

	eg, gctx := errgroup.WithContext(ctx)
	app, err := initUserApp(gctx)
	if err != nil {
		panic(err)
	}

	eg.Go(func() error {
		return app.Start(gctx)
	})

	defer func() {
		shutDownErr := app.Stop(context.Background())
		logger.Infof("API Server is shutdown with err=%v\n", shutDownErr)
	}()

	if err = eg.Wait(); err != nil {
		logger.Errorln(err)
	}
}
