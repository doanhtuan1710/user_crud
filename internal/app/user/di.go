package user

import (
	"user_crud/internal/pkg/api"
	"user_crud/internal/pkg/infra"
	"user_crud/internal/pkg/repo/mysql"

	"github.com/google/wire"
)

var deps = wire.NewSet(
	infra.GraphSet,
	mysql.GraphSet,
	api.GraphSet,
)

var GraphSet = wire.NewSet(
	deps,
	NewApp,
	NewHTTPServer,
)
