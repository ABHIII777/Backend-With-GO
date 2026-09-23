package repository

import (
	"errors"

	"restapi/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

// Store is implemented by memory and postgres backends.
// Services and router depend only on this interface.
type Store interface {
	CreateUser(name, email string) (models.User, error)
	GetUser(id int) (models.User, error)
	ListUsers() ([]models.User, error)
	UpdateUser(id int, name, email string) (models.User, error)
	PatchUser(id int, name, email *string) (models.User, error)
	DeleteUser(id int) error

	CreateTodo(userID int, title string) (models.Todo, error)
	GetTodo(id int) (models.Todo, error)
	ListTodos(userID *int, completed *bool) ([]models.Todo, error)
	UpdateTodoFull(id int, title string, completed bool) (models.Todo, error)
	PatchTodo(id int, title *string, completed *bool) (models.Todo, error)
	DeleteTodo(id int) error
}
