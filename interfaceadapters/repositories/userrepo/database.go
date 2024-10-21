package userrepo

import (
	"fmt"

	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DBRepo interface {
		ListRecords() ([]models.User, error)
		FindRecord(id string) (*models.User, error)
		CreateRecord(u *models.User) (*models.User, error)
		UpdateRecord(u *models.User) (*models.User, error)
		DeleteRecord(u *models.User) (*models.User, error)
	}

	MysqlRepository struct {
		db  *gorm.DB
		log *logger.Logger
	}
)

func New(l *logger.Logger, env *config.Env) *MysqlRepository {
	return &MysqlRepository{db: connect(env.DBConnectionString), log: l}
}

// Connect initializes the database connection
func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database :: %v", err))
	}

	// Auto-migrate the User model
	db.AutoMigrate(&models.User{})

	return db
}

func (repo *MysqlRepository) CreateRecord(u *models.User) (*models.User, error) {
	repo.log.Debugf("inserting record %v", *u)
	result := repo.db.Create(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}

func (repo *MysqlRepository) FindRecord(id string) (*models.User, error) {
	repo.log.Debugf("finding record with id %s", id)
	var user models.User
	result := repo.db.Where(id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("err :: %v", gorm.ErrRecordNotFound)
		} else {
			return nil, fmt.Errorf("err :: %v", result.Error)
		}
	}
	return &user, nil
}

func (repo *MysqlRepository) DeleteRecord(u *models.User) (*models.User, error) {
	repo.log.Debugf("deleting record with id %s", u.ID)
	result := repo.db.Delete(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}

func (repo *MysqlRepository) ListRecords() ([]models.User, error) {
	repo.log.Debugf("listing all records")
	var users []models.User
	result := repo.db.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return users, nil
}

func (repo *MysqlRepository) UpdateRecord(u *models.User) (*models.User, error) {
	repo.log.Debugf("updating record with id %d", u.ID)
	result := repo.db.Updates(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}
