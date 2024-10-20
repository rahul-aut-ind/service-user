package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rahul-aut-ind/service-user/interfaceadapters/controllers/usercontroller"
	handlers "github.com/rahul-aut-ind/service-user/interfaceadapters/handlers/requesthandler"
)

type Routes struct {
	handler    handlers.RequestHandler
	controller controllers.IController
}

func New(
	rh handlers.RequestHandler,
	c controllers.IController,
) *Routes {
	return &Routes{
		handler:    rh,
		controller: c,
	}
}

func (r *Routes) Setup() {
	// Private
	// r.handler.Gin.Group("/service-user/api/v3/users").
	// 	// // Query all users
	// 	// GET("", func(c *gin.Context) { r.controller.FindAllUsers(c) }).
	// 	// // Query specific
	// 	// GET("/:id", func(c *gin.Context) { r.controller.FindUser(c) })

	// Public
	r.handler.Gin.Group("/api/v1/users").
		// create user
		POST("", func(c *gin.Context) { r.controller.CreateUser(c) }).
		// Query specific
		GET("/:id", func(c *gin.Context) { r.controller.FindUser(c) })
	// // update by id
	// PUT("/:id", func(c *gin.Context) { r.controller.UpdateUser(c) }).
	// // delete by id
	// DELETE("/:id", func(c *gin.Context) { r.controller.DeleteUser(c) })

}
