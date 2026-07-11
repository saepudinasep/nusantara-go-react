package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/domain"
	"backend/pkg/jwt"
)

// SetupRouter mendaftarkan semua route: public (login) dan protected (butuh JWT + role tertentu)
func SetupRouter(
	jwtService *jwt.JWTService,
	authUsecase domain.AuthUsecase,
	kelasUsecase domain.KelasUsecase,
	dashboardUsecase domain.DashboardUsecase,
	profileUsecase domain.ProfileUsecase,
	sppUsecase domain.SppUsecase,
	studentUsecase domain.StudentUsecase,
	staffUsecase domain.StaffUsecase,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // origin frontend Vite
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authHandler := handler.NewAuthHandler(authUsecase)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase)
	kelasHandler := handler.NewKelasHandler(kelasUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	sppHandler := handler.NewSppHandler(sppUsecase)
	studentHandler := handler.NewStudentHandler(studentUsecase)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	api := r.Group("/api")
	{
		// ---- Public routes ----
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// ---- Protected routes (butuh Bearer token yang valid) ----
		protected := api.Group("")
		protected.Use(middleware.JWTAuthMiddleware(jwtService))
		{
			protected.GET("/auth/me", authHandler.Me)

			// dashboard umum: bisa diakses oleh SEMUA role yang sudah login
			protected.GET("/dashboard", func(c *gin.Context) {
				dashboardHandler.SiswaDashboard(c) // fallback generic, biasanya frontend redirect sesuai role
			})

			// ---- Role: admin ----
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware(string(domain.RoleAdmin)))
			{
				admin.GET("/dashboard", dashboardHandler.AdminDashboard)
				admin.GET("/profile", profileHandler.AdminProfile)

				// CRUD Kelas (khusus admin)
				admin.GET("/kelas", kelasHandler.List)
				admin.POST("/kelas", kelasHandler.Create)
				admin.GET("/kelas/:id", kelasHandler.Get)
				admin.PUT("/kelas/:id", kelasHandler.Update)
				admin.DELETE("/kelas/:id", kelasHandler.Delete)

				// CRUD SPP (khusus admin)
				admin.GET("/spp", sppHandler.List)
				admin.POST("/spp", sppHandler.Create)
				admin.GET("/spp/:id", sppHandler.Get)
				admin.PUT("/spp/:id", sppHandler.Update)
				admin.DELETE("/spp/:id", sppHandler.Delete)

				// CRUD Siswa (khusus admin) — Create sekaligus membuat akun login
				admin.GET("/siswa", studentHandler.List)
				admin.POST("/siswa", studentHandler.Create)
				admin.GET("/siswa/:id", studentHandler.Get)
				admin.PUT("/siswa/:id", studentHandler.Update)
				admin.DELETE("/siswa/:id", studentHandler.Delete)

				// CRUD Petugas (khusus admin) — Create sekaligus membuat akun login
				admin.GET("/petugas", staffHandler.List)
				admin.POST("/petugas", staffHandler.Create)
				admin.GET("/petugas/:id", staffHandler.Get)
				admin.PUT("/petugas/:id", staffHandler.Update)
				admin.DELETE("/petugas/:id", staffHandler.Delete)
			}

			// ---- Role: petugas ----
			petugas := protected.Group("/petugas")
			petugas.Use(middleware.RoleMiddleware(string(domain.RolePetugas)))
			{
				petugas.GET("/dashboard", dashboardHandler.PetugasDashboard)
				petugas.GET("/profile", profileHandler.PetugasProfile)

				// Data Kelas untuk petugas bersifat READ-ONLY (cari data kelas untuk keperluan tagihan siswa).
				// Sengaja hanya List & Get yang didaftarkan di sini — TIDAK ada Create/Update/Delete,
				// jadi petugas tidak akan pernah bisa mengubah data kelas walau tahu endpoint-nya.
				petugas.GET("/kelas", kelasHandler.List)
				petugas.GET("/kelas/:id", kelasHandler.Get)
			}

			// ---- Role: guru ----
			guru := protected.Group("/guru")
			guru.Use(middleware.RoleMiddleware(string(domain.RoleGuru)))
			{
				guru.GET("/dashboard", dashboardHandler.GuruDashboard)
				guru.GET("/profile", profileHandler.GuruProfile)
			}

			// ---- Role: siswa ----
			siswa := protected.Group("/siswa")
			siswa.Use(middleware.RoleMiddleware(string(domain.RoleSiswa)))
			{
				siswa.GET("/dashboard", dashboardHandler.SiswaDashboard)
				siswa.GET("/profile", profileHandler.SiswaProfile)
			}
		}
	}

	return r
}
