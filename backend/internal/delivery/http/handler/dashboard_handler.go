package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/pkg/response"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// StatCard merepresentasikan satu kartu ringkasan angka di dashboard (dipetakan ke komponen stat-card di frontend)
type StatCard struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Sub   string `json:"sub"`
	Color string `json:"color"` // blue | green | amber | red — dipetakan ke warna aksen kartu
	Icon  string `json:"icon"`  // dashboard | users | book | kelas | laporan | calendar | check
}

// ActivityItem merepresentasikan satu baris di panel "Aktivitas Terkini"
type ActivityItem struct {
	Label string `json:"label"`
	Sub   string `json:"sub"`
}

// NOTE: Angka pada stats & activity di bawah ini masih data contoh (belum terhubung ke query MySQL).
// Struktur response sudah didesain agar tinggal diisi hasil query nyata (COUNT siswa, total SPP terbayar, dst)
// tanpa perlu mengubah kontrak API atau tampilan frontend.

// AdminDashboard menangani GET /api/admin/dashboard (hanya role admin)
func (h *DashboardHandler) AdminDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard admin", gin.H{
		"stats": []StatCard{
			{Label: "Total Petugas", Value: "1", Sub: "Aktif bertugas", Color: "blue", Icon: "users"},
			{Label: "Total Guru", Value: "1", Sub: "Aktif mengajar", Color: "blue", Icon: "users"},
			{Label: "Total Siswa", Value: "1", Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
			{Label: "Total Kelas", Value: "3", Sub: "Tahun ajaran 2025/2026", Color: "amber", Icon: "kelas"},
		},
		"activities": []ActivityItem{
			{Label: "Kelas baru ditambahkan", Sub: "XI RPL · kemarin"},
			{Label: "Siswa baru mendaftar", Sub: "Siswa Satu · 2 hari lalu"},
			{Label: "Pembayaran SPP diterima", Sub: "Rp150.000 · 1 hari lalu"},
		},
	})
}

// PetugasDashboard menangani GET /api/petugas/dashboard (hanya role petugas)
func (h *DashboardHandler) PetugasDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard petugas", gin.H{
		"stats": []StatCard{
			{Label: "Transaksi Hari Ini", Value: "1", Sub: "Pembayaran SPP tercatat", Color: "blue", Icon: "check"},
			{Label: "Total Siswa", Value: "1", Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
			{Label: "Tunggakan", Value: "0", Sub: "Siswa belum bayar bulan ini", Color: "amber", Icon: "book"},
			{Label: "Total Diterima", Value: "Rp150.000", Sub: "Bulan ini", Color: "green", Icon: "calendar"},
		},
		"activities": []ActivityItem{
			{Label: "Pembayaran SPP diterima", Sub: "Siswa Satu · Rp150.000 · 1 hari lalu"},
		},
	})
}

// GuruDashboard menangani GET /api/guru/dashboard (hanya role guru)
func (h *DashboardHandler) GuruDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard guru", gin.H{
		"stats": []StatCard{
			{Label: "Kelas Diampu", Value: "0", Sub: "Belum ada jadwal (fitur V2)", Color: "blue", Icon: "kelas"},
			{Label: "Total Siswa", Value: "1", Sub: "Terdaftar aktif", Color: "green", Icon: "users"},
			{Label: "Materi Diunggah", Value: "0", Sub: "Fitur V2 - Learning Center", Color: "amber", Icon: "book"},
			{Label: "Kuis Aktif", Value: "0", Sub: "Fitur V2 - Learning Center", Color: "blue", Icon: "check"},
		},
		"activities": []ActivityItem{
			{Label: "Fitur pengajaran (jadwal, materi, kuis) akan hadir di V2", Sub: "Learning Center"},
		},
	})
}

// SiswaDashboard menangani GET /api/siswa/dashboard (hanya role siswa)
func (h *DashboardHandler) SiswaDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard siswa", gin.H{
		"stats": []StatCard{
			{Label: "Status SPP", Value: "Lunas", Sub: "Bulan Juli 2026", Color: "green", Icon: "check"},
			{Label: "Kelas", Value: "XA", Sub: "Tingkat 10", Color: "blue", Icon: "kelas"},
			{Label: "Tagihan Aktif", Value: "0", Sub: "Tidak ada tunggakan", Color: "blue", Icon: "book"},
			{Label: "Total Dibayar", Value: "Rp150.000", Sub: "Tahun ajaran 2025/2026", Color: "green", Icon: "calendar"},
		},
		"activities": []ActivityItem{
			{Label: "Pembayaran SPP berhasil", Sub: "Rp150.000 · 1 hari lalu"},
		},
	})
}
