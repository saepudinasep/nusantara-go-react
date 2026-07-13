package domain

import "context"

// ReportSummary adalah ringkasan angka untuk satu rentang tanggal (total transaksi + total nominal)
type ReportSummary struct {
	TanggalDari     string  `json:"tanggal_dari"`
	TanggalSampai   string  `json:"tanggal_sampai"`
	JumlahTransaksi int64   `json:"jumlah_transaksi"`
	TotalNominal    float64 `json:"total_nominal"`
}

// StaffReportBreakdown merepresentasikan rekap per-petugas (dipakai khusus Laporan Global admin)
type StaffReportBreakdown struct {
	StaffID         int64   `json:"staff_id"`
	StaffNama       string  `json:"staff_nama"`
	JumlahTransaksi int64   `json:"jumlah_transaksi"`
	TotalNominal    float64 `json:"total_nominal"`
}

// ReportRepository adalah interface (port) untuk query agregat laporan keuangan
type ReportRepository interface {
	GetSummary(ctx context.Context, staffID *int64, tanggalDari, tanggalSampai string) (ReportSummary, error)
	GetBreakdownByStaff(ctx context.Context, tanggalDari, tanggalSampai string) ([]StaffReportBreakdown, error)
}

// ReportUsecase adalah interface (port) untuk business logic penyusunan laporan
type ReportUsecase interface {
	// GetAdminReport: laporan GLOBAL — ringkasan seluruh sekolah + rekap per petugas + daftar transaksi
	GetAdminReport(ctx context.Context, tanggalDari, tanggalSampai string) (ReportSummary, []StaffReportBreakdown, []Payment, error)
	// GetPetugasReport: laporan HARIAN milik sendiri — ringkasan + daftar transaksi, dibatasi staffID sendiri
	GetPetugasReport(ctx context.Context, staffID int64, tanggalDari, tanggalSampai string) (ReportSummary, []Payment, error)
}
