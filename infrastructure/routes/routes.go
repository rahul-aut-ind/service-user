package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rahul-aut-ind/service-user/interfaceadapters/controllers"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/middlewares"
	handlers "github.com/rahul-aut-ind/service-user/interfaceadapters/requesthandler"
)

type Routes struct {
	handler    handlers.RequestHandler
	controller controllers.Handler
	validator  middlewares.Validator
}

func New(
	h handlers.RequestHandler,
	c controllers.Handler,
	v middlewares.Validator,
) *Routes {
	return &Routes{
		handler:    h,
		controller: c,
		validator:  v,
	}
}

// nolint:dupl // different route groups
func (r *Routes) Setup() {
	// Public
	r.handler.Gin.Group("/api/v1/users").
		Use(r.validator.ValidateRequest()).
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

	r.handler.Gin.Group("/api/v1/user-image").
		Use(r.validator.ValidateRequest()).
		// upload a image
		POST("", func(c *gin.Context) { r.controller.CreateUserImage(c) }).
		// get an user image
		GET("/:id", func(c *gin.Context) { r.controller.GetUserImage(c) }).
		// get all user images
		GET("", func(c *gin.Context) { r.controller.GetAllUserImages(c) }).
		// delete an user image
		DELETE("/:id", func(c *gin.Context) { r.controller.DeleteUserImage(c) }).
		// deletes all user images
		DELETE("", func(c *gin.Context) { r.controller.DeleteAllUserImages(c) })
}
