package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rahul-aut-ind/service-user/infrastructure/routes"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type App struct {
	route  *routes.Routes
	engine *gin.Engine
	env    *config.Env
	log    *logger.Logger
}

func newApp(r *routes.Routes, env *config.Env, l *logger.Logger, e *gin.Engine) *App {
	return &App{route: r, env: env, log: l, engine: e}
}

func (a *App) Start() {
	a.engine.Use(a.log.DefaultLogger())
	a.route.Setup()
	err := a.engine.Run(fmt.Sprintf("%s:%s", a.env.ServerHost, a.env.ServerPort))
	if err != nil {
		a.log.Fatalf("could not start the server | err :: %v", err)
	}
}
