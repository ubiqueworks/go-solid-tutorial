package repo

import (
	"sync"

	"github.com/ubiqueworks/go-solid-tutorial/domain"
)

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db:   make(map[string]domain.User),
		lock: &sync.RWMutex{},
	}
}

type UserRepository struct {
	db   map[string]domain.User
	lock *sync.RWMutex
}

func (r UserRepository) Create(u *domain.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.db[u.ID] = *u
	return nil
}

func (r UserRepository) DeleteByID(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, exists := r.db[id]; !exists {
		return ErrNotFound
	}

	delete(r.db, id)
	return nil
}

func (r UserRepository) GetByID(id string) (*domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if u, exists := r.db[id]; exists {
		return &u, nil
	}
	return nil, ErrNotFound
}

func (r UserRepository) List() ([]domain.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	users := make([]domain.User, 0, len(r.db))
	for _, user := range r.db {
		users = append(users, user)
	}
	return users, nil
}
