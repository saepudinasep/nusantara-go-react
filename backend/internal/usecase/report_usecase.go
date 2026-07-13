package usecase

import (
	"context"

	"backend/internal/domain"
)

type reportUsecase struct {
	reportRepo  domain.ReportRepository
	paymentRepo domain.PaymentRepository
}

// NewReportUsecase mengembalikan implementasi domain.ReportUsecase.
// Sengaja bergantung langsung pada domain.PaymentRepository (bukan PaymentUsecase) supaya daftar
// transaksi pada laporan bisa diambil tanpa batasan pagination kecil ala tabel UI biasa (laporan
// untuk dicetak butuh semua baris dalam rentang tanggal, bukan cuma 10 per halaman).
func NewReportUsecase(reportRepo domain.ReportRepository, paymentRepo domain.PaymentRepository) domain.ReportUsecase {
	return &reportUsecase{reportRepo: reportRepo, paymentRepo: paymentRepo}
}

// maxReportRows membatasi jumlah baris transaksi yang ikut dicetak dalam satu laporan,
// supaya tetap wajar untuk dicetak/di-render browser walau rentang tanggalnya sangat panjang.
const maxReportRows = 2000

func validateDateRange(tanggalDari, tanggalSampai string) error {
	if tanggalDari == "" || tanggalSampai == "" {
		return domain.ErrInvalidDateRange
	}
	if tanggalDari > tanggalSampai {
		return domain.ErrInvalidDateRange
	}
	return nil
}

func filterPaymentsByDateRange(transactions []domain.Payment, tanggalDari, tanggalSampai string) []domain.Payment {
	if len(transactions) == 0 {
		return nil
	}

	filtered := make([]domain.Payment, 0, len(transactions))
	for _, transaction := range transactions {
		if transaction.TanggalBayar < tanggalDari || transaction.TanggalBayar > tanggalSampai {
			continue
		}
		filtered = append(filtered, transaction)
	}
	return filtered
}

func (u *reportUsecase) GetAdminReport(ctx context.Context, tanggalDari, tanggalSampai string) (domain.ReportSummary, []domain.StaffReportBreakdown, []domain.Payment, error) {
	if err := validateDateRange(tanggalDari, tanggalSampai); err != nil {
		return domain.ReportSummary{}, nil, nil, err
	}

	summary, err := u.reportRepo.GetSummary(ctx, nil, tanggalDari, tanggalSampai)
	if err != nil {
		return domain.ReportSummary{}, nil, nil, err
	}

	breakdown, err := u.reportRepo.GetBreakdownByStaff(ctx, tanggalDari, tanggalSampai)
	if err != nil {
		return domain.ReportSummary{}, nil, nil, err
	}

	transactions, _, err := u.paymentRepo.FindAll(ctx, 1, maxReportRows, nil, domain.PaymentFilter{})
	if err != nil {
		return domain.ReportSummary{}, nil, nil, err
	}

	transactions = filterPaymentsByDateRange(transactions, tanggalDari, tanggalSampai)
	return summary, breakdown, transactions, nil
}

func (u *reportUsecase) GetPetugasReport(ctx context.Context, staffID int64, tanggalDari, tanggalSampai string) (domain.ReportSummary, []domain.Payment, error) {
	if err := validateDateRange(tanggalDari, tanggalSampai); err != nil {
		return domain.ReportSummary{}, nil, err
	}

	summary, err := u.reportRepo.GetSummary(ctx, &staffID, tanggalDari, tanggalSampai)
	if err != nil {
		return domain.ReportSummary{}, nil, err
	}

	transactions, _, err := u.paymentRepo.FindAll(ctx, 1, maxReportRows, &staffID, domain.PaymentFilter{})
	if err != nil {
		return domain.ReportSummary{}, nil, err
	}

	transactions = filterPaymentsByDateRange(transactions, tanggalDari, tanggalSampai)
	return summary, transactions, nil
}
