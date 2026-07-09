package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/delivery/http/middleware"
	"backend/internal/domain"
	"backend/pkg/response"
)

type DashboardHandler struct {
	dashboardUsecase domain.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase domain.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{dashboardUsecase: dashboardUsecase}
}

type dashboardResponse struct {
	Stats      []domain.StatCard     `json:"stats"`
	Activities []domain.ActivityItem `json:"activities"`
}

// AdminDashboard menangani GET /api/admin/dashboard (hanya role admin)
func (h *DashboardHandler) AdminDashboard(c *gin.Context) {
	stats, activities, err := h.dashboardUsecase.GetAdminDashboard(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data dashboard")
		return
	}

	response.Success(c, http.StatusOK, "selamat datang di dashboard admin", dashboardResponse{
		Stats:      stats,
		Activities: activities,
	})
}

// PetugasDashboard menangani GET /api/petugas/dashboard (hanya role petugas)
func (h *DashboardHandler) PetugasDashboard(c *gin.Context) {
	stats, activities, err := h.dashboardUsecase.GetPetugasDashboard(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data dashboard")
		return
	}

	response.Success(c, http.StatusOK, "selamat datang di dashboard petugas", dashboardResponse{
		Stats:      stats,
		Activities: activities,
	})
}

// GuruDashboard menangani GET /api/guru/dashboard (hanya role guru)
func (h *DashboardHandler) GuruDashboard(c *gin.Context) {
	stats, activities, err := h.dashboardUsecase.GetGuruDashboard(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data dashboard")
		return
	}

	response.Success(c, http.StatusOK, "selamat datang di dashboard guru", dashboardResponse{
		Stats:      stats,
		Activities: activities,
	})
}

// SiswaDashboard menangani GET /api/siswa/dashboard (hanya role siswa)
func (h *DashboardHandler) SiswaDashboard(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)

	stats, activities, err := h.dashboardUsecase.GetSiswaDashboard(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data dashboard")
		return
	}

	response.Success(c, http.StatusOK, "selamat datang di dashboard siswa", dashboardResponse{
		Stats:      stats,
		Activities: activities,
	})
}
