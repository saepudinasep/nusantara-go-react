package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type PaymentHandler struct {
	paymentUsecase domain.PaymentUsecase
	staffUsecase   domain.StaffUsecase
	studentUsecase domain.StudentUsecase
}

func NewPaymentHandler(paymentUsecase domain.PaymentUsecase, staffUsecase domain.StaffUsecase, studentUsecase domain.StudentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase, staffUsecase: staffUsecase, studentUsecase: studentUsecase}
}

type paymentAdminCreateRequest struct {
	StaffID      int64   `json:"staff_id" binding:"required"`
	StudentID    int64   `json:"student_id" binding:"required"`
	SppID        int64   `json:"spp_id" binding:"required"`
	BulanDibayar string  `json:"bulan_dibayar" binding:"required"`
	TanggalBayar string  `json:"tanggal_bayar"`
	JumlahBayar  float64 `json:"jumlah_bayar" binding:"required"`
}

type paymentPetugasCreateRequest struct {
	StudentID    int64   `json:"student_id" binding:"required"`
	SppID        int64   `json:"spp_id" binding:"required"`
	BulanDibayar string  `json:"bulan_dibayar" binding:"required"`
	TanggalBayar string  `json:"tanggal_bayar"`
	JumlahBayar  float64 `json:"jumlah_bayar" binding:"required"`
}

type paymentListResponse struct {
	Items      []domain.Payment  `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

// ListAll menangani GET /api/admin/transaksi?page=1&limit=10 — TANPA filter, lihat semua transaksi
func (h *PaymentHandler) ListAll(c *gin.Context) {
	h.list(c, nil)
}

// ListOwn menangani GET /api/petugas/transaksi?page=1&limit=10 — hanya transaksi milik petugas yang login
func (h *PaymentHandler) ListOwn(c *gin.Context) {
	staffID, err := h.staffUsecase.GetOwnStaffID(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrPetugasNotFound) {
			response.Error(c, http.StatusNotFound, "profil petugas kamu belum dibuat oleh admin")
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	h.list(c, &staffID)
}

// ListOwnStudent menangani GET /api/siswa/riwayat?page=1&limit=10 — HANYA riwayat pembayaran milik
// siswa yang sedang login sendiri (student_id diambil dari profil siswa terkait user JWT, tidak
// pernah dari input pemanggil).
func (h *PaymentHandler) ListOwnStudent(c *gin.Context) {
	student, err := h.studentUsecase.GetOwnProfile(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrStudentProfileMissing) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data riwayat pembayaran")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	list, pagination, err := h.paymentUsecase.GetAllByStudent(c.Request.Context(), student.ID, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data riwayat pembayaran")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil riwayat pembayaran", paymentListResponse{
		Items:      list,
		Pagination: pagination,
	})
}

func (h *PaymentHandler) list(c *gin.Context, staffID *int64) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := domain.PaymentFilter{
		Search:        c.Query("search"),
		BulanDibayar:  c.Query("bulan"),
		TanggalDari:   c.Query("tanggal_dari"),
		TanggalSampai: c.Query("tanggal_sampai"),
	}

	list, pagination, err := h.paymentUsecase.GetAll(c.Request.Context(), page, limit, staffID, filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data transaksi", paymentListResponse{
		Items:      list,
		Pagination: pagination,
	})
}

// Get menangani GET /api/admin/transaksi/:id — admin boleh lihat detail transaksi siapa saja
func (h *PaymentHandler) Get(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	payment, err := h.paymentUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data transaksi", payment)
}

// GetOwn menangani GET /api/petugas/transaksi/:id — HANYA boleh lihat detail transaksi milik sendiri.
// Kalau ID-nya valid tapi milik petugas lain, sengaja dibalas 404 (bukan 403) supaya tidak
// mengonfirmasi ke pemanggil bahwa transaksi dengan ID tersebut memang ada.
func (h *PaymentHandler) GetOwn(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	staffID, err := h.staffUsecase.GetOwnStaffID(c.Request.Context(), currentUserID(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	payment, err := h.paymentUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data transaksi")
		return
	}

	if payment.StaffID != staffID {
		response.Error(c, http.StatusNotFound, domain.ErrPaymentNotFound.Error())
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data transaksi", payment)
}

// CreateAsAdmin menangani POST /api/admin/transaksi — admin BEBAS memilih petugas mana yang tercatat
// memproses transaksi ini (dikirim eksplisit lewat staff_id di body).
func (h *PaymentHandler) CreateAsAdmin(c *gin.Context) {
	var req paymentAdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	payment := &domain.Payment{
		StaffID:      req.StaffID,
		StudentID:    req.StudentID,
		SppID:        req.SppID,
		BulanDibayar: req.BulanDibayar,
		TanggalBayar: req.TanggalBayar,
		JumlahBayar:  req.JumlahBayar,
	}

	h.create(c, payment)
}

// CreateAsPetugas menangani POST /api/petugas/transaksi — staff_id TIDAK diambil dari body sama sekali,
// selalu dipaksa jadi ID petugas yang sedang login. Petugas tidak mungkin mencatat transaksi atas nama
// petugas lain walau memodifikasi request secara manual.
func (h *PaymentHandler) CreateAsPetugas(c *gin.Context) {
	var req paymentPetugasCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	staffID, err := h.staffUsecase.GetOwnStaffID(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrPetugasNotFound) {
			response.Error(c, http.StatusNotFound, "profil petugas kamu belum dibuat oleh admin, tidak bisa memproses pembayaran")
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal memproses pembayaran")
		return
	}

	payment := &domain.Payment{
		StaffID:      staffID,
		StudentID:    req.StudentID,
		SppID:        req.SppID,
		BulanDibayar: req.BulanDibayar,
		TanggalBayar: req.TanggalBayar,
		JumlahBayar:  req.JumlahBayar,
	}

	h.create(c, payment)
}

func (h *PaymentHandler) create(c *gin.Context, payment *domain.Payment) {
	created, err := h.paymentUsecase.Create(c.Request.Context(), payment)
	if err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "pembayaran berhasil dicatat", created)
}

// Delete menangani DELETE /api/admin/transaksi/:id — pembatalan transaksi, KHUSUS admin
// (tidak didaftarkan sama sekali di grup route petugas, lihat router.go).
func (h *PaymentHandler) Delete(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "id tidak valid")
		return
	}

	if err := h.paymentUsecase.Delete(c.Request.Context(), id); err != nil {
		h.handleWriteError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "transaksi berhasil dibatalkan", nil)
}

func (h *PaymentHandler) handleWriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPaymentNotFound):
		response.Error(c, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrPaymentDuplicate):
		response.Error(c, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrStudentInvalid),
		errors.Is(err, domain.ErrSppInvalid),
		errors.Is(err, domain.ErrStaffInvalid):
		response.Error(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrDatabaseBusy):
		response.Error(c, http.StatusServiceUnavailable, err.Error())
	default:
		response.Error(c, http.StatusBadRequest, err.Error())
	}
}
