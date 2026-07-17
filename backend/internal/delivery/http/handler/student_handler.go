package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type StudentHandler struct {
	studentUsecase domain.StudentUsecase
}

func NewStudentHandler(studentUsecase domain.StudentUsecase) *StudentHandler {
	return &StudentHandler{studentUsecase: studentUsecase}
}

type studentCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nisn     string `json:"nisn" binding:"required"`
	Nama     string `json:"nama" binding:"required"`
	ClassID  int64  `json:"class_id" binding:"required"`
	Alamat   string `json:"alamat"`
	NoTelp   string `json:"no_telp"`
}

type studentUpdateRequest struct {
	Nisn    string `json:"nisn" binding:"required"`
	Nama    string `json:"nama" binding:"required"`
	ClassID int64  `json:"class_id" binding:"required"`
	Alamat  string `json:"alamat"`
	NoTelp  string `json:"no_telp"`
}

type studentListResponse struct {
	Items      []domain.Student  `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

// List menangani GET /api/admin/siswa?page=1&limit=10
// List menangani GET /api/{role}/siswa?page=1&limit=10
// Query opsional: with_status=true untuk menyertakan status Lunas/Belum Bayar SPP aktif bulan
// berjalan per siswa; only_unpaid=true (butuh with_status=true) untuk cuma menampilkan yang nunggak.
func (h *StudentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	withStatus := c.Query("with_status") == "true"
	onlyUnpaid := c.Query("only_unpaid") == "true"

	var list []domain.Student
	var pagination domain.Pagination
	var err error

	if withStatus {
		list, pagination, err = h.studentUsecase.GetAllWithTunggakanStatus(c.Request.Context(), page, limit, onlyUnpaid)
	} else {
		list, pagination, err = h.studentUsecase.GetAll(c.Request.Context(), page, limit)
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data siswa")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data siswa", studentListResponse{
		Items:      list,
		Pagination: pagination,
	})
}

// Get menangani GET /api/admin/siswa/:id
func (h *StudentHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	student, err := h.studentUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrSiswaNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data siswa")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data siswa", student)
}

// Create menangani POST /api/admin/siswa — sekaligus membuat akun login (role siswa) untuk siswa ini
func (h *StudentHandler) Create(c *gin.Context) {
	var req studentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	student := &domain.Student{
		Username: req.Username,
		Password: req.Password,
		Nisn:     req.Nisn,
		Nama:     req.Nama,
		ClassID:  req.ClassID,
		Alamat:   req.Alamat,
		NoTelp:   req.NoTelp,
	}

	created, err := h.studentUsecase.Create(c.Request.Context(), student)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "siswa dan akun login berhasil dibuat", created)
}

// Update menangani PUT /api/admin/siswa/:id — TIDAK mengubah username/password login
func (h *StudentHandler) Update(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	var req studentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	student := &domain.Student{
		Nisn:    req.Nisn,
		Nama:    req.Nama,
		ClassID: req.ClassID,
		Alamat:  req.Alamat,
		NoTelp:  req.NoTelp,
	}

	updated, err := h.studentUsecase.Update(c.Request.Context(), id, student)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data siswa berhasil diperbarui", updated)
}

// Delete menangani DELETE /api/admin/siswa/:id — turut menghapus akun login siswa (cascade)
func (h *StudentHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.studentUsecase.Delete(c.Request.Context(), id); err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "siswa dan akun login berhasil dihapus", nil)
}

// handleWriteError memetakan error domain ke status HTTP yang sesuai — error SQL mentah
// tidak pernah diteruskan langsung ke client (lihat catatan yang sama di spp_handler.go).
// SearchByNisn menangani GET /api/{role}/siswa/cari?nisn=XXX — dipakai alur "cari siswa" saat proses pembayaran
func (h *StudentHandler) SearchByNisn(c *gin.Context) {
	nisn := c.Query("nisn")
	if nisn == "" {
		response.Error(c, http.StatusBadRequest, "parameter nisn wajib diisi")
		return
	}

	student, err := h.studentUsecase.SearchByNisn(c.Request.Context(), nisn)
	if err != nil {
		if errors.Is(err, domain.ErrStudentInvalid) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mencari data siswa")
		return
	}

	response.Success(c, http.StatusOK, "siswa ditemukan", student)
}

func (h *StudentHandler) handleWriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrSiswaNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrUsernameTaken):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrNisnTaken):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrSiswaInUse):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrKelasTidakValid):
		response.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrDatabaseBusy):
		response.Error(c, http.StatusServiceUnavailable, err.Error())
	default:
		response.Error(c, http.StatusBadRequest, err.Error())
	}
}
