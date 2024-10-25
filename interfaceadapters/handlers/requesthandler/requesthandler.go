package requesthandler

import (
	"github.com/gin-gonic/gin"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

// New creates a new request handler
func New(e *gin.Engine) RequestHandler {
	return RequestHandler{Gin: e}
}
