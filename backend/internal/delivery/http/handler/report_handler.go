package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type ReportHandler struct {
	reportUsecase domain.ReportUsecase
	staffUsecase  domain.StaffUsecase
}

func NewReportHandler(reportUsecase domain.ReportUsecase, staffUsecase domain.StaffUsecase) *ReportHandler {
	return &ReportHandler{reportUsecase: reportUsecase, staffUsecase: staffUsecase}
}

type adminReportResponse struct {
	Summary      domain.ReportSummary          `json:"summary"`
	Breakdown    []domain.StaffReportBreakdown `json:"breakdown"`
	Transactions []domain.Payment              `json:"transactions"`
}

type petugasReportResponse struct {
	Summary      domain.ReportSummary `json:"summary"`
	Transactions []domain.Payment     `json:"transactions"`
}

// AdminReport menangani GET /api/admin/laporan?tanggal_dari=YYYY-MM-DD&tanggal_sampai=YYYY-MM-DD
// Laporan GLOBAL: ringkasan seluruh sekolah + rekap per petugas.
func (h *ReportHandler) AdminReport(c *gin.Context) {
	tanggalDari := c.Query("tanggal_dari")
	tanggalSampai := c.Query("tanggal_sampai")

	summary, breakdown, transactions, err := h.reportUsecase.GetAdminReport(c.Request.Context(), tanggalDari, tanggalSampai)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDateRange) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal menyusun laporan")
		return
	}

	response.Success(c, http.StatusOK, "berhasil menyusun laporan", adminReportResponse{
		Summary:      summary,
		Breakdown:    breakdown,
		Transactions: transactions,
	})
}

// PetugasReport menangani GET /api/petugas/laporan?tanggal_dari=YYYY-MM-DD&tanggal_sampai=YYYY-MM-DD
// Laporan HARIAN: dibatasi ke rekap milik petugas yang sedang login saja.
func (h *ReportHandler) PetugasReport(c *gin.Context) {
	tanggalDari := c.Query("tanggal_dari")
	tanggalSampai := c.Query("tanggal_sampai")

	staffID, err := h.staffUsecase.GetOwnStaffID(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrPetugasNotFound) {
			response.Error(c, http.StatusNotFound, "profil petugas kamu belum dibuat oleh admin")
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal menyusun laporan")
		return
	}

	summary, transactions, err := h.reportUsecase.GetPetugasReport(c.Request.Context(), staffID, tanggalDari, tanggalSampai)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDateRange) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal menyusun laporan")
		return
	}

	response.Success(c, http.StatusOK, "berhasil menyusun laporan", petugasReportResponse{
		Summary:      summary,
		Transactions: transactions,
	})
}
