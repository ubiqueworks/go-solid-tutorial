//go:generate mockgen -package service -source=service.go -destination service_mock.go

package service

import (
	"github.com/google/uuid"
	"github.com/ubiqueworks/go-solid-tutorial/domain"
)

type UserRepository interface {
	Create(u *domain.User) error
	DeleteByID(id string) error
	GetByID(id string) (*domain.User, error)
	List() ([]domain.User, error)
}

func New(userRepo UserRepository) (*Service, error) {
	return &Service{
		userRepo: userRepo,
	}, nil
}

type Service struct {
	userRepo UserRepository
}

func (s Service) CreateUser(u *domain.User) error {
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

func (s Service) DeleteUser(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidData
	}

	if err := s.userRepo.DeleteByID(uid.String()); err != nil {
		return err
	}
	return nil
}

func (s Service) ListUsers() ([]domain.User, error) {
	return s.userRepo.List()
}
