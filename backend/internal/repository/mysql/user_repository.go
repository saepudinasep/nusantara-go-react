package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository mengembalikan implementasi domain.UserRepository yang berbasis MySQL
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, password, role FROM users WHERE username = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, username)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, password, role FROM users WHERE id = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
