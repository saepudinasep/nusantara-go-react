package usecase

import (
	"context"

	"backend/internal/domain"
)

type tagihanUsecase struct {
	studentRepo domain.StudentRepository
	sppRepo     domain.SppRepository
	paymentRepo domain.PaymentRepository
}

// NewTagihanUsecase mengembalikan implementasi domain.TagihanUsecase.
// Sengaja tidak punya repository sendiri — murni menyusun ulang (compose) data dari 3 repository
// yang sudah ada, tanpa query SQL baru yang rumit.
func NewTagihanUsecase(studentRepo domain.StudentRepository, sppRepo domain.SppRepository, paymentRepo domain.PaymentRepository) domain.TagihanUsecase {
	return &tagihanUsecase{studentRepo: studentRepo, sppRepo: sppRepo, paymentRepo: paymentRepo}
}

var bulanList = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// maxSppEntries membatasi jumlah jenis SPP yang diikutkan dalam satu tampilan tagihan —
// dalam praktiknya sekolah biasanya cuma punya beberapa jenis SPP aktif (per tahun ajaran),
// jadi batas ini jauh di atas kebutuhan wajar dan hanya jaring pengaman.
const maxSppEntries = 50

func (u *tagihanUsecase) GetTagihan(ctx context.Context, userID int64) (*domain.Student, []domain.SppTagihan, error) {
	student, err := u.studentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	sppList, _, err := u.sppRepo.FindAll(ctx, 1, maxSppEntries)
	if err != nil {
		return nil, nil, err
	}

	result := make([]domain.SppTagihan, 0, len(sppList))
	for _, spp := range sppList {
		payments, err := u.paymentRepo.FindByStudentAndSpp(ctx, student.ID, spp.ID)
		if err != nil {
			return nil, nil, err
		}

		// petakan bulan yang sudah dibayar -> detail pembayarannya, supaya lookup di bawah O(1)
		paidMonths := make(map[string]domain.Payment, len(payments))
		for _, p := range payments {
			paidMonths[p.BulanDibayar] = p
		}

		bulanan := make([]domain.MonthlyBill, 0, len(bulanList))
		for _, bulan := range bulanList {
			bill := domain.MonthlyBill{Bulan: bulan, Nominal: spp.Nominal}
			if p, ok := paidMonths[bulan]; ok {
				bill.Status = "Lunas"
				bill.TanggalBayar = p.TanggalBayar
			} else {
				bill.Status = "Belum Bayar"
			}
			bulanan = append(bulanan, bill)
		}

		result = append(result, domain.SppTagihan{
			SppID:       spp.ID,
			TahunAjaran: spp.TahunAjaran,
			Nominal:     spp.Nominal,
			Bulanan:     bulanan,
		})
	}

	return student, result, nil
}
