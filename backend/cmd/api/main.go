package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/delivery/http/router"
	"backend/internal/repository/mysql"
	"backend/internal/usecase"
	"backend/pkg/jwt"
)

func main() {
	// 1. Load konfigurasi dari .env
	cfg := config.LoadConfig()

	// 2. Buka koneksi database
	db := config.NewMySQLConnection(cfg)
	defer db.Close()

	// 3. Inisialisasi service JWT
	jwtService := jwt.NewJWTService(cfg.JWTSecret, cfg.JWTExpireHours)

	// 4. Wiring dependency: repository -> usecase -> handler/router
	//    (dependency injection manual, arah dependency selalu dari luar ke domain)
	userRepository := mysql.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepository, jwtService)

	kelasRepository := mysql.NewKelasRepository(db)
	kelasUsecase := usecase.NewKelasUsecase(kelasRepository)

	profileRepository := mysql.NewProfileRepository(db)
	profileUsecase := usecase.NewProfileUsecase(profileRepository)

	sppRepository := mysql.NewSppRepository(db)
	sppUsecase := usecase.NewSppUsecase(sppRepository)

	studentRepository := mysql.NewStudentRepository(db)
	studentUsecase := usecase.NewStudentUsecase(studentRepository, sppRepository)

	staffRepository := mysql.NewStaffRepository(db)
	staffUsecase := usecase.NewStaffUsecase(staffRepository)

	paymentRepository := mysql.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository)

	reportRepository := mysql.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepository, paymentRepository)

	tagihanUsecase := usecase.NewTagihanUsecase(studentRepository, sppRepository, paymentRepository)

	// dashboardUsecase dibuat PALING TERAKHIR karena bergantung pada tagihanUsecase di atas
	// (kartu "Tagihan Aktif" pada dashboard siswa memakai logika yang sama persis dengan
	// halaman Tagihan & Riwayat, supaya angkanya selalu konsisten).
	dashboardRepository := mysql.NewDashboardRepository(db)
	dashboardUsecase := usecase.NewDashboardUsecase(dashboardRepository, tagihanUsecase)

	// 5. Setup router dan jalankan server
	r := router.SetupRouter(jwtService, authUsecase, kelasUsecase, dashboardUsecase, profileUsecase, sppUsecase, studentUsecase, staffUsecase, paymentUsecase, reportUsecase, tagihanUsecase)

	log.Printf("server berjalan di port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
