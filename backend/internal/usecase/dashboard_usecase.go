package usecase

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain"
)

type dashboardUsecase struct {
	repo           domain.DashboardRepository
	tagihanUsecase domain.TagihanUsecase
}

// NewDashboardUsecase mengembalikan implementasi domain.DashboardUsecase.
// tagihanUsecase dipakai khusus GetSiswaDashboard supaya kartu "Tagihan Aktif" konsisten dengan
// halaman Tagihan & Riwayat (hitung per-bulan per-SPP yang sesungguhnya, bukan cek 1 bulan saja).
func NewDashboardUsecase(repo domain.DashboardRepository, tagihanUsecase domain.TagihanUsecase) domain.DashboardUsecase {
	return &dashboardUsecase{repo: repo, tagihanUsecase: tagihanUsecase}
}

var indonesianMonths = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func currentMonthName() string {
	return indonesianMonths[int(time.Now().Month())]
}

func formatRupiah(amount float64) string {
	return "Rp" + formatThousands(int64(amount))
}

// formatThousands menyisipkan titik sebagai pemisah ribuan ala format Indonesia (150000 -> "150.000")
func formatThousands(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result []byte
	offset := len(s) % 3
	if offset == 0 {
		offset = 3
	}
	result = append(result, s[:offset]...)
	for i := offset; i < len(s); i += 3 {
		result = append(result, '.')
		result = append(result, s[i:i+3]...)
	}
	return string(result)
}

