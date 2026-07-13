package domain

import "context"

// Student merepresentasikan entity siswa. Field Username/Password hanya dipakai saat Create
// (untuk sekaligus membuat akun login di tabel users) — Password tidak pernah dikembalikan
// ke client (json:"-"), dan Username diabaikan saat Update (kredensial login tidak diubah lewat form ini).
type Student struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Nisn      string `json:"nisn"`
	Nama      string `json:"nama"`
	ClassID   int64  `json:"class_id"`
	NamaKelas string `json:"nama_kelas,omitempty"`
	Tingkat   int    `json:"tingkat,omitempty"`
	Alamat    string `json:"alamat"`
	NoTelp    string `json:"no_telp"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// StudentRepository adalah interface (port) yang harus diimplementasikan oleh layer repository.
// Create WAJIB dijalankan dalam satu transaksi DB (insert users + insert students, all-or-nothing).
type StudentRepository interface {
	Create(ctx context.Context, s *Student, hashedPassword string) (int64, error)
	FindAll(ctx context.Context, page, limit int) ([]Student, int64, error)
	FindByID(ctx context.Context, id int64) (*Student, error)
	FindByNisn(ctx context.Context, nisn string) (*Student, error)
	// FindByUserID dipakai siswa yang sedang login untuk melihat profil/tagihan/riwayat miliknya sendiri
	FindByUserID(ctx context.Context, userID int64) (*Student, error)
	Update(ctx context.Context, s *Student) error
	Delete(ctx context.Context, id int64) error
}

// StudentUsecase adalah interface (port) untuk business logic pengelolaan siswa
type StudentUsecase interface {
	Create(ctx context.Context, s *Student) (*Student, error)
	GetAll(ctx context.Context, page, limit int) ([]Student, Pagination, error)
	GetByID(ctx context.Context, id int64) (*Student, error)
	SearchByNisn(ctx context.Context, nisn string) (*Student, error)
	GetOwnProfile(ctx context.Context, userID int64) (*Student, error)
	Update(ctx context.Context, id int64, s *Student) (*Student, error)
	Delete(ctx context.Context, id int64) error
}
