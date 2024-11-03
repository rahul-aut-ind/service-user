package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type (
	Validator struct {
		log *logger.Logger
	}
)

func New(l *logger.Logger) Validator {
	return Validator{log: l}
}

func (v *Validator) RequestValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if IDToken := ctx.GetHeader(config.HeaderIDToken); IDToken != "" {
			ctx.Next()
		} else {
			e := fmt.Errorf("required header %s not available", config.HeaderIDToken)
			v.log.Warnf("err :: %s", e)
			_ = ctx.Error(e)
			ctx.AbortWithStatusJSON(http.StatusForbidden, errors.New(errors.ErrCodeInvalidUserIDHeader, e))
		}
	}
}
