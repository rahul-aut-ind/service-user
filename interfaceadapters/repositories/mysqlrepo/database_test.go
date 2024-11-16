package mysqlrepo

import (
	"strconv"
	"testing"

	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/integrationtest"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
	repo    *MysqlClient
	dbSetup *integrationtest.MySQLSetup
}

func (s *RepoTestSuite) SetupSuite() {
	s.dbSetup = integrationtest.NewMySQLSetup()
}

func (s *RepoTestSuite) SetupTest() {
	s.repo = &MysqlClient{
		client: connect(s.dbSetup.ConnString),
		log:    logger.New(),
	}
}

func (s *RepoTestSuite) TearDownSuite() {
	s.dbSetup.Stop()
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (s *RepoTestSuite) TestShouldCreateUser() {
	res, err := s.repo.CreateRecord(&models.User{
		Name:  "test1",
		Email: "test1@test.com",
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test1", res.Name)
	assert.Equal(s.T(), "test1@test.com", res.Email)
}

func (s *RepoTestSuite) TestShouldGetUser() {
	res, _ := s.repo.CreateRecord(&models.User{
		Name:  "test2",
		Email: "test2@test.com",
	})
	u, err := s.repo.FindRecord(strconv.Itoa(int(res.ID)))

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test2", u.Name)
	assert.Equal(s.T(), "test2@test.com", u.Email)
}

func (s *RepoTestSuite) TestShouldUpdateUser() {
	res, _ := s.repo.CreateRecord(&models.User{
		Name:  "test3",
		Email: "test3@test.com",
	})

	u, err := s.repo.UpdateRecord(&models.User{
		ID:    res.ID,
		Name:  "test3_updated",
		Email: "test3@test.com",
	})

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test3_updated", u.Name)
	assert.Equal(s.T(), "test3@test.com", u.Email)
}

func (s *RepoTestSuite) TestShouldDeleteUser() {
	res, _ := s.repo.CreateRecord(&models.User{
		Name:  "test4",
		Email: "test4@test.com",
	})

	_, err := s.repo.DeleteRecord(&models.User{
		ID: res.ID,
	})

	assert.Nil(s.T(), err)
}

func (s *RepoTestSuite) TestShouldGetAllUsers() {
	res1, _ := s.repo.CreateRecord(&models.User{
		Name:  "test5",
		Email: "test5@test.com",
	})
	res2, _ := s.repo.CreateRecord(&models.User{
		Name:  "test6",
		Email: "test6@test.com",
	})

	records, err := s.repo.ListRecords()

	assert.Nil(s.T(), err)

	results := make([]string, 0, len(records))
	for _, i := range records {
		results = append(results, i.Email, i.Name)
	}

	assert.Contains(s.T(), results, res1.Email)
	assert.Contains(s.T(), results, res1.Name)
	assert.Contains(s.T(), results, res2.Email)
	assert.Contains(s.T(), results, res2.Name)

}
