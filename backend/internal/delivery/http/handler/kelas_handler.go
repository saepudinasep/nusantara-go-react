package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type KelasHandler struct {
	kelasUsecase domain.KelasUsecase
}

func NewKelasHandler(kelasUsecase domain.KelasUsecase) *KelasHandler {
	return &KelasHandler{kelasUsecase: kelasUsecase}
}

type kelasRequest struct {
	NamaKelas string `json:"nama_kelas" binding:"required"`
	Tingkat   int    `json:"tingkat" binding:"required"`
}

// List menangani GET /api/admin/kelas
func (h *KelasHandler) List(c *gin.Context) {
	list, err := h.kelasUsecase.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data kelas")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data kelas", list)
}

// Get menangani GET /api/admin/kelas/:id
func (h *KelasHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	kelas, err := h.kelasUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrKelasNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data kelas")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data kelas", kelas)
}

// Create menangani POST /api/admin/kelas
func (h *KelasHandler) Create(c *gin.Context) {
	var req kelasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	kelas := &domain.Kelas{NamaKelas: req.NamaKelas, Tingkat: req.Tingkat}

	created, err := h.kelasUsecase.Create(c.Request.Context(), kelas)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEntry) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "kelas berhasil dibuat", created)
}

// Update menangani PUT /api/admin/kelas/:id
func (h *KelasHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req kelasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	kelas := &domain.Kelas{NamaKelas: req.NamaKelas, Tingkat: req.Tingkat}

	updated, err := h.kelasUsecase.Update(c.Request.Context(), id, kelas)
	if err != nil {
		if errors.Is(err, domain.ErrKelasNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, domain.ErrDuplicateEntry) {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "kelas berhasil diperbarui", updated)
}

// Delete menangani DELETE /api/admin/kelas/:id
func (h *KelasHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.kelasUsecase.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrKelasNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal menghapus kelas")
		return
	}

	response.Success(c, http.StatusOK, "kelas berhasil dihapus", nil)
}

func parseIDParam(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}
