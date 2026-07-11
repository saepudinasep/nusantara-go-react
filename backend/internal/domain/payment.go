package domain

import "context"

// Payment merepresentasikan satu transaksi pembayaran SPP. Field ...Nama/Nisn/TahunAjaran
// bersifat read-only hasil JOIN (dipakai untuk tampilan), yang benar-benar disimpan ke tabel
// payments hanya StaffID, StudentID, SppID, BulanDibayar, TanggalBayar, JumlahBayar.
type Payment struct {
	ID           int64   `json:"id"`
	StaffID      int64   `json:"staff_id"`
	StaffNama    string  `json:"staff_nama,omitempty"`
	StudentID    int64   `json:"student_id"`
	StudentNama  string  `json:"student_nama,omitempty"`
	Nisn         string  `json:"nisn,omitempty"`
	NamaKelas    string  `json:"nama_kelas,omitempty"`
	SppID        int64   `json:"spp_id"`
	TahunAjaran  string  `json:"tahun_ajaran,omitempty"`
	BulanDibayar string  `json:"bulan_dibayar"`
	TanggalBayar string  `json:"tanggal_bayar"`
	JumlahBayar  float64 `json:"jumlah_bayar"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

// PaymentRepository adalah interface (port) yang harus diimplementasikan oleh layer repository.
// staffID pada FindAll bersifat opsional: nil berarti tanpa filter (dipakai admin, lihat semua
// transaksi), sedangkan nilai non-nil membatasi hasil hanya milik petugas tersebut.
type PaymentRepository interface {
	Create(ctx context.Context, p *Payment) (int64, error)
	FindAll(ctx context.Context, page, limit int, staffID *int64) ([]Payment, int64, error)
	FindByID(ctx context.Context, id int64) (*Payment, error)
	Delete(ctx context.Context, id int64) error
	HasPaidForPeriod(ctx context.Context, studentID, sppID int64, bulanDibayar string) (bool, error)
}

// PaymentUsecase adalah interface (port) untuk business logic pemrosesan pembayaran
type PaymentUsecase interface {
	Create(ctx context.Context, p *Payment) (*Payment, error)
	GetAll(ctx context.Context, page, limit int, staffID *int64) ([]Payment, Pagination, error)
	GetByID(ctx context.Context, id int64) (*Payment, error)
	Delete(ctx context.Context, id int64) error
}
