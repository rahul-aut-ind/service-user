package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/middlewares"
	handlers "github.com/rahul-aut-ind/service-user/interfaceadapters/requesthandler"
	controllers "github.com/rahul-aut-ind/service-user/interfaceadapters/usercontroller"
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
		POST("", func(c *gin.Context) { r.controller.AddUserImage(c) })
}

// GET("", r.controller.GetAll).
// GET(fmt.Sprintf("/:%s", controller.ScanID), r.controller.Get).
// DELETE(fmt.Sprintf("/:%s", controller.ScanID), r.controller.Delete)
