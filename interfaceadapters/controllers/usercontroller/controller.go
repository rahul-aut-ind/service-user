package usercontroller

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/services/userservice"
)

type (
	IController interface {
		FindUser(c Context)
		FindAllUsers(c Context)
		CreateUser(c Context)
		UpdateUser(c Context)
		DeleteUser(c Context)
	}

	Controller struct {
		service userservice.IService
		log     *logger.Logger
	}

	Context interface {
		JSON(code int, obj interface{})
		ShouldBindJSON(obj interface{}) error
		GetHeader(key string) string
		Param(key string) string
		Query(key string) string
	}
)

var (
	userIDRegExp = regexp.MustCompile(`^\d+$`)
)

func New(s userservice.IService, l *logger.Logger) *Controller {
	return &Controller{
		service: s,
		log:     l,
	}
}

func (uc *Controller) CreateUser(c Context) {
	newUser := &models.User{}

	err := c.ShouldBindJSON(newUser)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	user, err := uc.service.Add(newUser)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusAccepted, &models.Response{Data: user})
}

func (uc *Controller) FindUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	user, err := uc.service.Get(userID)
	if err != nil {
		if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
			uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
			return
		}
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusOK, &models.Response{Data: user})
}

func (uc *Controller) DeleteUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	err := uc.service.Delete(userID)
	if err != nil {
		if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
			uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
			return
		}
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusAccepted, &models.Response{Data: models.RequestAccepted})
}

func (uc *Controller) UpdateUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	updatedUserInfo := &models.User{}

	err := c.ShouldBindJSON(updatedUserInfo)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	user, err := uc.service.Update(userID, updatedUserInfo)
	if err != nil {
		if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
			uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
			return
		}
		uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusOK, &models.Response{Data: user})
}

func (uc *Controller) FindAllUsers(c Context) {
	user, err := uc.service.GetAll()
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusOK, &models.Response{Data: user})
}

func (uc *Controller) handleError(c Context, err error) {
	var apiErr errors.Error
	if e, ok := err.(errors.Error); ok {
		apiErr = e
	} else {
		apiErr = errors.New(errors.ErrCodeGeneric, err)
	}
	uc.log.Errorf("error :: %v", err)
	c.JSON(apiErr.HTTPCode(), apiErr)
}
