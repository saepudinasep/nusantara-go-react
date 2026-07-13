package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"backend/internal/domain"
)

type studentRepository struct {
	db *sql.DB
}

// NewStudentRepository mengembalikan implementasi domain.StudentRepository yang berbasis MySQL
func NewStudentRepository(db *sql.DB) domain.StudentRepository {
	return &studentRepository{db: db}
}

// Create menjalankan 2 insert (users lalu students) dalam SATU transaksi.
// Kalau insert ke students gagal (mis. NISN duplikat, class_id tidak valid), insert user yang
// baru saja dibuat ikut di-rollback — supaya tidak ada akun login "nyangkut" tanpa profil siswa.
func (r *studentRepository) Create(ctx context.Context, s *domain.Student, hashedPassword string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // no-op kalau sudah di-Commit sebelumnya

	// 1. Insert akun login ke tabel users
	userResult, err := tx.ExecContext(ctx,
		`INSERT INTO users (username, password, role) VALUES (?, ?, 'siswa')`,
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

	// 2. Insert profil siswa, terhubung ke user yang baru saja dibuat
	studentResult, err := tx.ExecContext(ctx,
		`INSERT INTO students (user_id, class_id, nisn, nama, alamat, no_telp) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, s.ClassID, s.Nisn, s.Nama, s.Alamat, s.NoTelp,
	)
	if err != nil {
		if isDuplicateEntryError(err) {
			return 0, domain.ErrNisnTaken
		}
		if isForeignKeyChildRowError(err) {
			return 0, domain.ErrKelasTidakValid
		}
		if isLockWaitTimeoutError(err) {
			return 0, domain.ErrDatabaseBusy
		}
		return 0, err
	}

	studentID, err := studentResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return studentID, nil
}

func (r *studentRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Student, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM students`).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT s.id, s.user_id, u.username, s.nisn, s.nama, s.class_id, c.nama_kelas, c.tingkat,
		       s.alamat, s.no_telp, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN classes c ON c.id = s.class_id
		ORDER BY s.nama ASC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []domain.Student{}
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Nisn, &s.Nama, &s.ClassID, &s.NamaKelas, &s.Tingkat,
			&s.Alamat, &s.NoTelp, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}

	return list, total, rows.Err()
}

func (r *studentRepository) FindByID(ctx context.Context, id int64) (*domain.Student, error) {
	query := `
		SELECT s.id, s.user_id, u.username, s.nisn, s.nama, s.class_id, c.nama_kelas, c.tingkat,
		       s.alamat, s.no_telp, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN classes c ON c.id = s.class_id
		WHERE s.id = ?
		LIMIT 1
	`

	var s domain.Student
	err := r.db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.UserID, &s.Username, &s.Nisn, &s.Nama, &s.ClassID,
		&s.NamaKelas, &s.Tingkat, &s.Alamat, &s.NoTelp, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSiswaNotFound
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// FindByNisn dipakai fitur "cari siswa by NISN" pada proses pembayaran (admin & petugas)
func (r *studentRepository) FindByNisn(ctx context.Context, nisn string) (*domain.Student, error) {
	query := `
		SELECT s.id, s.user_id, u.username, s.nisn, s.nama, s.class_id, c.nama_kelas, c.tingkat,
		       s.alamat, s.no_telp, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN classes c ON c.id = s.class_id
		WHERE s.nisn = ?
		LIMIT 1
	`

	var s domain.Student
	err := r.db.QueryRowContext(ctx, query, nisn).Scan(&s.ID, &s.UserID, &s.Username, &s.Nisn, &s.Nama, &s.ClassID,
		&s.NamaKelas, &s.Tingkat, &s.Alamat, &s.NoTelp, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrStudentInvalid
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// FindByUserID dipakai siswa yang sedang login untuk melihat profil/tagihan/riwayat miliknya sendiri
func (r *studentRepository) FindByUserID(ctx context.Context, userID int64) (*domain.Student, error) {
	query := `
		SELECT s.id, s.user_id, u.username, s.nisn, s.nama, s.class_id, c.nama_kelas, c.tingkat,
		       s.alamat, s.no_telp, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		JOIN classes c ON c.id = s.class_id
		WHERE s.user_id = ?
		LIMIT 1
	`

	var su domain.Student
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&su.ID, &su.UserID, &su.Username, &su.Nisn, &su.Nama, &su.ClassID,
		&su.NamaKelas, &su.Tingkat, &su.Alamat, &su.NoTelp, &su.CreatedAt, &su.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrStudentProfileMissing
	}
	if err != nil {
		return nil, err
	}

	return &su, nil
}

// Update hanya mengubah kolom di tabel students (nisn, nama, class_id, alamat, no_telp).
// Username/password login TIDAK diubah lewat endpoint ini.
func (r *studentRepository) Update(ctx context.Context, s *domain.Student) error {
	query := `UPDATE students SET nisn = ?, nama = ?, class_id = ?, alamat = ?, no_telp = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, s.Nisn, s.Nama, s.ClassID, s.Alamat, s.NoTelp, s.ID)
	if err != nil {
		if isDuplicateEntryError(err) {
			return domain.ErrNisnTaken
		}
		if isForeignKeyChildRowError(err) {
			return domain.ErrKelasTidakValid
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
		return domain.ErrSiswaNotFound
	}

	return nil
}

// Delete menghapus lewat tabel users (bukan langsung DELETE FROM students), supaya akun login
// siswa ikut terhapus otomatis lewat ON DELETE CASCADE di foreign key students.user_id.
// Kalau siswa masih punya riwayat di tabel payments (FK RESTRICT), MySQL akan menolak — kita
// petakan jadi domain.ErrSiswaInUse, bukan 500 bocoran SQL.
func (r *studentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE u FROM users u JOIN students s ON s.user_id = u.id WHERE s.id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if isForeignKeyConstraintError(err) {
			return domain.ErrSiswaInUse
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
		return domain.ErrSiswaNotFound
	}

	return nil
}

// isForeignKeyChildRowError mendeteksi error MySQL 1452 (baris yang di-insert/update mereferensikan
// parent yang tidak ada, mis. class_id yang salah) — beda kasus dengan 1451 (isForeignKeyConstraintError,
// dipakai untuk DELETE yang ditolak karena masih direferensikan baris lain).
func isForeignKeyChildRowError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Cannot add or update a child row")
}
