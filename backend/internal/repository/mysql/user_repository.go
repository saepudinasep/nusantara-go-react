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

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE email = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email, password, role FROM users WHERE id = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, id)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
