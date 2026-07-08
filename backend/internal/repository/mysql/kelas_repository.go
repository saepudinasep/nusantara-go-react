package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"backend/internal/domain"
)

type kelasRepository struct {
	db *sql.DB
}

// NewKelasRepository mengembalikan implementasi domain.KelasRepository yang berbasis MySQL,
// menyasar tabel fisik "classes" sesuai skema V1.
func NewKelasRepository(db *sql.DB) domain.KelasRepository {
	return &kelasRepository{db: db}
}

func (r *kelasRepository) Create(ctx context.Context, k *domain.Kelas) (int64, error) {
	query := `INSERT INTO classes (nama_kelas, tingkat) VALUES (?, ?)`

	result, err := r.db.ExecContext(ctx, query, k.NamaKelas, k.Tingkat)
	if err != nil {
		if isDuplicateEntryError(err) {
			return 0, domain.ErrDuplicateEntry
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (r *kelasRepository) FindAll(ctx context.Context) ([]domain.Kelas, error) {
	query := `SELECT id, nama_kelas, tingkat, created_at, updated_at FROM classes ORDER BY nama_kelas ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Kelas
	for rows.Next() {
		var k domain.Kelas
		if err := rows.Scan(&k.ID, &k.NamaKelas, &k.Tingkat, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, k)
	}

	return list, rows.Err()
}

func (r *kelasRepository) FindByID(ctx context.Context, id int64) (*domain.Kelas, error) {
	query := `SELECT id, nama_kelas, tingkat, created_at, updated_at FROM classes WHERE id = ? LIMIT 1`

	row := r.db.QueryRowContext(ctx, query, id)

	var k domain.Kelas
	err := row.Scan(&k.ID, &k.NamaKelas, &k.Tingkat, &k.CreatedAt, &k.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrKelasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func (r *kelasRepository) Update(ctx context.Context, k *domain.Kelas) error {
	query := `UPDATE classes SET nama_kelas = ?, tingkat = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, k.NamaKelas, k.Tingkat, k.ID)
	if err != nil {
		if isDuplicateEntryError(err) {
			return domain.ErrDuplicateEntry
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrKelasNotFound
	}

	return nil
}

func (r *kelasRepository) Delete(ctx context.Context, id int64) error {
	// Tabel "classes" pada skema V1 tidak punya kolom deleted_at, jadi ini hard delete.
	query := `DELETE FROM classes WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrKelasNotFound
	}

	return nil
}

// isDuplicateEntryError mendeteksi error MySQL 1062 (Duplicate entry) tanpa perlu import driver spesifik
func isDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
