package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rahul-aut-ind/service-user/infrastructure/app"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

func main() {
	log := logger.New()
	log.Info(">>>>>   service-user   <<<<<<")
	gin.SetMode("debug")
	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(log.DefaultLogger())

	a, err := app.New(e)
	if err != nil {
		log.Fatal(err)
	}

	a.Start()
}
