package repository

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"restapi/internal/models"
)

type PostgresStore struct {
	db *sql.DB
}

// Compile-time guarantee that *PostgresStore satisfies Store.
var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func NewPostgresDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func (s *PostgresStore) ListUsers() ([]models.User, error) {
	rows, err := s.db.Query(
		`SELECT id, name, email, created_at
		FROM users
		ORDER BY id
		`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var user models.User

		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *PostgresStore) GetUser(id int) (models.User, error) {
	var user models.User

	err := s.db.QueryRow(`
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}


func (s *PostgresStore) CreateUser(name, email string) (models.User, error) {
	return models.User{}, errNotImplemented
}

func (s *PostgresStore) UpdateUser(id int, name, email string) (models.User, error) {
	return models.User{}, errNotImplemented
}

func (s *PostgresStore) DeleteUser(id int) error { return errNotImplemented }

func (s *PostgresStore) CreateTodo(userID int, title string) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (s *PostgresStore) GetTodo(id int) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (s *PostgresStore) ListTodos(userID *int, completed *bool) ([]models.Todo, error) {
	return nil, errNotImplemented
}

func (s *PostgresStore) UpdateTodoFull(id int, title string, completed bool) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (s *PostgresStore) PatchTodo(id int, title *string, completed *bool) (models.Todo, error) {
	return models.Todo{}, errNotImplemented
}

func (s *PostgresStore) DeleteTodo(id int) error { return errNotImplemented }
