package usecase

import (
	"context"
	"errors"

	"backend/internal/domain"
)

type kelasUsecase struct {
	kelasRepo domain.KelasRepository
}

// NewKelasUsecase mengembalikan implementasi domain.KelasUsecase
func NewKelasUsecase(kelasRepo domain.KelasRepository) domain.KelasUsecase {
	return &kelasUsecase{kelasRepo: kelasRepo}
}

func validateKelas(k *domain.Kelas) error {
	if k.NamaKelas == "" {
		return errors.New("nama_kelas wajib diisi")
	}
	if k.Tingkat < 1 {
		return errors.New("tingkat wajib diisi dengan angka valid (contoh: 10, 11, 12)")
	}
	return nil
}

func (u *kelasUsecase) Create(ctx context.Context, k *domain.Kelas) (*domain.Kelas, error) {
	if err := validateKelas(k); err != nil {
		return nil, err
	}

	id, err := u.kelasRepo.Create(ctx, k)
	if err != nil {
		return nil, err
	}

	return u.kelasRepo.FindByID(ctx, id)
}

func (u *kelasUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Kelas, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.kelasRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Kelas{}
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

func (u *kelasUsecase) GetByID(ctx context.Context, id int64) (*domain.Kelas, error) {
	return u.kelasRepo.FindByID(ctx, id)
}

func (u *kelasUsecase) Update(ctx context.Context, id int64, k *domain.Kelas) (*domain.Kelas, error) {
	if err := validateKelas(k); err != nil {
		return nil, err
	}

	k.ID = id
	if err := u.kelasRepo.Update(ctx, k); err != nil {
		return nil, err
	}

	return u.kelasRepo.FindByID(ctx, id)
}

func (u *kelasUsecase) Delete(ctx context.Context, id int64) error {
	return u.kelasRepo.Delete(ctx, id)
}
