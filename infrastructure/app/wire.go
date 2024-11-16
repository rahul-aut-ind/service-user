//go:build wireinject
// +build wireinject

package app

import (
	"github.com/rahul-aut-ind/service-user/infrastructure/caching"
	"github.com/rahul-aut-ind/service-user/infrastructure/routes"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/middlewares"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/dynamorepo"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/mysqlrepo"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/s3repo"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/requesthandler"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/usercontroller"
	usercontroller2 "github.com/rahul-aut-ind/service-user/interfaceadapters/usercontroller"
	"github.com/rahul-aut-ind/service-user/internal/awsconfig"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/rahul-aut-ind/service-user/services/imageservice"
	"github.com/rahul-aut-ind/service-user/services/userservice"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func New(e *gin.Engine) (*App, error) {
	wire.Build(
		logger.Wired,

		config.Wired,

		awsconfig.Wired,

		requesthandler.Wired,

		middlewares.Wired,

		caching.Wired,
		wire.Bind(new(caching.CacheHandler), new(*caching.RedisClient)),

		mysqlrepo.Wired,
		wire.Bind(new(mysqlrepo.DataHandler), new(*mysqlrepo.MysqlClient)),

		s3repo.Wired,
		wire.Bind(new(s3repo.S3Handler), new(*s3repo.S3Repo)),

		dynamorepo.Wired,
		wire.Bind(new(dynamorepo.DataHandler), new(*dynamorepo.DynamoDBRepo)),

		userservice.Wired,
		wire.Bind(new(userservice.UserService), new(*userservice.Service)),

		imageservice.Wired,
		wire.Bind(new(imageservice.UserImageService), new(*imageservice.Service)),

		usercontroller.Wired,
		wire.Bind(new(usercontroller2.Handler), new(*usercontroller2.Controller)),

		routes.Wired,

		newApp,
	)

	return nil, nil
}
