package userservice

import (
	"fmt"

	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/userrepo"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type (
	Services interface {
		GetUserWithID(id string) (*models.User, error)
		GetAllUsers() ([]models.User, error)
		AddUser(u *models.User) (*models.User, error)
		UpdateUser(id string, u *models.User) (*models.User, error)
		DeleteUser(id string) error
		UploadProfilePicture(id string) error
	}

	Service struct {
		db  userrepo.DBRepo
		log *logger.Logger
	}
)

func New(r userrepo.DBRepo, l *logger.Logger) *Service {
	return &Service{db: r, log: l}
}

func (s *Service) AddUser(user *models.User) (*models.User, error) {
	res, err := s.db.CreateRecord(user)
	if err != nil {
		msg := fmt.Sprintf("error creating user :: %s", err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}
	return res, nil
}

func (s *Service) GetUserWithID(id string) (*models.User, error) {
	res, err := s.db.FindRecord(id)
	if err != nil {
		msg := fmt.Sprintf("error :: %s", err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}

func (s *Service) DeleteUser(id string) error {
	res, err := s.db.FindRecord(id)
	if err != nil {
		msg := fmt.Sprintf("error %s :: %s", models.ErrMsgNoUserfound, err.Error())
		s.log.Errorf(msg)
		return fmt.Errorf("%s", msg)
	}

	res, err = s.db.DeleteRecord(res)
	if err != nil {
		msg := fmt.Sprintf("error deleting user %s :: %s", id, err.Error())
		s.log.Errorf(msg)
		return fmt.Errorf("%s", msg)
	}
	s.log.Debugf("deleted user %d", res.ID)

	return nil
}

func (s *Service) GetAllUsers() ([]models.User, error) {
	res, err := s.db.ListRecords()
	if err != nil {
		msg := fmt.Sprintf("error getting all users :: %s", err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}

func (s *Service) UpdateUser(id string, u *models.User) (*models.User, error) {
	rec, err := s.db.FindRecord(id)
	if err != nil {
		msg := fmt.Sprintf("error %s :: %s", models.ErrMsgNoUserfound, err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	// primary key mapping
	u.ID = rec.ID
	// updating of email should be prohibited, maybe is used for login
	u.Email = rec.Email

	res, err := s.db.UpdateRecord(u)
	if err != nil {
		msg := fmt.Sprintf("error updating user %s :: %s", id, err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}

func (s *Service) UploadProfilePicture(id string) error {
	s.log.Debugf("uploading profile pic with id %s", id)
	return nil
}
