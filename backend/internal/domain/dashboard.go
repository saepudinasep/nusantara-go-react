package domain

import "context"

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

// DashboardRepository adalah interface (port) untuk query agregat dashboard ke MySQL.
// Sengaja dibuat "dumb" (murni ambil angka/baris mentah) — logika format & penyusunan StatCard
// ada di usecase, bukan di sini.
type DashboardRepository interface {
	CountStaffs(ctx context.Context) (int64, error)
	CountGuru(ctx context.Context) (int64, error)
	CountStudents(ctx context.Context) (int64, error)
	CountClasses(ctx context.Context) (int64, error)
	CountPaymentsToday(ctx context.Context) (int64, error)
	CountPaidStudentsForMonth(ctx context.Context, monthName string, year int) (int64, error)
	SumPaymentsInMonth(ctx context.Context, month int, year int) (float64, error)

	RecentActivitiesAdmin(ctx context.Context, limit int) ([]ActivityItem, error)
	RecentPayments(ctx context.Context, limit int) ([]ActivityItem, error)
	RecentPaymentsByStudent(ctx context.Context, studentID int64, limit int) ([]ActivityItem, error)

	// FindStudentByUserID mengambil profil siswa (id, nama kelas, tingkat) berdasarkan user_id yang login.
	// Mengembalikan studentID = 0 jika user tersebut belum punya profil siswa.
	FindStudentByUserID(ctx context.Context, userID int64) (studentID int64, namaKelas string, tingkat int, err error)
	HasStudentPaidForMonth(ctx context.Context, studentID int64, monthName string, year int) (bool, error)
	SumPaymentsByStudent(ctx context.Context, studentID int64) (float64, error)
}

// DashboardUsecase adalah interface (port) untuk business logic penyusunan data dashboard per role
type DashboardUsecase interface {
	GetAdminDashboard(ctx context.Context) ([]StatCard, []ActivityItem, error)
	GetPetugasDashboard(ctx context.Context) ([]StatCard, []ActivityItem, error)
	GetGuruDashboard(ctx context.Context) ([]StatCard, []ActivityItem, error)
	GetSiswaDashboard(ctx context.Context, userID int64) ([]StatCard, []ActivityItem, error)
}
