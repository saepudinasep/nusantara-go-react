package usecase

import (
	"context"
	"errors"

	"backend/internal/domain"
	"backend/pkg/hash"
)

type studentUsecase struct {
	studentRepo domain.StudentRepository
}

// NewStudentUsecase mengembalikan implementasi domain.StudentUsecase
func NewStudentUsecase(studentRepo domain.StudentRepository) domain.StudentUsecase {
	return &studentUsecase{studentRepo: studentRepo}
}

func validateStudentCreate(s *domain.Student) error {
	if s.Username == "" {
		return errors.New("username wajib diisi")
	}
	if len(s.Password) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	if s.Nisn == "" {
		return errors.New("nisn wajib diisi")
	}
	if s.Nama == "" {
		return errors.New("nama wajib diisi")
	}
	if s.ClassID <= 0 {
		return errors.New("kelas wajib dipilih")
	}
	return nil
}

func validateStudentUpdate(s *domain.Student) error {
	if s.Nisn == "" {
		return errors.New("nisn wajib diisi")
	}
	if s.Nama == "" {
		return errors.New("nama wajib diisi")
	}
	if s.ClassID <= 0 {
		return errors.New("kelas wajib dipilih")
	}
	return nil
}

// Create memvalidasi input, meng-hash password, lalu mendelegasikan ke repository
// yang akan membuat baris di users (akun login) dan students (profil) dalam satu transaksi.
func (u *studentUsecase) Create(ctx context.Context, s *domain.Student) (*domain.Student, error) {
	if err := validateStudentCreate(s); err != nil {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(s.Password)
	if err != nil {
		return nil, err
	}

	id, err := u.studentRepo.Create(ctx, s, hashedPassword)
	if err != nil {
		return nil, err
	}

	return u.studentRepo.FindByID(ctx, id)
}

func (u *studentUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Student, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.studentRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Student{}
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

func (u *studentUsecase) GetByID(ctx context.Context, id int64) (*domain.Student, error) {
	return u.studentRepo.FindByID(ctx, id)
}

// SearchByNisn dipakai fitur "cari siswa by NISN" pada proses pembayaran (admin & petugas)
func (u *studentUsecase) SearchByNisn(ctx context.Context, nisn string) (*domain.Student, error) {
	return u.studentRepo.FindByNisn(ctx, nisn)
}

func (u *studentUsecase) Update(ctx context.Context, id int64, s *domain.Student) (*domain.Student, error) {
	if err := validateStudentUpdate(s); err != nil {
		return nil, err
	}

	s.ID = id
	if err := u.studentRepo.Update(ctx, s); err != nil {
		return nil, err
	}

	return u.studentRepo.FindByID(ctx, id)
}

func (u *studentUsecase) Delete(ctx context.Context, id int64) error {
	return u.studentRepo.Delete(ctx, id)
}
