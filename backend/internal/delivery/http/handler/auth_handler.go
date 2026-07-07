package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/delivery/http/middleware"
	"backend/internal/domain"
	"backend/pkg/response"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

// Login menangani POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
		return
	}

	token, user, err := h.authUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	response.Success(c, http.StatusOK, "login berhasil", loginResponse{
		Token: token,
		User:  *user,
	})
}

// Me menangani GET /api/auth/me (memerlukan JWT), mengembalikan profil user yang sedang login
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)

	user, err := h.authUsecase.GetProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil profil", user)
}
