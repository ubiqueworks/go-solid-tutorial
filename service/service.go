//go:generate mockgen -package service -source=service.go -destination service_mock.go

package service

import (
	"github.com/google/uuid"
	"github.com/ubiqueworks/go-interface-usage/domain"
	"github.com/ubiqueworks/go-interface-usage/repo"
)

type Service interface {
	CreateUser(u *domain.User) error
	DeleteUser(userId string) error
	ListUsers() ([]domain.User, error)
}

func New(userRepo repo.UserRepository) (Service, error) {
	return &serviceImpl{
		userRepo: userRepo,
	}, nil
}

type serviceImpl struct {
	userRepo repo.UserRepository
}

func (s serviceImpl) CreateUser(u *domain.User) error {
	if u == nil {
		return ErrInvalidData
	}

	if u.Name == "" {
		return ErrInvalidData
	}

	u.ID = uuid.New().String()
	if err := s.userRepo.Create(u); err != nil {
		return err
	}
	return nil
}

func (s serviceImpl) DeleteUser(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidData
	}

	if err := s.userRepo.DeleteByID(uid.String()); err != nil {
		return err
	}
	return nil
}

func (s serviceImpl) ListUsers() ([]domain.User, error) {
	return s.userRepo.List()
}
