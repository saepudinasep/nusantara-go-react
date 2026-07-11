package domain

import "context"

// Staff merepresentasikan entity petugas SPP. Field Username/Password hanya dipakai saat Create
// (sekaligus membuat akun login role=petugas di tabel users) — Password tidak pernah dikembalikan
// ke client (json:"-"), dan Username diabaikan saat Update.
type Staff struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Nama      string `json:"nama"`
	Posisi    string `json:"posisi"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// StaffRepository adalah interface (port) yang harus diimplementasikan oleh layer repository.
// Create WAJIB dijalankan dalam satu transaksi DB (insert users + insert staffs, all-or-nothing).
type StaffRepository interface {
	Create(ctx context.Context, s *Staff, hashedPassword string) (int64, error)
	FindAll(ctx context.Context, page, limit int) ([]Staff, int64, error)
	FindByID(ctx context.Context, id int64) (*Staff, error)
	Update(ctx context.Context, s *Staff) error
	Delete(ctx context.Context, id int64) error
}

// StaffUsecase adalah interface (port) untuk business logic pengelolaan petugas
type StaffUsecase interface {
	Create(ctx context.Context, s *Staff) (*Staff, error)
	GetAll(ctx context.Context, page, limit int) ([]Staff, Pagination, error)
	GetByID(ctx context.Context, id int64) (*Staff, error)
	Update(ctx context.Context, id int64, s *Staff) (*Staff, error)
	Delete(ctx context.Context, id int64) error
}
