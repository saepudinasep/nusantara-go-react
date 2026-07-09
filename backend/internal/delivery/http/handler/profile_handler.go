package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/delivery/http/middleware"
	"backend/internal/domain"
	"backend/pkg/response"
)

type ProfileHandler struct {
	profileUsecase domain.ProfileUsecase
}

func NewProfileHandler(profileUsecase domain.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{profileUsecase: profileUsecase}
}

func currentUserID(c *gin.Context) int64 {
	userID, _ := c.Get(middleware.ContextUserID)
	return userID.(int64)
}

// AdminProfile menangani GET /api/admin/profile (murni dari tabel users)
func (h *ProfileHandler) AdminProfile(c *gin.Context) {
	profile, err := h.profileUsecase.GetAdminProfile(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data profil")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil profil", profile)
}

// GuruProfile menangani GET /api/guru/profile (murni dari tabel users, belum ada tabel profil guru di V1)
func (h *ProfileHandler) GuruProfile(c *gin.Context) {
	profile, err := h.profileUsecase.GetGuruProfile(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data profil")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil profil", profile)
}

// PetugasProfile menangani GET /api/petugas/profile (gabungan tabel users + staffs)
func (h *ProfileHandler) PetugasProfile(c *gin.Context) {
	profile, err := h.profileUsecase.GetPetugasProfile(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data profil")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil profil", profile)
}

// SiswaProfile menangani GET /api/siswa/profile (gabungan tabel users + students + classes)
func (h *ProfileHandler) SiswaProfile(c *gin.Context) {
	profile, err := h.profileUsecase.GetSiswaProfile(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data profil")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil profil", profile)
}
