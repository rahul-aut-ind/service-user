package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/infrastructure/app"
)

func main() {
	log := logger.New()

	log.Info(">>>>>   service-user   <<<<<<")

	gin.SetMode("debug")
	r := gin.New()

	a, err := app.New(r)
	if err != nil {
		log.Fatal(err)
	}

	a.Start()
}
