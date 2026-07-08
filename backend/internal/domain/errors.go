package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("anda tidak memiliki akses ke resource ini")
	ErrKelasNotFound      = errors.New("kelas tidak ditemukan")
	ErrDuplicateEntry     = errors.New("data dengan nama_kelas tersebut sudah ada")
)
