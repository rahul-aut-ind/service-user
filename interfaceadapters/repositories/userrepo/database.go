package userrepo

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rahul-aut-ind/service-user/domain/errors"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	schema "gorm.io/gorm/schema"
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
	db, err := gorm.Open(mysql.Open(dsn), initConfig())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database :: %v", err))
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(fmt.Sprintf("could not initialize tables | err :: %v", err))
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(5)

	return db
}

// initConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

// initLog Connection Log Configuration
func initLog() glog.Interface {
	newLogger := glog.New(log.New(os.Stdout, "\r\n", log.LstdFlags), glog.Config{
		Colorful:      true,
		LogLevel:      glog.Warn,
		SlowThreshold: time.Second,
	})
	return newLogger
}

// initNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
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
			return nil, fmt.Errorf("err :: %v", errors.ErrCodeNoUser)
		}
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return &user, nil
}

func (repo *MysqlRepository) DeleteRecord(u *models.User) (*models.User, error) {
	repo.log.Debugf("deleting record with id %d", u.ID)
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
