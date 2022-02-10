//go:build wireinject

package main

import (
	"context"
	"user_crud/internal/app/user"

	"github.com/google/wire"
)

func initUserApp(ctx context.Context) (userApp user.App, err error) {
	wire.Build(user.GraphSet)
	return
}
