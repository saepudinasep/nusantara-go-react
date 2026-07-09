package domain

import "context"

// Kelas merepresentasikan entity kelas (tabel fisik di database: "classes", sesuai skema V1)
type Kelas struct {
	ID        int64  `json:"id"`
	NamaKelas string `json:"nama_kelas"`
	Tingkat   int    `json:"tingkat"` // contoh: 10, 11, 12
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Pagination merepresentasikan parameter & hasil pagination
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// KelasRepository adalah interface (port) yang harus diimplementasikan oleh layer repository
type KelasRepository interface {
	Create(ctx context.Context, k *Kelas) (int64, error)
	FindAll(ctx context.Context, page, limit int) ([]Kelas, int64, error)
	FindByID(ctx context.Context, id int64) (*Kelas, error)
	Update(ctx context.Context, k *Kelas) error
	Delete(ctx context.Context, id int64) error
}

// KelasUsecase adalah interface (port) untuk business logic pengelolaan kelas
type KelasUsecase interface {
	Create(ctx context.Context, k *Kelas) (*Kelas, error)
	GetAll(ctx context.Context, page, limit int) ([]Kelas, Pagination, error)
	GetByID(ctx context.Context, id int64) (*Kelas, error)
	Update(ctx context.Context, id int64, k *Kelas) (*Kelas, error)
	Delete(ctx context.Context, id int64) error
}
