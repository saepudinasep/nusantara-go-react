package usecase

import (
	"context"
	"errors"
	"time"

	"backend/internal/domain"
)

type paymentUsecase struct {
	paymentRepo domain.PaymentRepository
}

// NewPaymentUsecase mengembalikan implementasi domain.PaymentUsecase
func NewPaymentUsecase(paymentRepo domain.PaymentRepository) domain.PaymentUsecase {
	return &paymentUsecase{paymentRepo: paymentRepo}
}

func validatePayment(p *domain.Payment) error {
	if p.StudentID <= 0 {
		return errors.New("siswa wajib dipilih (cari berdasarkan NISN terlebih dahulu)")
	}
	if p.SppID <= 0 {
		return errors.New("jenis SPP wajib dipilih")
	}
	if p.StaffID <= 0 {
		return errors.New("petugas wajib dipilih")
	}
	if p.BulanDibayar == "" {
		return errors.New("bulan yang dibayar wajib diisi")
	}
	if p.TanggalBayar == "" {
		return errors.New("tanggal bayar wajib diisi")
	}
	if p.JumlahBayar <= 0 {
		return errors.New("jumlah bayar harus lebih besar dari 0")
	}
	return nil
}

// Create memvalidasi input, mengecek duplikasi pembayaran pada periode yang sama,
// baru kemudian mendelegasikan insert ke repository.
func (u *paymentUsecase) Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	if err := validatePayment(p); err != nil {
		return nil, err
	}

	sudahBayar, err := u.paymentRepo.HasPaidForPeriod(ctx, p.StudentID, p.SppID, p.BulanDibayar)
	if err != nil {
		return nil, err
	}
	if sudahBayar {
		return nil, domain.ErrPaymentDuplicate
	}

	if p.TanggalBayar == "" {
		p.TanggalBayar = time.Now().Format("2006-01-02")
	}

	id, err := u.paymentRepo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return u.paymentRepo.FindByID(ctx, id)
}

func (u *paymentUsecase) GetAll(ctx context.Context, page, limit int, staffID *int64, filter domain.PaymentFilter) ([]domain.Payment, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.paymentRepo.FindAll(ctx, page, limit, staffID, filter)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Payment{}
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

func (u *paymentUsecase) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	return u.paymentRepo.FindByID(ctx, id)
}

func (u *paymentUsecase) Delete(ctx context.Context, id int64) error {
	return u.paymentRepo.Delete(ctx, id)
}

// GetAllByStudent dipakai halaman "Riwayat Pembayaran" siswa
func (u *paymentUsecase) GetAllByStudent(ctx context.Context, studentID int64, page, limit int) ([]domain.Payment, domain.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	list, total, err := u.paymentRepo.FindAllByStudent(ctx, studentID, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	if list == nil {
		list = []domain.Payment{}
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
