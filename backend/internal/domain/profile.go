package domain

import "context"

// AdminProfile diambil murni dari tabel users (admin tidak punya tabel profil terpisah)
type AdminProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// GuruProfile — CATATAN: skema V1 belum punya tabel khusus untuk guru (mis. "teachers"),
// jadi untuk saat ini datanya sama seperti AdminProfile (murni dari tabel users).
// Rencananya tabel profil guru akan ditambahkan di V2 (Learning Center).
type GuruProfile struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// PetugasProfile digabung dari tabel users + staffs
type PetugasProfile struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Nama       string `json:"nama"`
	Posisi     string `json:"posisi"`
	HasProfile bool   `json:"has_profile"` // false kalau user role=petugas tapi belum ada baris di tabel staffs
}

// SiswaProfile digabung dari tabel users + students + classes
type SiswaProfile struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	Nama       string `json:"nama"`
	Nisn       string `json:"nisn"`
	Alamat     string `json:"alamat"`
	NoTelp     string `json:"no_telp"`
	NamaKelas  string `json:"nama_kelas"`
	Tingkat    int    `json:"tingkat"`
	HasProfile bool   `json:"has_profile"` // false kalau user role=siswa tapi belum ada baris di tabel students
}

// ProfileRepository adalah interface (port) untuk mengambil data profil per role dari MySQL
type ProfileRepository interface {
	FindAdminProfile(ctx context.Context, userID int64) (*AdminProfile, error)
	FindGuruProfile(ctx context.Context, userID int64) (*GuruProfile, error)
	FindPetugasProfile(ctx context.Context, userID int64) (*PetugasProfile, error)
	FindSiswaProfile(ctx context.Context, userID int64) (*SiswaProfile, error)
}

// ProfileUsecase adalah interface (port) untuk business logic pengambilan profil per role
type ProfileUsecase interface {
	GetAdminProfile(ctx context.Context, userID int64) (*AdminProfile, error)
	GetGuruProfile(ctx context.Context, userID int64) (*GuruProfile, error)
	GetPetugasProfile(ctx context.Context, userID int64) (*PetugasProfile, error)
	GetSiswaProfile(ctx context.Context, userID int64) (*SiswaProfile, error)
}
