package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rahul-aut-ind/service-user/interfaceadapters/controllers/usercontroller"
	handlers "github.com/rahul-aut-ind/service-user/interfaceadapters/handlers/requesthandler"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/middlewares"
)

type Routes struct {
	handler    handlers.RequestHandler
	controller controllers.UserHandler
	validator  middlewares.Validator
}

func New(
	h handlers.RequestHandler,
	c controllers.UserHandler,
	v middlewares.Validator,
) *Routes {
	return &Routes{
		handler:    h,
		controller: c,
		validator:  v,
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
		Use(r.validator.RequestValidator()).
		// create user
		POST("", func(c *gin.Context) { r.controller.CreateUser(c) }).
		// update user by id
		PUT("/:id", func(c *gin.Context) { r.controller.UpdateUser(c) }).
		// Query specific
		GET("/:id", func(c *gin.Context) { r.controller.FindUser(c) }).
		// Query all users
		GET("", func(c *gin.Context) { r.controller.FindAllUsers(c) }).
		// delete by id
		DELETE("/:id", func(c *gin.Context) { r.controller.DeleteUser(c) })
}
