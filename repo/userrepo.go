//go:generate mockgen -package repo -source=userrepo.go -destination userrepo_mock.go

package repo

import (
	"sync"

	"github.com/ubiqueworks/go-interface-usage/domain"
)

type UserRepository interface {
	Create(u *domain.User) error
	DeleteByID(id string) error
	GetByID(id string) (*domain.User, error)
	List() ([]domain.User, error)
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db:   make(map[string]domain.User),
		lock: &sync.RWMutex{},
	}
}

type userRepository struct {
	db   map[string]domain.User
	lock *sync.RWMutex
}

func (r userRepository) Create(u *domain.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.db[u.ID] = *u
	return nil
}

func (r userRepository) DeleteByID(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, exists := r.db[id]; !exists {
		return ErrNotFound
	}

	delete(r.db, id)
	return nil
}

func (r userRepository) GetByID(id string) (*domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if u, exists := r.db[id]; exists {
		return &u, nil
	}
	return nil, ErrNotFound
}

func (r userRepository) List() ([]domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	users := make([]domain.User, len(r.db))
	for _, user := range r.db {
		users = append(users, user)
	}
	return users, nil
}
