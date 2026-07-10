package usecase

import (
	"context"
	"errors"
	"regexp"

	"backend/internal/domain"
)

type sppUsecase struct {
	sppRepo domain.SppRepository
}

// NewSppUsecase mengembalikan implementasi domain.SppUsecase
func NewSppUsecase(sppRepo domain.SppRepository) domain.SppUsecase {
	return &sppUsecase{sppRepo: sppRepo}
}

// tahunAjaranPattern memvalidasi format "2025/2026"
var tahunAjaranPattern = regexp.MustCompile(`^\d{4}/\d{4}$`)

func validateSpp(s *domain.Spp) error {
	if s.TahunAjaran == "" {
		return errors.New("tahun_ajaran wajib diisi")
	}
	if !tahunAjaranPattern.MatchString(s.TahunAjaran) {
		return errors.New("tahun_ajaran harus berformat YYYY/YYYY (contoh: 2025/2026)")
	}
	if s.Nominal <= 0 {
		return errors.New("nominal harus lebih besar dari 0")
	}
	return nil
}

func (u *sppUsecase) Create(ctx context.Context, s *domain.Spp) (*domain.Spp, error) {
	if err := validateSpp(s); err != nil {
		return nil, err
	}

	id, err := u.sppRepo.Create(ctx, s)
	if err != nil {
		return nil, err
	}

	return u.sppRepo.FindByID(ctx, id)
}

func (u *sppUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Spp, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.sppRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Spp{}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	pagination := domain.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return list, pagination, nil
}

func (u *sppUsecase) GetByID(ctx context.Context, id int64) (*domain.Spp, error) {
	return u.sppRepo.FindByID(ctx, id)
}

func (u *sppUsecase) Update(ctx context.Context, id int64, s *domain.Spp) (*domain.Spp, error) {
	if err := validateSpp(s); err != nil {
		return nil, err
	}

	s.ID = id
	if err := u.sppRepo.Update(ctx, s); err != nil {
		return nil, err
	}

	return u.sppRepo.FindByID(ctx, id)
}

func (u *sppUsecase) Delete(ctx context.Context, id int64) error {
	return u.sppRepo.Delete(ctx, id)
}
