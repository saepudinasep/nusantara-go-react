package usecase

import (
	"context"

	"backend/internal/domain"
	"backend/pkg/hash"
	"backend/pkg/jwt"
)

type authUsecase struct {
	userRepo   domain.UserRepository
	jwtService *jwt.JWTService
}

// NewAuthUsecase mengembalikan implementasi domain.AuthUsecase
func NewAuthUsecase(userRepo domain.UserRepository, jwtService *jwt.JWTService) domain.AuthUsecase {
	return &authUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login memvalidasi credential lalu mengembalikan JWT bearer token + data user
func (u *authUsecase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Disamarkan sebagai invalid credentials agar tidak bocor info user terdaftar atau tidak
		return "", nil, domain.ErrInvalidCredentials
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", nil, domain.ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GetProfile mengambil data user berdasarkan ID (dipakai untuk endpoint /me atau dashboard)
func (u *authUsecase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
