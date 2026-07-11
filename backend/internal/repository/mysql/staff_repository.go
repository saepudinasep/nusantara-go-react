package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type staffRepository struct {
	db *sql.DB
}

// NewStaffRepository mengembalikan implementasi domain.StaffRepository yang berbasis MySQL
func NewStaffRepository(db *sql.DB) domain.StaffRepository {
	return &staffRepository{db: db}
}

// Create menjalankan 2 insert (users lalu staffs) dalam SATU transaksi.
// Kalau insert ke staffs gagal, insert user yang baru saja dibuat ikut di-rollback —
// supaya tidak ada akun login "nyangkut" tanpa profil petugas.
func (r *staffRepository) Create(ctx context.Context, s *domain.Staff, hashedPassword string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // no-op kalau sudah di-Commit sebelumnya

	// 1. Insert akun login ke tabel users
	userResult, err := tx.ExecContext(ctx,
		`INSERT INTO users (username, password, role) VALUES (?, ?, 'petugas')`,
		s.Username, hashedPassword,
	)
	if err != nil {
		if isDuplicateEntryError(err) {
			return 0, domain.ErrUsernameTaken
		}
		if isLockWaitTimeoutError(err) {
			return 0, domain.ErrDatabaseBusy
		}
		return 0, err
	}

	userID, err := userResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 2. Insert profil petugas, terhubung ke user yang baru saja dibuat
	staffResult, err := tx.ExecContext(ctx,
		`INSERT INTO staffs (user_id, nama, posisi) VALUES (?, ?, ?)`,
		userID, s.Nama, s.Posisi,
	)
	if err != nil {
		if isLockWaitTimeoutError(err) {
			return 0, domain.ErrDatabaseBusy
		}
		return 0, err
	}

	staffID, err := staffResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return staffID, nil
}

func (r *staffRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Staff, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM staffs`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT s.id, s.user_id, u.username, s.nama, s.posisi, s.created_at, s.updated_at
		FROM staffs s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.nama ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []domain.Staff{}
	for rows.Next() {
		var s domain.Staff
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Nama, &s.Posisi, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}

	return list, total, rows.Err()
}

func (r *staffRepository) FindByID(ctx context.Context, id int64) (*domain.Staff, error) {
	query := `
		SELECT s.id, s.user_id, u.username, s.nama, s.posisi, s.created_at, s.updated_at
		FROM staffs s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = ?
		LIMIT 1
	`

	var s domain.Staff
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&s.ID, &s.UserID, &s.Username, &s.Nama, &s.Posisi, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPetugasNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Update hanya mengubah kolom di tabel staffs (nama, posisi).
// Username/password login TIDAK diubah lewat endpoint ini.
func (r *staffRepository) Update(ctx context.Context, s *domain.Staff) error {
	query := `UPDATE staffs SET nama = ?, posisi = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, s.Nama, s.Posisi, s.ID)
	if err != nil {
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
		return domain.ErrPetugasNotFound
	}

	return nil
}

// Delete menghapus lewat tabel users (bukan langsung DELETE FROM staffs), supaya akun login
// petugas ikut terhapus otomatis lewat ON DELETE CASCADE di foreign key staffs.user_id.
// Kalau petugas masih punya riwayat di tabel payments (FK RESTRICT), MySQL akan menolak —
// kita petakan jadi domain.ErrPetugasInUse, bukan 500 bocoran SQL.
func (r *staffRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE u FROM users u JOIN staffs s ON s.user_id = u.id WHERE s.id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if isForeignKeyConstraintError(err) {
			return domain.ErrPetugasInUse
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
		return domain.ErrPetugasNotFound
	}

	return nil
}
