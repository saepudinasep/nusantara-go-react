package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type StaffHandler struct {
	staffUsecase domain.StaffUsecase
}

func NewStaffHandler(staffUsecase domain.StaffUsecase) *StaffHandler {
	return &StaffHandler{staffUsecase: staffUsecase}
}

type staffCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nama     string `json:"nama" binding:"required"`
	Posisi   string `json:"posisi" binding:"required"`
}

type staffUpdateRequest struct {
	Nama   string `json:"nama" binding:"required"`
	Posisi string `json:"posisi" binding:"required"`
}

type staffListResponse struct {
	Items      []domain.Staff    `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

// List menangani GET /api/admin/petugas?page=1&limit=10
func (h *StaffHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	list, pagination, err := h.staffUsecase.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data petugas")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data petugas", staffListResponse{
		Items:      list,
		Pagination: pagination,
	})
}

// Get menangani GET /api/admin/petugas/:id
func (h *StaffHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	staff, err := h.staffUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPetugasNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data petugas")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data petugas", staff)
}

// Create menangani POST /api/admin/petugas — sekaligus membuat akun login (role petugas)
func (h *StaffHandler) Create(c *gin.Context) {
	var req staffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	staff := &domain.Staff{
		Username: req.Username,
		Password: req.Password,
		Nama:     req.Nama,
		Posisi:   req.Posisi,
	}

	created, err := h.staffUsecase.Create(c.Request.Context(), staff)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "petugas dan akun login berhasil dibuat", created)
}

// Update menangani PUT /api/admin/petugas/:id — TIDAK mengubah username/password login
func (h *StaffHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req staffUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	staff := &domain.Staff{Nama: req.Nama, Posisi: req.Posisi}

	updated, err := h.staffUsecase.Update(c.Request.Context(), id, staff)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data petugas berhasil diperbarui", updated)
}

// Delete menangani DELETE /api/admin/petugas/:id — turut menghapus akun login petugas (cascade)
func (h *StaffHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.staffUsecase.Delete(c.Request.Context(), id); err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "petugas dan akun login berhasil dihapus", nil)
}

// handleWriteError memetakan error domain ke status HTTP yang sesuai — error SQL mentah
// tidak pernah diteruskan langsung ke client.
func (h *StaffHandler) handleWriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPetugasNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrUsernameTaken):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrPetugasInUse):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrDatabaseBusy):
		response.Error(c, http.StatusServiceUnavailable, err.Error())
	default:
		response.Error(c, http.StatusBadRequest, err.Error())
	}
}
