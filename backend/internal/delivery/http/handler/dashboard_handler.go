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
// Struktur response sudah didesain agar tinggal diisi hasil query nyata (COUNT siswa, rata-rata nilai, dst)
// tanpa perlu mengubah kontrak API atau tampilan frontend.

// AdminDashboard menangani GET /api/admin/dashboard (hanya role admin)
func (h *DashboardHandler) AdminDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard admin", gin.H{
		"stats": []StatCard{
			{Label: "Total Guru", Value: "12", Sub: "Aktif mengajar", Color: "blue", Icon: "users"},
			{Label: "Total Siswa", Value: "248", Sub: "6 kelas aktif", Color: "green", Icon: "users"},
			{Label: "Kelas Aktif", Value: "6", Sub: "Tahun ajaran 2025/2026", Color: "amber", Icon: "kelas"},
			{Label: "Akun Terdaftar", Value: "261", Sub: "Semua role", Color: "blue", Icon: "check"},
		},
		"activities": []ActivityItem{
			{Label: "Guru baru ditambahkan", Sub: "Budi Guru · 2 jam lalu"},
			{Label: "Kelas baru dibuat", Sub: "Kelas 9C · kemarin"},
			{Label: "Siswa baru mendaftar", Sub: "Siti Siswa · 2 hari lalu"},
		},
	})
}

// GuruDashboard menangani GET /api/guru/dashboard (hanya role guru)
func (h *DashboardHandler) GuruDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard guru", gin.H{
		"stats": []StatCard{
			{Label: "Kelas Diampu", Value: "3", Sub: "Semester genap", Color: "blue", Icon: "kelas"},
			{Label: "Total Siswa", Value: "84", Sub: "Di semua kelas", Color: "green", Icon: "users"},
			{Label: "Tugas Belum Dinilai", Value: "5", Sub: "Perlu ditindaklanjuti", Color: "amber", Icon: "book"},
			{Label: "Kehadiran Hari Ini", Value: "92%", Sub: "Kelas 9A", Color: "green", Icon: "calendar"},
		},
		"activities": []ActivityItem{
			{Label: "Nilai ulangan diinput", Sub: "Kelas 9A · Matematika · 1 jam lalu"},
			{Label: "Presensi tercatat", Sub: "Kelas 9B · hari ini"},
			{Label: "Tugas baru diberikan", Sub: "Kelas 9C · kemarin"},
		},
	})
}

// SiswaDashboard menangani GET /api/siswa/dashboard (hanya role siswa)
func (h *DashboardHandler) SiswaDashboard(c *gin.Context) {
	response.Success(c, http.StatusOK, "selamat datang di dashboard siswa", gin.H{
		"stats": []StatCard{
			{Label: "Nilai Rata-rata", Value: "87", Sub: "Semester ini", Color: "blue", Icon: "check"},
			{Label: "Kehadiran", Value: "96%", Sub: "Bulan ini", Color: "green", Icon: "calendar"},
			{Label: "Tugas Aktif", Value: "2", Sub: "Menunggu dikumpulkan", Color: "amber", Icon: "book"},
			{Label: "Jadwal Hari Ini", Value: "4", Sub: "Mata pelajaran", Color: "blue", Icon: "kelas"},
		},
		"activities": []ActivityItem{
			{Label: "Nilai ulangan terbit", Sub: "Matematika · 87 · 1 hari lalu"},
			{Label: "Presensi tercatat hadir", Sub: "Hari ini · 07:15"},
			{Label: "Tugas baru diberikan", Sub: "Bahasa Indonesia · kemarin"},
		},
	})
}
