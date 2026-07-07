package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("anda tidak memiliki akses ke resource ini")
)
