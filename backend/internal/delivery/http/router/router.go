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
func SetupRouter(jwtService *jwt.JWTService, authUsecase domain.AuthUsecase) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // origin frontend Vite
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authHandler := handler.NewAuthHandler(authUsecase)
	dashboardHandler := handler.NewDashboardHandler()

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

			// ---- Role-specific routes ----
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware(string(domain.RoleAdmin)))
			{
				admin.GET("/dashboard", dashboardHandler.AdminDashboard)
			}

			guru := protected.Group("/guru")
			guru.Use(middleware.RoleMiddleware(string(domain.RoleGuru)))
			{
				guru.GET("/dashboard", dashboardHandler.GuruDashboard)
			}

			siswa := protected.Group("/siswa")
			siswa.Use(middleware.RoleMiddleware(string(domain.RoleSiswa)))
			{
				siswa.GET("/dashboard", dashboardHandler.SiswaDashboard)
			}
		}
	}

	return r
}
