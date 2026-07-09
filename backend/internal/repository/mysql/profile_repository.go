package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/domain"
)

type profileRepository struct {
	db *sql.DB
}

// NewProfileRepository mengembalikan implementasi domain.ProfileRepository yang berbasis MySQL
func NewProfileRepository(db *sql.DB) domain.ProfileRepository {
	return &profileRepository{db: db}
}

// FindAdminProfile murni dari tabel users, tidak ada join sama sekali
func (r *profileRepository) FindAdminProfile(ctx context.Context, userID int64) (*domain.AdminProfile, error) {
	query := `SELECT id, username, role FROM users WHERE id = ? LIMIT 1`

	var p domain.AdminProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&p.ID, &p.Username, &p.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// FindGuruProfile juga murni dari tabel users — belum ada tabel profil guru terpisah di skema V1
func (r *profileRepository) FindGuruProfile(ctx context.Context, userID int64) (*domain.GuruProfile, error) {
	query := `SELECT id, username, role FROM users WHERE id = ? LIMIT 1`

	var p domain.GuruProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&p.ID, &p.Username, &p.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// FindPetugasProfile digabung dari users + staffs (LEFT JOIN supaya tetap dapat data users
// walau baris di staffs belum dibuat admin)
func (r *profileRepository) FindPetugasProfile(ctx context.Context, userID int64) (*domain.PetugasProfile, error) {
	query := `
		SELECT u.id, u.username, u.role, COALESCE(s.nama, ''), COALESCE(s.posisi, ''), (s.id IS NOT NULL) AS has_profile
		FROM users u
		LEFT JOIN staffs s ON s.user_id = u.id
		WHERE u.id = ?
		LIMIT 1
	`

	var p domain.PetugasProfile
	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&p.ID, &p.Username, &p.Role, &p.Nama, &p.Posisi, &p.HasProfile)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// FindSiswaProfile digabung dari users + students + classes (LEFT JOIN supaya tetap dapat data users
// walau baris di students belum dibuat admin)
func (r *profileRepository) FindSiswaProfile(ctx context.Context, userID int64) (*domain.SiswaProfile, error) {
	query := `
		SELECT u.id, u.username, u.role,
		       COALESCE(st.nama, ''), COALESCE(st.nisn, ''), COALESCE(st.alamat, ''), COALESCE(st.no_telp, ''),
		       COALESCE(c.nama_kelas, ''), COALESCE(c.tingkat, 0),
		       (st.id IS NOT NULL) AS has_profile
		FROM users u
		LEFT JOIN students st ON st.user_id = u.id
		LEFT JOIN classes c ON c.id = st.class_id
		WHERE u.id = ?
		LIMIT 1
	`

	var p domain.SiswaProfile
	err := r.db.QueryRowContext(ctx, query, userID).
		Scan(&p.ID, &p.Username, &p.Role, &p.Nama, &p.Nisn, &p.Alamat, &p.NoTelp, &p.NamaKelas, &p.Tingkat, &p.HasProfile)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}