func (u *dashboardUsecase) GetAdminDashboard(ctx context.Context) ([]domain.StatCard, []domain.ActivityItem, error) {
	totalStaffs, err := u.repo.CountStaffs(ctx)
	if err != nil {
		return nil, nil, err
	}
	totalGuru, err := u.repo.CountGuru(ctx)
	if err != nil {
		return nil, nil, err
	}
	totalStudents, err := u.repo.CountStudents(ctx)
	if err != nil {
		return nil, nil, err
	}
	totalClasses, err := u.repo.CountClasses(ctx)
	if err != nil {
		return nil, nil, err
	}

	stats := []domain.StatCard{
		{Label: "Total Petugas", Value: fmt.Sprintf("%d", totalStaffs), Sub: "Aktif bertugas", Color: "blue", Icon: "users"},
		{Label: "Total Guru", Value: fmt.Sprintf("%d", totalGuru), Sub: "Aktif mengajar", Color: "blue", Icon: "users"},
		{Label: "Total Siswa", Value: fmt.Sprintf("%d", totalStudents), Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
		{Label: "Total Kelas", Value: fmt.Sprintf("%d", totalClasses), Sub: "Tahun ajaran berjalan", Color: "amber", Icon: "kelas"},
	}

	activities, err := u.repo.RecentActivitiesAdmin(ctx, 5)
	if err != nil {
		return nil, nil, err
	}

	return stats, activities, nil
}

func (u *dashboardUsecase) GetPetugasDashboard(ctx context.Context) ([]domain.StatCard, []domain.ActivityItem, error) {
	now := time.Now()

	paymentsToday, err := u.repo.CountPaymentsToday(ctx)
	if err != nil {
		return nil, nil, err
	}
	totalStudents, err := u.repo.CountStudents(ctx)
	if err != nil {
		return nil, nil, err
	}
	paidThisMonth, err := u.repo.CountPaidStudentsForMonth(ctx, currentMonthName(), now.Year())
	if err != nil {
		return nil, nil, err
	}
	tunggakan := totalStudents - paidThisMonth
	if tunggakan < 0 {
		tunggakan = 0
	}
	totalReceived, err := u.repo.SumPaymentsInMonth(ctx, int(now.Month()), now.Year())
	if err != nil {
		return nil, nil, err
	}

	stats := []domain.StatCard{
		{Label: "Transaksi Hari Ini", Value: fmt.Sprintf("%d", paymentsToday), Sub: "Pembayaran SPP tercatat", Color: "blue", Icon: "check"},
		{Label: "Total Siswa", Value: fmt.Sprintf("%d", totalStudents), Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
		{Label: "Tunggakan", Value: fmt.Sprintf("%d", tunggakan), Sub: "Siswa belum bayar bulan ini", Color: "amber", Icon: "book"},
		{Label: "Total Diterima", Value: formatRupiah(totalReceived), Sub: "Bulan ini", Color: "green", Icon: "calendar"},
	}

	activities, err := u.repo.RecentPayments(ctx, 5)
	if err != nil {
		return nil, nil, err
	}

	return stats, activities, nil
}

func (u *dashboardUsecase) GetGuruDashboard(ctx context.Context) ([]domain.StatCard, []domain.ActivityItem, error) {
	// Kelas Diampu, Materi, dan Kuis masih placeholder karena tabel pendukungnya
	// (jadwal mengajar, materi, kuis) belum ada di skema V1 — direncanakan untuk V2 (Learning Center).
	totalStudents, err := u.repo.CountStudents(ctx)
	if err != nil {
		return nil, nil, err
	}

	stats := []domain.StatCard{
		{Label: "Kelas Diampu", Value: "0", Sub: "Belum ada jadwal (fitur V2)", Color: "blue", Icon: "kelas"},
		{Label: "Total Siswa", Value: fmt.Sprintf("%d", totalStudents), Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
		{Label: "Materi Diunggah", Value: "0", Sub: "Fitur V2 - Learning Center", Color: "amber", Icon: "book"},
		{Label: "Kuis Aktif", Value: "0", Sub: "Fitur V2 - Learning Center", Color: "blue", Icon: "check"},
	}

	activities := []domain.ActivityItem{
		{Label: "Fitur pengajaran (jadwal, materi, kuis) akan hadir di V2", Sub: "Learning Center"},
	}

	return stats, activities, nil
}

func (u *dashboardUsecase) GetSiswaDashboard(ctx context.Context, userID int64) ([]domain.StatCard, []domain.ActivityItem, error) {
	studentID, namaKelas, tingkat, err := u.repo.FindStudentByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// User login sebagai siswa tapi belum punya profil di tabel students — tampilkan state kosong, bukan error.
	if studentID == 0 {
		stats := []domain.StatCard{
			{Label: "Status SPP", Value: "-", Sub: "Profil siswa belum terhubung", Color: "amber", Icon: "check"},
			{Label: "Kelas", Value: "-", Sub: "Hubungi admin", Color: "blue", Icon: "kelas"},
			{Label: "Tagihan Aktif", Value: "-", Sub: "-", Color: "blue", Icon: "book"},
			{Label: "Total Dibayar", Value: "-", Sub: "-", Color: "green", Icon: "calendar"},
		}
		return stats, []domain.ActivityItem{}, nil
	}

	totalDibayar, err := u.repo.SumPaymentsByStudent(ctx, studentID)
	if err != nil {
		return nil, nil, err
	}

	// Pakai TagihanUsecase yang sama dengan halaman "Tagihan & Riwayat" supaya angkanya SELALU
	// konsisten di semua tempat — bukan cuma cek 1 bulan seperti versi sebelumnya (yang bisa salah
	// kalau siswa punya tunggakan di bulan-bulan lampau atau di jenis SPP lain).
	_, tagihanList, err := u.tagihanUsecase.GetTagihan(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	tunggakanBulan := 0
	tahunAjaranAktif := "-"
	if len(tagihanList) > 0 {
		// SppRepository.FindAll mengurutkan tahun_ajaran DESC, jadi entri pertama = SPP paling baru/aktif.
		sppAktif := tagihanList[0]
		tahunAjaranAktif = sppAktif.TahunAjaran
		currentMonthIndex := int(time.Now().Month()) // 1-12

		for i, bill := range sppAktif.Bulanan {
			bulanIndex := i + 1
			if bulanIndex > currentMonthIndex {
				break // bulan yang belum jatuh tempo tidak dihitung sebagai tunggakan
			}
			if bill.Status == "Belum Bayar" {
				tunggakanBulan++
			}
		}
	}

	statusSPP := "Lunas"
	statusColor := "green"
	statusSub := fmt.Sprintf("Tahun ajaran %s", tahunAjaranAktif)
	if tunggakanBulan > 0 {
		statusSPP = "Ada Tunggakan"
		statusColor = "amber"
	}

	tagihanSub := "Tidak ada tunggakan"
	if tunggakanBulan > 0 {
		tagihanSub = fmt.Sprintf("%d bulan belum dibayar", tunggakanBulan)
	}

	tagihanColor := "blue"
	if tunggakanBulan > 0 {
		tagihanColor = "amber"
	}

	stats := []domain.StatCard{
		{Label: "Status SPP", Value: statusSPP, Sub: statusSub, Color: statusColor, Icon: "check"},
		{Label: "Kelas", Value: namaKelas, Sub: fmt.Sprintf("Tingkat %d", tingkat), Color: "blue", Icon: "kelas"},
		{Label: "Tagihan Aktif", Value: fmt.Sprintf("%d", tunggakanBulan), Sub: tagihanSub, Color: tagihanColor, Icon: "book"},
		{Label: "Total Dibayar", Value: formatRupiah(totalDibayar), Sub: "Akumulasi seluruh transaksi", Color: "green", Icon: "calendar"},
	}

	activities, err := u.repo.RecentPaymentsByStudent(ctx, studentID, 5)
	if err != nil {
		return nil, nil, err
	}

	return stats, activities, nil
}
