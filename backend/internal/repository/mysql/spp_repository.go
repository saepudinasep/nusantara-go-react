package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type sppRepository struct {
	db *sql.DB
}

// NewSppRepository mengembalikan implementasi domain.SppRepository yang berbasis MySQL
func NewSppRepository(db *sql.DB) domain.SppRepository {
	return &sppRepository{db: db}
}

func (r *sppRepository) Create(ctx context.Context, s *domain.Spp) (int64, error) {
	query := `INSERT INTO spp (tahun_ajaran, nominal) VALUES (?, ?)`

	result, err := r.db.ExecContext(ctx, query, s.TahunAjaran, s.Nominal)
	if err != nil {
		if isDuplicateEntryError(err) {
			return 0, domain.ErrSppDuplicate
		}
		if isLockWaitTimeoutError(err) {
			return 0, domain.ErrDatabaseBusy
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (r *sppRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Spp, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM spp`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `SELECT id, tahun_ajaran, nominal, created_at, updated_at
	          FROM spp ORDER BY tahun_ajaran DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []domain.Spp{}
	for rows.Next() {
		var s domain.Spp
		if err := rows.Scan(&s.ID, &s.TahunAjaran, &s.Nominal, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}

	return list, total, rows.Err()
}

func (r *sppRepository) FindByID(ctx context.Context, id int64) (*domain.Spp, error) {
	query := `SELECT id, tahun_ajaran, nominal, created_at, updated_at FROM spp WHERE id = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, id)

	var s domain.Spp
	err := row.Scan(&s.ID, &s.TahunAjaran, &s.Nominal, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSppNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *sppRepository) Update(ctx context.Context, s *domain.Spp) error {
	query := `UPDATE spp SET tahun_ajaran = ?, nominal = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, s.TahunAjaran, s.Nominal, s.ID)
	if err != nil {
		if isDuplicateEntryError(err) {
			return domain.ErrSppDuplicate
		}
		if isLockWaitTimeoutError(err) {
			return domain.ErrDatabaseBusy
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSppNotFound
	}

	return nil
}

func (r *sppRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM spp WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if isForeignKeyConstraintError(err) {
			return domain.ErrSppInUse
		}
		if isLockWaitTimeoutError(err) {
			return domain.ErrDatabaseBusy
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrSppNotFound
	}

	return nil
}
