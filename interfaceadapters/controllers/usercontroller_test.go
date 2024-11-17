package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/rahul-aut-ind/service-user/services/imageservice"
	"net/http"
	"testing"

	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/infrastructure/caching"
	"github.com/rahul-aut-ind/service-user/mocks"
	"github.com/rahul-aut-ind/service-user/services/userservice"
)

var (
	testUserResp = &models.User{
		ID:    1,
		Name:  "TestUser",
		Email: "testuser@test.com",
	}
)

func TestController_FindUser_Success(t *testing.T) {
	repoMoc := new(mocks.DBRepo)
	contextMoc := new(mocks.Context)
	cacheMoc := new(mocks.CacheHandler)

	contextMoc.On("Param", "id").Return("1")
	contextMoc.On("JSON", http.StatusOK, &models.Response{Data: testUserResp})

	testService := userservice.New(repoMoc, logger.New())
	testImageService := imageservice.New(nil, nil, logger.New())
	testContrlr := New(cacheMoc, testService, testImageService, logger.New())

	cacheMoc.On("Get", contextMoc, "1").Return("", errors.New("err : %s", fmt.Errorf("no data in cache")))
	repoMoc.On("FindRecord", "1").Return(testUserResp, nil)
	u, _ := json.Marshal(testUserResp)
	cacheMoc.On("Set", contextMoc, "1", string(u), caching.DefaultTTL).Return(nil)

	// When
	testContrlr.FindUser(contextMoc)

	// Then
	repoMoc.AssertNumberOfCalls(t,
		"FindRecord",
		1,
	)
	repoMoc.AssertExpectations(t)
}

func TestController_FindUser_NoRecordsErr(t *testing.T) {
	repoMoc := new(mocks.DBRepo)
	contextMoc := new(mocks.Context)
	cacheMoc := new(mocks.CacheHandler)

	contextMoc.On("Param", "id").Return("9999")

	repoFindErr := fmt.Errorf("%s", errors.ErrCodeNoUser)
	respErr := errors.New(errors.ErrCodeNoUser, fmt.Errorf("error :: error :: %v", repoFindErr))

	contextMoc.On("JSON", http.StatusNotFound, respErr)

	testService := userservice.New(repoMoc, logger.New())
	testImageService := imageservice.New(nil, nil, logger.New())
	testContrlr := New(cacheMoc, testService, testImageService, logger.New())

	cacheMoc.On("Get", contextMoc, "9999").Return("", fmt.Errorf("no data in cache"))
	repoMoc.On("FindRecord", "9999").Return(nil, repoFindErr)

	// When
	testContrlr.FindUser(contextMoc)

	// Then
	repoMoc.AssertNumberOfCalls(t,
		"FindRecord",
		1,
	)
	repoMoc.AssertExpectations(t)
}

func TestController_FindUser_RegexBadReq(t *testing.T) {
	repoMoc := new(mocks.DBRepo)
	contextMoc := new(mocks.Context)
	cacheMoc := new(mocks.CacheHandler)

	// param doesn't match regex
	contextMoc.On("Param", "id").Return("-1")

	respErr := errors.New(errors.ErrCodeBadRequest, fmt.Errorf("bad request"))

	contextMoc.On("JSON", http.StatusBadRequest, respErr)

	testService := userservice.New(repoMoc, logger.New())
	testImageService := imageservice.New(nil, nil, logger.New())
	testContrlr := New(cacheMoc, testService, testImageService, logger.New())

	// When
	testContrlr.FindUser(contextMoc)

	// Then
	repoMoc.AssertNumberOfCalls(t,
		"FindRecord",
		0,
	)
	repoMoc.AssertExpectations(t)
}

func TestController_FindUser_RepoErr(t *testing.T) {
	repoMoc := new(mocks.DBRepo)
	contextMoc := new(mocks.Context)
	cacheMoc := new(mocks.CacheHandler)

	contextMoc.On("Param", "id").Return("1")

	repoErr := fmt.Errorf("some internal err")
	respErr := errors.New(errors.ErrCodeGeneric, fmt.Errorf("error :: error :: %v", repoErr))

	contextMoc.On("JSON", http.StatusInternalServerError, respErr)

	testService := userservice.New(repoMoc, logger.New())
	testImageService := imageservice.New(nil, nil, logger.New())
	testContrlr := New(cacheMoc, testService, testImageService, logger.New())

	cacheMoc.On("Get", contextMoc, "1").Return("", fmt.Errorf("no data in cache"))
	repoMoc.On("FindRecord", "1").Return(nil, repoErr)

	// When
	testContrlr.FindUser(contextMoc)

	// Then
	repoMoc.AssertNumberOfCalls(t,
		"FindRecord",
		1,
	)
	repoMoc.AssertExpectations(t)
}
