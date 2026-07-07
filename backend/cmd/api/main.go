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

	// 5. Setup router dan jalankan server
	r := router.SetupRouter(jwtService, authUsecase)

	log.Printf("server berjalan di port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
