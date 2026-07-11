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
	ErrSiswaNotFound      = errors.New("data siswa tidak ditemukan")
	ErrUsernameTaken      = errors.New("username sudah digunakan, pilih username lain")
	ErrNisnTaken          = errors.New("NISN sudah terdaftar untuk siswa lain")
	ErrKelasTidakValid    = errors.New("kelas yang dipilih tidak valid")
	ErrSiswaInUse         = errors.New("siswa tidak dapat dihapus karena masih memiliki riwayat transaksi pembayaran")
	ErrPetugasNotFound    = errors.New("data petugas tidak ditemukan")
	ErrPetugasInUse       = errors.New("petugas tidak dapat dihapus karena masih memiliki riwayat transaksi pembayaran")
	ErrPaymentNotFound    = errors.New("data transaksi tidak ditemukan")
	ErrPaymentDuplicate   = errors.New("siswa ini sudah tercatat membayar SPP untuk bulan dan jenis SPP tersebut")
	ErrStudentInvalid     = errors.New("siswa dengan NISN tersebut tidak ditemukan")
	ErrSppInvalid         = errors.New("data SPP yang dipilih tidak valid")
	ErrStaffInvalid       = errors.New("petugas yang dipilih tidak valid")
)
