package domain

import "context"

// MonthlyBill merepresentasikan status tagihan SPP untuk SATU bulan pada satu jenis SPP
type MonthlyBill struct {
	Bulan        string  `json:"bulan"`
	Status       string  `json:"status"` // "Lunas" | "Belum Bayar"
	Nominal      float64 `json:"nominal"`
	TanggalBayar string  `json:"tanggal_bayar,omitempty"` // hanya terisi kalau Status = Lunas
}

// SppTagihan mengelompokkan status 12 bulan (MonthlyBill) di bawah satu jenis SPP (tahun ajaran)
type SppTagihan struct {
	SppID       int64         `json:"spp_id"`
	TahunAjaran string        `json:"tahun_ajaran"`
	Nominal     float64       `json:"nominal"`
	Bulanan     []MonthlyBill `json:"bulanan"`
}

// TagihanUsecase adalah interface (port) untuk business logic penyusunan status tagihan per bulan.
// Sengaja tidak punya Repository sendiri — usecase ini murni menyusun ulang (compose) data yang
// sudah diambil lewat StudentRepository, SppRepository, dan PaymentRepository yang sudah ada.
type TagihanUsecase interface {
	GetTagihan(ctx context.Context, userID int64) (*Student, []SppTagihan, error)
}
