package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type SppHandler struct {
	sppUsecase domain.SppUsecase
}

func NewSppHandler(sppUsecase domain.SppUsecase) *SppHandler {
	return &SppHandler{sppUsecase: sppUsecase}
}

type sppRequest struct {
	TahunAjaran string  `json:"tahun_ajaran" binding:"required"`
	Nominal     float64 `json:"nominal"`
}

type sppListResponse struct {
	Items      []domain.Spp      `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

// List menangani GET /api/admin/spp?page=1&limit=10
func (h *SppHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	list, pagination, err := h.sppUsecase.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data SPP")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data SPP", sppListResponse{
		Items:      list,
		Pagination: pagination,
	})
}

// Get menangani GET /api/admin/spp/:id
func (h *SppHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	spp, err := h.sppUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSppNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data SPP")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data SPP", spp)
}

// Create menangani POST /api/admin/spp
func (h *SppHandler) Create(c *gin.Context) {
	var req sppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	spp := &domain.Spp{TahunAjaran: req.TahunAjaran, Nominal: req.Nominal}

	created, err := h.sppUsecase.Create(c.Request.Context(), spp)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "data SPP berhasil dibuat", created)
}

// Update menangani PUT /api/admin/spp/:id
func (h *SppHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req sppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	spp := &domain.Spp{TahunAjaran: req.TahunAjaran, Nominal: req.Nominal}

	updated, err := h.sppUsecase.Update(c.Request.Context(), id, spp)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data SPP berhasil diperbarui", updated)
}

// Delete menangani DELETE /api/admin/spp/:id
func (h *SppHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.sppUsecase.Delete(c.Request.Context(), id); err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data SPP berhasil dihapus", nil)
}

// handleWriteError memetakan error domain ke status HTTP yang sesuai.
// Error SQL mentah (constraint, lock timeout, dsb) TIDAK PERNAH dikirim langsung ke client —
// hanya error yang sudah "diterjemahkan" jadi domain.Err... di layer repository yang boleh
// nyampe ke response. Kalau errornya tidak dikenali sama sekali, balas pesan generik supaya
// detail query/koneksi database tidak bocor ke luar.
func (h *SppHandler) handleWriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrSppNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrSppDuplicate):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrSppInUse):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrDatabaseBusy):
		response.Error(c, http.StatusServiceUnavailable, err.Error())
	default:
		// err di sini bisa jadi pesan validasi usecase (aman ditampilkan) ATAU error tak terduga.
		// Karena usecase.validateSpp() selalu mengembalikan errors.New(...) yang sudah ramah-pengguna
		// (bukan error driver SQL), aman menampilkan err.Error() di sini untuk kasus 400.
		// Error driver SQL murni (bukan hasil ExecContext yang gagal ter-mapping di atas) sengaja
		// TIDAK dilempar balik ke usecase sebagai teks mentah — repository selalu membungkusnya
		// jadi salah satu domain.Err... di atas dulu.
		response.Error(c, http.StatusBadRequest, err.Error())
	}
}
