package usercontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/infrastructure/caching"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/requestparser"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/rahul-aut-ind/service-user/services/imageservice"
	"github.com/rahul-aut-ind/service-user/services/userservice"
)

type (
	Handler interface {
		FindUser(c Context)
		FindAllUsers(c Context)
		CreateUser(c Context)
		UpdateUser(c Context)
		DeleteUser(c Context)
		AddUserImage(c Context)
	}

	Controller struct {
		rc           caching.CacheHandler
		userService  userservice.UserService
		imageService imageservice.UserImageService
		log          *logger.Logger
		val          *validator.Validate
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
		Copy() *gin.Context
	}
)

const (
	RequestAccepted = "ok"
)

var (
	userIDRegExp = regexp.MustCompile(`^\d+$`)
)

func New(rc caching.CacheHandler, us userservice.UserService, is imageservice.UserImageService, l *logger.Logger) *Controller {
	return &Controller{
		rc:           rc,
		userService:  us,
		imageService: is,
		log:          l,
		val:          validator.New(),
	}
}

func (uc *Controller) CreateUser(c Context) {
	req := &models.Request{}

	err := c.ShouldBindJSON(req)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	if err := uc.validateInput(*req); err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}
	newUser := &models.User{
		Name:    req.FirstName + " " + req.LastName,
		Email:   req.Email,
		Address: req.Address,
		Age:     req.Age,
	}

	user, err := uc.userService.AddUser(newUser)
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
		user, err := uc.userService.GetUserWithID(userID)
		if err != nil {
			if strings.Contains(err.Error(), errors.ErrCodeNoUser) {
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

	err := uc.userService.DeleteUser(userID)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrCodeNoUser) {
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
	c.JSON(http.StatusAccepted, &models.Response{Data: RequestAccepted})
}

func (uc *Controller) UpdateUser(c Context) {
	userID := c.Param("id")
	if !(userIDRegExp.MatchString(userID)) {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request")))
		return
	}

	req := &models.Request{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}

	if err := uc.validateInput(*req); err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request. Err :: %v", err)))
		return
	}
	updatedUserInfo := &models.User{
		Name:    req.FirstName + " " + req.LastName,
		Email:   req.Email,
		Address: req.Address,
		Age:     req.Age,
	}

	user, err := uc.userService.UpdateUser(userID, updatedUserInfo)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrCodeNoUser) {
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
	user, err := uc.userService.GetAllUsers()
	if err != nil {
		uc.handleError(c, errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: %v", err)))
		return
	}
	c.JSON(http.StatusOK, &models.Response{Data: user})
}

func (uc *Controller) AddUserImage(c Context) {
	body, err := io.ReadAll(c.Copy().Request.Body)

	if err != nil {
		uc.handleError(c, err)
		return
	}

	rp := &requestparser.RequestParser{
		Log:         uc.log,
		Body:        body,
		ContentType: c.GetHeader(config.HeaderContentType),
	}

	req, err := rp.ParseMultipart()
	if err != nil {
		uc.handleError(c, err)
		return
	}

	resp, err := uc.imageService.SaveImage(c.GetHeader(config.HeaderUserID), req)
	if err != nil {
		uc.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (uc *Controller) validateInput(input models.Request) error {
	return uc.val.Struct(input)
}

func (uc *Controller) handleError(c Context, err error) {
	var apiErr errors.Error
	if e, ok := err.(errors.Error); ok {
		apiErr = e
	} else {
		apiErr = errors.New(errors.ErrCodeGeneric, err)
	}
	uc.log.Errorf("error :: %s", err)
	c.JSON(apiErr.HTTPCode(), apiErr)
}
