package userservice

import (
	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/userrepo"
)

type (
	IService interface {
		Find(id string) (*models.User, error)
		// FindAll() ([]models.User, error)
		Add(u *models.User) (*models.User, error)
		// Update(u *models.User) (*models.User, error)
		// Delete(id string) error
	}
	Service struct {
		db  userrepo.DBRepo
		Log *logger.Logger
	}
)

func New(r userrepo.DBRepo, l *logger.Logger) *Service {
	return &Service{db: r, Log: l}
}

func (s *Service) Add(user *models.User) (*models.User, error) {
	res := s.db.Create(user)
	s.Log.Infof("service %v", res)
	return user, nil
}

func (s *Service) Find(id string) (*models.User, error) {
	res := s.db.Find(id)
	s.Log.Infof("service %v", res)
	return res, nil
}
