package mysql

import (
	"context"
	"database/sql"

	"backend/internal/domain"
)

type reportRepository struct {
	db *sql.DB
}

// NewReportRepository mengembalikan implementasi domain.ReportRepository yang berbasis MySQL
func NewReportRepository(db *sql.DB) domain.ReportRepository {
	return &reportRepository{db: db}
}

// GetSummary menghitung total transaksi & total nominal dalam rentang tanggal.
// staffID nil = seluruh sekolah (admin), non-nil = dibatasi ke satu petugas saja.
func (r *reportRepository) GetSummary(ctx context.Context, staffID *int64, tanggalDari, tanggalSampai string) (domain.ReportSummary, error) {
	query := `SELECT COUNT(*), COALESCE(SUM(jumlah_bayar), 0) FROM payments
	          WHERE tanggal_bayar BETWEEN ? AND ?`
	args := []interface{}{tanggalDari, tanggalSampai}

	if staffID != nil {
		query += ` AND staff_id = ?`
		args = append(args, *staffID)
	}

	var summary domain.ReportSummary
	summary.TanggalDari = tanggalDari
	summary.TanggalSampai = tanggalSampai

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&summary.JumlahTransaksi, &summary.TotalNominal)
	if err != nil {
		return domain.ReportSummary{}, err
	}

	return summary, nil
}

// GetBreakdownByStaff merekap jumlah transaksi & total nominal PER PETUGAS dalam rentang tanggal.
// Pakai LEFT JOIN supaya petugas yang tidak ada transaksi sama sekali di rentang itu tetap muncul
// (dengan angka 0), bukan hilang dari laporan — penting untuk akuntabilitas "siapa yang tidak setor".
func (r *reportRepository) GetBreakdownByStaff(ctx context.Context, tanggalDari, tanggalSampai string) ([]domain.StaffReportBreakdown, error) {
	query := `
		SELECT sf.id, sf.nama, COUNT(p.id), COALESCE(SUM(p.jumlah_bayar), 0)
		FROM staffs sf
		LEFT JOIN payments p ON p.staff_id = sf.id AND p.tanggal_bayar BETWEEN ? AND ?
		GROUP BY sf.id, sf.nama
		ORDER BY sf.nama ASC
	`

	rows, err := r.db.QueryContext(ctx, query, tanggalDari, tanggalSampai)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []domain.StaffReportBreakdown{}
	for rows.Next() {
		var b domain.StaffReportBreakdown
		if err := rows.Scan(&b.StaffID, &b.StaffNama, &b.JumlahTransaksi, &b.TotalNominal); err != nil {
			return nil, err
		}
		list = append(list, b)
	}

	return list, rows.Err()
}
