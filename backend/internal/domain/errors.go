package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("username atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("anda tidak memiliki akses ke resource ini")
	ErrKelasNotFound      = errors.New("kelas tidak ditemukan")
	ErrDuplicateEntry     = errors.New("data dengan nama_kelas tersebut sudah ada")
	ErrKelasInUse         = errors.New("kelas tidak dapat dihapus karena masih memiliki data siswa yang terhubung")
	ErrDatabaseBusy       = errors.New("sistem sedang sibuk memproses operasi lain, silakan coba lagi dalam beberapa saat")
	ErrSppNotFound        = errors.New("data SPP tidak ditemukan")
	ErrSppDuplicate       = errors.New("data SPP untuk tahun ajaran tersebut sudah ada")
	ErrSppInUse           = errors.New("data SPP tidak dapat dihapus karena masih memiliki transaksi pembayaran yang terhubung")
)
