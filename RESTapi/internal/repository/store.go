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
}
