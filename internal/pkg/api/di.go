package api

import "github.com/google/wire"

var GraphSet = wire.NewSet(
	NewUserHandler,
	NewBaseHandler,
)
