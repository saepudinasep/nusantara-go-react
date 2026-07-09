package usecase

import (
	"context"

	"backend/internal/domain"
)

type profileUsecase struct {
	repo domain.ProfileRepository
}

// NewProfileUsecase mengembalikan implementasi domain.ProfileUsecase
func NewProfileUsecase(repo domain.ProfileRepository) domain.ProfileUsecase {
	return &profileUsecase{repo: repo}
}

func (u *profileUsecase) GetAdminProfile(ctx context.Context, userID int64) (*domain.AdminProfile, error) {
	return u.repo.FindAdminProfile(ctx, userID)
}

func (u *profileUsecase) GetGuruProfile(ctx context.Context, userID int64) (*domain.GuruProfile, error) {
	return u.repo.FindGuruProfile(ctx, userID)
}

func (u *profileUsecase) GetPetugasProfile(ctx context.Context, userID int64) (*domain.PetugasProfile, error) {
	return u.repo.FindPetugasProfile(ctx, userID)
}

func (u *profileUsecase) GetSiswaProfile(ctx context.Context, userID int64) (*domain.SiswaProfile, error) {
	return u.repo.FindSiswaProfile(ctx, userID)
}
