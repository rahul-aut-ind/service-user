package userrepo

import (
	"log"

	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DBRepo interface {
		Create(u *models.User) *models.User
		Find(id string) *models.User
	}

	MysqlRepository struct {
		DB  *gorm.DB
		Log *logger.Logger
	}
)

func New(log *logger.Logger, env *config.Env) *MysqlRepository {
	return &MysqlRepository{DB: connect(env.DBConnectionString), Log: log}
}

// Connect initializes the database connection
func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil
	}

	// Auto-migrate the User model
	db.AutoMigrate(&models.User{})

	return db
}

func (db *MysqlRepository) Create(u *models.User) *models.User {
	db.Log.Info("inserting ", u.Email)
	db.DB.Create(&u)
	return u
}

func (db *MysqlRepository) Find(id string) *models.User {
	db.Log.Info("finding user with id ", id)
	var user models.User
	db.DB.Where(id).First(&user)
	return &user
}
