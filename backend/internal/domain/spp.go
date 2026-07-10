package domain

import "context"

// Spp merepresentasikan entity master SPP (tabel fisik: "spp", sesuai skema V1)
type Spp struct {
	ID          int64   `json:"id"`
	TahunAjaran string  `json:"tahun_ajaran"`
	Nominal     float64 `json:"nominal"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

// SppRepository adalah interface (port) yang harus diimplementasikan oleh layer repository
type SppRepository interface {
	Create(ctx context.Context, s *Spp) (int64, error)
	FindAll(ctx context.Context, page, limit int) ([]Spp, int64, error)
	FindByID(ctx context.Context, id int64) (*Spp, error)
	Update(ctx context.Context, s *Spp) error
	Delete(ctx context.Context, id int64) error
}

// SppUsecase adalah interface (port) untuk business logic pengelolaan SPP
type SppUsecase interface {
	Create(ctx context.Context, s *Spp) (*Spp, error)
	GetAll(ctx context.Context, page, limit int) ([]Spp, Pagination, error)
	GetByID(ctx context.Context, id int64) (*Spp, error)
	Update(ctx context.Context, id int64, s *Spp) (*Spp, error)
	Delete(ctx context.Context, id int64) error
}
