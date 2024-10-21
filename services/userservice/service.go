package userservice

import (
	"fmt"

	"github.com/rahul-aut-ind/service-user/domain/logger"
	"github.com/rahul-aut-ind/service-user/domain/models"
	"github.com/rahul-aut-ind/service-user/interfaceadapters/repositories/userrepo"
)

type (
	IService interface {
		Get(id string) (*models.User, error)
		GetAll() ([]models.User, error)
		Add(u *models.User) (*models.User, error)
		Update(uID string, u *models.User) (*models.User, error)
		Delete(id string) error
	}

	Service struct {
		db  userrepo.DBRepo
		log *logger.Logger
	}
)

func New(r userrepo.DBRepo, l *logger.Logger) *Service {
	return &Service{db: r, log: l}
}

func (s *Service) Add(user *models.User) (*models.User, error) {
	res, err := s.db.CreateRecord(user)
	if err != nil {
		msg := fmt.Sprintf("error creating user :: %s", err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}
	return res, nil
}

func (s *Service) Get(id string) (*models.User, error) {
	res, err := s.db.FindRecord(id)
	if err != nil {
		msg := fmt.Sprintf("error %s :: %s", models.ErrMsgNoUserfound, err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}

func (s *Service) Delete(id string) error {
	res, err := s.db.FindRecord(id)
	if err != nil {
		msg := fmt.Sprintf("error %s :: %s", models.ErrMsgNoUserfound, err.Error())
		s.log.Errorf(msg)
		return fmt.Errorf("%s", msg)
	}

	res, err = s.db.DeleteRecord(res)
	if err != nil {
		msg := fmt.Sprintf("error deleting user %d :: %s", res.ID, err.Error())
		s.log.Errorf(msg)
		return fmt.Errorf("%s", msg)
	}

	return nil
}

func (s *Service) GetAll() ([]models.User, error) {
	res, err := s.db.ListRecords()
	if err != nil {
		msg := fmt.Sprintf("error getting all users :: %s", err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}

func (s *Service) Update(uID string, u *models.User) (*models.User, error) {
	rec, err := s.db.FindRecord(uID)
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
		msg := fmt.Sprintf("error updating user %d :: %s", res.ID, err.Error())
		s.log.Errorf(msg)
		return nil, fmt.Errorf("%s", msg)
	}

	return res, nil
}
