package usecase

import (
	"context"
	"errors"

	"backend/internal/domain"
	"backend/pkg/hash"
)

type staffUsecase struct {
	staffRepo domain.StaffRepository
}

// NewStaffUsecase mengembalikan implementasi domain.StaffUsecase
func NewStaffUsecase(staffRepo domain.StaffRepository) domain.StaffUsecase {
	return &staffUsecase{staffRepo: staffRepo}
}

func validateStaffCreate(s *domain.Staff) error {
	if s.Username == "" {
		return errors.New("username wajib diisi")
	}
	if len(s.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	if s.Nama == "" {
		return errors.New("nama wajib diisi")
	}
	if s.Posisi == "" {
		return errors.New("posisi / jabatan wajib diisi")
	}
	return nil
}

func validateStaffUpdate(s *domain.Staff) error {
	if s.Nama == "" {
		return errors.New("nama wajib diisi")
	}
	if s.Posisi == "" {
		return errors.New("posisi / jabatan wajib diisi")
	}
	return nil
}

// Create memvalidasi input, meng-hash password, lalu mendelegasikan ke repository
// yang akan membuat baris di users (akun login) dan staffs (profil) dalam satu transaksi.
func (u *staffUsecase) Create(ctx context.Context, s *domain.Staff) (*domain.Staff, error) {
	if err := validateStaffCreate(s); err != nil {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(s.Password)
	if err != nil {
		return nil, err
	}

	id, err := u.staffRepo.Create(ctx, s, hashedPassword)
	if err != nil {
		return nil, err
	}

	return u.staffRepo.FindByID(ctx, id)
}

func (u *staffUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Staff, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.staffRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Staff{}
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

func (u *staffUsecase) GetByID(ctx context.Context, id int64) (*domain.Staff, error) {
	return u.staffRepo.FindByID(ctx, id)
}

func (u *staffUsecase) Update(ctx context.Context, id int64, s *domain.Staff) (*domain.Staff, error) {
	if err := validateStaffUpdate(s); err != nil {
		return nil, err
	}

	s.ID = id
	if err := u.staffRepo.Update(ctx, s); err != nil {
		return nil, err
	}

	return u.staffRepo.FindByID(ctx, id)
}

func (u *staffUsecase) Delete(ctx context.Context, id int64) error {
	return u.staffRepo.Delete(ctx, id)
}
