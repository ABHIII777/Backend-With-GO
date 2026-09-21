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
	todos  map[int]models.Todo
	nextU  int
	nextT  int
	emails map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:  make(map[int]models.User),
		todos:  make(map[int]models.Todo),
		nextU:  1,
		nextT:  1,
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

func (m *MemoryStore) DeleteUser(id int) error { return errNotImplemented }

func (m *MemoryStore) CreateTodo(userID int, title string) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (m *MemoryStore) GetTodo(id int) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (m *MemoryStore) ListTodos(userID *int, completed *bool) ([]models.Todo, error) {
	return nil, errNotImplemented
}

func (m *MemoryStore) UpdateTodoFull(id int, title string, completed bool) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (m *MemoryStore) PatchTodo(id int, title *string, completed *bool) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (m *MemoryStore) DeleteTodo(id int) error { return errNotImplemented }
