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
	"gorm.io/gorm/schema"
)

type (
	DataHandler interface {
		ListRecords() ([]models.User, error)
		FindRecord(id string) (*models.User, error)
		CreateRecord(u *models.User) (*models.User, error)
		UpdateRecord(u *models.User) (*models.User, error)
		DeleteRecord(u *models.User) (*models.User, error)
	}

	MysqlClient struct {
		client *gorm.DB
		log    *logger.Logger
	}
)

const (
	LogLevel = glog.Warn
)

func New(l *logger.Logger, env *config.Env) *MysqlClient {
	return &MysqlClient{client: connect(env.DBConnectionString), log: l}
}

// Connect initializes the database connection
func connect(dsn string) *gorm.DB {
	client, err := gorm.Open(mysql.Open(dsn), initConfig())
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database :: %v", err))
	}

	// Auto-migrate the User model
	err = client.AutoMigrate(&models.User{})
	if err != nil {
		panic(fmt.Sprintf("could not initialize tables | err :: %v", err))
	}

	sqlDB, _ := client.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(5)

	return client
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
		LogLevel:      LogLevel,
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

func (db *MysqlClient) CreateRecord(u *models.User) (*models.User, error) {
	db.log.Debugf("inserting record %v", *u)
	result := db.client.Create(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}

func (db *MysqlClient) FindRecord(id string) (*models.User, error) {
	db.log.Debugf("finding record with id %s", id)
	var user models.User
	result := db.client.Where(id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("err :: %v", errors.ErrCodeNoUser)
		}
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return &user, nil
}

func (db *MysqlClient) DeleteRecord(u *models.User) (*models.User, error) {
	db.log.Debugf("deleting record with id %d", u.ID)
	result := db.client.Delete(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}

func (db *MysqlClient) ListRecords() ([]models.User, error) {
	db.log.Debugf("listing all records")
	var users []models.User
	result := db.client.Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return users, nil
}

func (db *MysqlClient) UpdateRecord(u *models.User) (*models.User, error) {
	db.log.Debugf("updating record with id %d", u.ID)
	result := db.client.Updates(&u)
	if result.Error != nil {
		return nil, fmt.Errorf("err :: %v", result.Error)
	}
	return u, nil
}
