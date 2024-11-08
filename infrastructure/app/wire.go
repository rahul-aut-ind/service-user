//go:build wireinject
// +build wireinject

package app

import (
	"github.com/rahul-aut-ind/service-user/infrastructure/caching"
	"github.com/rahul-aut-ind/service-user/infrastructure/routes"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/controllers/usercontroller"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/handlers/requesthandler"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/middlewares"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/userrepo"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/rahul-aut-ind/service-user/services/userservice"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func New(e *gin.Engine) (*App, error) {
	wire.Build(
		logger.Wired,

		config.Wired,

		requesthandler.Wired,

		middlewares.Wired,

		caching.Wired,
		wire.Bind(new(caching.CacheHandler), new(*caching.RedisClient)),

		userrepo.Wired,
		wire.Bind(new(userrepo.DBRepo), new(*userrepo.MysqlRepository)),

		userservice.Wired,
		wire.Bind(new(userservice.Services), new(*userservice.Service)),

		usercontroller.Wired,
		wire.Bind(new(usercontroller.UserHandler), new(*usercontroller.Controller)),

		routes.Wired,

		newApp,
	)

	return nil, nil
}
