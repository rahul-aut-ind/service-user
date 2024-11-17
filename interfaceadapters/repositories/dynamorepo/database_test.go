package dynamorepo

import (
	"testing"
	"time"

	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/integrationtest"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
	repo        *DynamoDBRepo
	dynamoSetup *integrationtest.DynamoDBSetup
}

func (s *RepoTestSuite) SetupSuite() {
	s.dynamoSetup = integrationtest.NewDynamoDbSetup()
}

func (s *RepoTestSuite) SetupTest() {
	s.repo = &DynamoDBRepo{
		TableName: integrationtest.UserImageTable,
		Client:    s.dynamoSetup.Client,
		Log:       logger.New(),
	}
}

func (s *RepoTestSuite) TearDownSuite() {
	s.dynamoSetup.Stop()
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) TestShouldSaveImage() {

	createReq := &models.UserImage{
		IsDeleted: false,
		UserID:    "123",
		ImageID:   "128a68e4-a10a-11ef-ba63-c689f470ad55",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReq)

	result, _ := s.repo.GetImage("123", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "128a68e4-a10a-11ef-ba63-c689f470ad55", result.ImageID)
}

func (s *RepoTestSuite) TestShouldNotSaveImageWithoutUserID() {

	createReqWithoutUserID := &models.UserImage{
		IsDeleted: false,
		ImageID:   "128a68e4-a10a-11ef-ba63-c689f470ad55",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReqWithoutUserID)

	assert.NotNil(s.T(), err)
}

func (s *RepoTestSuite) TestShouldNotSaveImageWithoutImageID() {

	createReqWithoutImageID := &models.UserImage{
		IsDeleted: false,
		UserID:    "123",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReqWithoutImageID)

	assert.NotNil(s.T(), err)
}

func (s *RepoTestSuite) TestShouldDeleteImage() {

	createReq := &models.UserImage{
		IsDeleted: false,
		UserID:    "123",
		ImageID:   "128a68e4-a10a-11ef-ba63-c689f470ad55",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReq)

	err = s.repo.DeleteImage("123", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.Nil(s.T(), err)

	_, err = s.repo.GetImage("123", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "image not found", err.Error())
}

func (s *RepoTestSuite) TestShouldDeleteAllImage() {

	createReq := &models.UserImage{
		IsDeleted: false,
		UserID:    "123",
		ImageID:   "128a68e4-a10a-11ef-ba63-c689f470ad55",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReq)

	err = s.repo.DeleteAllImages("123")
	assert.Nil(s.T(), err)

	_, err = s.repo.GetImage("123", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "image not found", err.Error())
}

func (s *RepoTestSuite) TestShouldNotDeleteImageIfNotExist() {

	err := s.repo.DeleteImage("345", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "image not found", err.Error())

	_, err = s.repo.GetImage("345", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "image not found", err.Error())
}

func (s *RepoTestSuite) TestShouldGetImage() {

	createReq := &models.UserImage{
		IsDeleted: false,
		UserID:    "123",
		ImageID:   "128a68e4-a10a-11ef-ba63-c689f470ad55",
		Path:      "story-image/123/128a68e4-a10a-11ef-ba63-c689f470ad55.jpg",
		TakenAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateOrUpdateImage(createReq)

	result, _ := s.repo.GetImage("123", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "128a68e4-a10a-11ef-ba63-c689f470ad55", result.ImageID)
}

func (s *RepoTestSuite) TestShouldNotGetImageIfNotExist() {

	_, err := s.repo.GetImage("345", "128a68e4-a10a-11ef-ba63-c689f470ad55")
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), "image not found", err.Error())
}

func (s *RepoTestSuite) TestShouldGetAllImagePaginated() {

	todayTime := time.Now()
	yesterdayTime := todayTime.AddDate(0, 0, -1)

	createReq1 := &models.UserImage{
		IsDeleted: false,
		UserID:    "999",
		ImageID:   "1111111-1111111",
		Path:      "story-image/123/1111111-1111111.jpg",
		TakenAt:   todayTime,
		UpdatedAt: todayTime,
	}
	err := s.repo.CreateOrUpdateImage(createReq1)

	createReq2 := &models.UserImage{
		IsDeleted: false,
		UserID:    "999",
		ImageID:   "22222222-22222222",
		Path:      "story-image/123/22222222-22222222.jpg",
		TakenAt:   yesterdayTime,
		UpdatedAt: yesterdayTime,
	}
	err = s.repo.CreateOrUpdateImage(createReq2)

	getReq1 := models.PaginatedInput{
		UserID:           "999",
		LastImageID:      "",
		LastImageTakenAt: "",
		Limit:            1,
	}

	data, err := s.repo.GetAllImagesPaginated(getReq1)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data.UserImages))
	assert.Equal(s.T(), "1111111-1111111", data.UserImages[0].ImageID)
	assert.Equal(s.T(), "1111111-1111111", data.Page.LastEvaluatedKey[config.QueryParamLastKey])

	getReq2 := models.PaginatedInput{
		UserID:           "999",
		LastImageID:      data.Page.LastEvaluatedKey[config.QueryParamLastKey],
		LastImageTakenAt: data.Page.LastEvaluatedKey[config.QueryParamlastKeyDate],
		Limit:            1,
	}

	data, err = s.repo.GetAllImagesPaginated(getReq2)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data.UserImages))
	assert.Equal(s.T(), "22222222-22222222", data.UserImages[0].ImageID)
	assert.Equal(s.T(), "22222222-22222222", data.Page.LastEvaluatedKey[config.QueryParamLastKey])

	getReq3 := models.PaginatedInput{
		UserID:           "999",
		LastImageID:      data.Page.LastEvaluatedKey[config.QueryParamLastKey],
		LastImageTakenAt: data.Page.LastEvaluatedKey[config.QueryParamlastKeyDate],
		Limit:            1,
	}

	data, err = s.repo.GetAllImagesPaginated(getReq3)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data.UserImages))
	assert.Equal(s.T(), 0, len(data.Page.LastEvaluatedKey))
}

func (s *RepoTestSuite) TestShouldNotGetAllImagePaginatedIfNotExist() {

	req := models.PaginatedInput{
		UserID:           "234",
		LastImageID:      "128a68e4-a10a-11ef-ba63-c689f470ad55",
		LastImageTakenAt: "2024-10-11T00:00:00Z",
		Limit:            1,
	}

	data, err := s.repo.GetAllImagesPaginated(req)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data.UserImages))
	assert.Equal(s.T(), 0, len(data.Page.LastEvaluatedKey))
}
