package repository

import (
	"sync"

	"restapi/internal/models"
)

// MemoryStore is the in-process Store used by default until
// postgres is wired. Full CRUD lands in the next phase; the
// plumbing (constructor + interface satisfaction) lands now
// so main -> HandleConnection threading compiles.
type MemoryStore struct {
	mu     sync.RWMutex
	users  map[int]models.User
	nextU  int
	emails map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:  make(map[int]models.User),
		nextU:  1,
		emails: make(map[string]int),
	}
}

func (m *MemoryStore) CreateUser(name, email string) (models.User, error) {
	return models.User{}, errNotImplemented
}

func (m *MemoryStore) GetUser(id int) (models.User, error) {
	return models.User{}, errNotImplemented
}

func (m *MemoryStore) ListUsers() ([]models.User, error) {
	return nil, errNotImplemented
}

func (m *MemoryStore) UpdateUser(id int, name, email string) (models.User, error) {
	return models.User{}, errNotImplemented
}

func (m *MemoryStore) PatchUser(id int, name, email *string) (models.User, error){
	return models.User{}, errNotImplemented
}

func (m *MemoryStore) DeleteUser(id int) error { return errNotImplemented }
