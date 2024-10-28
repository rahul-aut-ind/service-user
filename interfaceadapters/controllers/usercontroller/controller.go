package usercontroller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/infrastructure/caching"
	"github.com/rahul-aut-ind/service-user/services/userservice"
)

type (
	UserHandler interface {
		FindUser(c Context)
		FindAllUsers(c Context)
		CreateUser(c Context)
		UpdateUser(c Context)
		DeleteUser(c Context)
	}

	Controller struct {
		rc      caching.CacheHandler
		service userservice.Services
		log     *logger.Logger
	}

	Context interface {
		JSON(code int, obj interface{})
		ShouldBindJSON(obj interface{}) error
		GetHeader(key string) string
		Param(key string) string
		Query(key string) string
		Value(key any) any
		Err() error
		Done() <-chan struct{}
		Deadline() (deadline time.Time, ok bool)
	}
)

var (
	userIDRegExp = regexp.MustCompile(`^\d+$`)
	validate     = validator.New()
)

func validateInput(input models.User) error {
	return validate.Struct(input)
}

func New(rc caching.CacheHandler, s userservice.Services, l *logger.Logger) *Controller {
	return &Controller{
		rc:      rc,
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

	if err := validateInput(*newUser); err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	user, err := uc.service.AddUser(newUser)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
		return
	}
	u, _ := json.Marshal(user)
	err = uc.rc.Set(c, strconv.Itoa(int(user.ID)), string(u), caching.DefaultTTL)
	if err != nil {
		uc.log.Warnf("err updating cache :: %s", err)
	}
	c.JSON(http.StatusAccepted, &models.Response{Data: user})
}

func (uc *Controller) FindUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	// check if data exists in redis
	cachedData, err := uc.rc.Get(c, userID)
	if err != nil {
		uc.log.Debug("cache miss")
		user, err := uc.service.GetUserWithID(userID)
		if err != nil {
			if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
				uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
				return
			}
			uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
			return
		}
		u, _ := json.Marshal(user)
		err = uc.rc.Set(c, userID, string(u), caching.DefaultTTL)
		if err != nil {
			uc.log.Warnf("err updating cache :: %s", err)
		}
		c.JSON(http.StatusOK, &models.Response{Data: user})
		return
	}
	uc.log.Debug("serving data from cache..")
	data := models.User{}
	err = json.Unmarshal([]byte(cachedData), &data)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
	}
	c.JSON(http.StatusOK, &models.Response{Data: &data})
}

func (uc *Controller) DeleteUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	err := uc.service.DeleteUser(userID)
	if err != nil {
		if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
			uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
			return
		}
		uc.handleError(c, errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: %v", err)))
		return
	}
	err = uc.rc.Delete(c, userID)
	if err != nil {
		uc.log.Warnf("err updating cache :: %s", err)
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

	if err := validateInput(*updatedUserInfo); err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	user, err := uc.service.UpdateUser(userID, updatedUserInfo)
	if err != nil {
		if strings.Contains(err.Error(), models.ErrMsgNoUserfound) {
			uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
			return
		}
		uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
		return
	}
	u, _ := json.Marshal(user)
	err = uc.rc.Set(c, userID, string(u), caching.DefaultTTL)
	if err != nil {
		uc.log.Warnf("err updating cache :: %s", err)
	}
	c.JSON(http.StatusOK, &models.Response{Data: user})
}

func (uc *Controller) FindAllUsers(c Context) {
	user, err := uc.service.GetAllUsers()
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
