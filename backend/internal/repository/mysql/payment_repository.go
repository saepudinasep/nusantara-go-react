package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"backend/internal/domain"
)

type paymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository mengembalikan implementasi domain.PaymentRepository yang berbasis MySQL
func NewPaymentRepository(db *sql.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *domain.Payment) (int64, error) {
	query := `INSERT INTO payments (staff_id, student_id, spp_id, bulan_dibayar, tanggal_bayar, jumlah_bayar)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, p.StaffID, p.StudentID, p.SppID, p.BulanDibayar, p.TanggalBayar, p.JumlahBayar)
	if err != nil {
		if isLockWaitTimeoutError(err) {
			return 0, domain.ErrDatabaseBusy
		}
		if isForeignKeyChildRowError(err) {
			// Petakan ke error yang spesifik sesuai FK mana yang gagal, supaya pesannya
			// tepat sasaran (bukan cuma "data tidak valid" generik).
			switch {
			case strings.Contains(err.Error(), "fk_payments_student"):
				return 0, domain.ErrStudentInvalid
			case strings.Contains(err.Error(), "fk_payments_spp"):
				return 0, domain.ErrSppInvalid
			case strings.Contains(err.Error(), "fk_payments_staff"):
				return 0, domain.ErrStaffInvalid
			}
		}
		return 0, err
	}

	return result.LastInsertId()
}

func (r *paymentRepository) FindAll(ctx context.Context, page, limit int, staffID *int64, filter domain.PaymentFilter) ([]domain.Payment, int64, error) {
	conditions := []string{}
	args := []interface{}{}

	if staffID != nil {
		conditions = append(conditions, "p.staff_id = ?")
		args = append(args, *staffID)
	}
	if filter.Search != "" {
		conditions = append(conditions, "(st.nama LIKE ? OR st.nisn LIKE ?)")
		like := "%" + filter.Search + "%"
		args = append(args, like, like)
	}
	if filter.BulanDibayar != "" {
		conditions = append(conditions, "p.bulan_dibayar = ?")
		args = append(args, filter.BulanDibayar)
	}
	if filter.TanggalDari != "" {
		conditions = append(conditions, "p.tanggal_bayar >= ?")
		args = append(args, filter.TanggalDari)
	}
	if filter.TanggalSampai != "" {
		conditions = append(conditions, "p.tanggal_bayar <= ?")
		args = append(args, filter.TanggalSampai)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// JOIN ke students WAJIB ada di count query juga karena filter Search membaca kolom st.nama/st.nisn
	joins := `
		JOIN staffs sf ON sf.id = p.staff_id
		JOIN students st ON st.id = p.student_id
		JOIN classes c ON c.id = st.class_id
		JOIN spp sp ON sp.id = p.spp_id
	`

	var total int64
	countQuery := `SELECT COUNT(*) FROM payments p ` + joins + whereClause
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT p.id, p.staff_id, sf.nama, p.student_id, st.nama, st.nisn, c.nama_kelas,
		       p.spp_id, sp.tahun_ajaran, p.bulan_dibayar, p.tanggal_bayar, p.jumlah_bayar,
		       p.created_at, p.updated_at
		FROM payments p
		` + joins + whereClause + `
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []domain.Payment{}
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.StaffID, &p.StaffNama, &p.StudentID, &p.StudentNama, &p.Nisn, &p.NamaKelas,
			&p.SppID, &p.TahunAjaran, &p.BulanDibayar, &p.TanggalBayar, &p.JumlahBayar,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}

	return list, total, rows.Err()
}

func (r *paymentRepository) FindByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := `
		SELECT p.id, p.staff_id, sf.nama, p.student_id, st.nama, st.nisn, c.nama_kelas,
		       p.spp_id, sp.tahun_ajaran, p.bulan_dibayar, p.tanggal_bayar, p.jumlah_bayar,
		       p.created_at, p.updated_at
		FROM payments p
		JOIN staffs sf ON sf.id = p.staff_id
		JOIN students st ON st.id = p.student_id
		JOIN classes c ON c.id = st.class_id
		JOIN spp sp ON sp.id = p.spp_id
		WHERE p.id = ?
		LIMIT 1
	`

	var p domain.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.StaffID, &p.StaffNama, &p.StudentID, &p.StudentNama,
		&p.Nisn, &p.NamaKelas, &p.SppID, &p.TahunAjaran, &p.BulanDibayar, &p.TanggalBayar, &p.JumlahBayar,
		&p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Delete = "batalkan transaksi". Skema payments belum punya kolom status/canceled_at,
// jadi pembatalan diimplementasikan sebagai hard delete. Tidak ada tabel lain yang
// mereferensikan payments di skema V1, jadi tidak ada risiko FK constraint saat delete ini.
func (r *paymentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM payments WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		if isLockWaitTimeoutError(err) {
			return domain.ErrDatabaseBusy
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}

// HasPaidForPeriod mengecek apakah siswa sudah pernah membayar SPP tertentu untuk bulan yang sama —
// dipakai usecase untuk mencegah pencatatan pembayaran ganda pada periode yang sama.
func (r *paymentRepository) HasPaidForPeriod(ctx context.Context, studentID, sppID int64, bulanDibayar string) (bool, error) {
	query := `SELECT COUNT(*) FROM payments WHERE student_id = ? AND spp_id = ? AND bulan_dibayar = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, studentID, sppID, bulanDibayar).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindAllByStudent dipakai halaman "Riwayat Pembayaran" siswa — HANYA transaksi milik siswa itu sendiri.
func (r *paymentRepository) FindAllByStudent(ctx context.Context, studentID int64, page, limit int) ([]domain.Payment, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments WHERE student_id = ?`, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT p.id, p.staff_id, sf.nama, p.student_id, st.nama, st.nisn, c.nama_kelas,
		       p.spp_id, sp.tahun_ajaran, p.bulan_dibayar, p.tanggal_bayar, p.jumlah_bayar,
		       p.created_at, p.updated_at
		FROM payments p
		JOIN staffs sf ON sf.id = p.staff_id
		JOIN students st ON st.id = p.student_id
		JOIN classes c ON c.id = st.class_id
		JOIN spp sp ON sp.id = p.spp_id
		WHERE p.student_id = ?
		ORDER BY p.tanggal_bayar DESC, p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	list := []domain.Payment{}
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.StaffID, &p.StaffNama, &p.StudentID, &p.StudentNama, &p.Nisn, &p.NamaKelas,
			&p.SppID, &p.TahunAjaran, &p.BulanDibayar, &p.TanggalBayar, &p.JumlahBayar,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, p)
	}

	return list, total, rows.Err()
}

// FindByStudentAndSpp dipakai untuk menyusun status Lunas/Belum per bulan pada satu jenis SPP —
// TANPA pagination karena hasilnya paling banyak 12 baris (1 per bulan).
func (r *paymentRepository) FindByStudentAndSpp(ctx context.Context, studentID, sppID int64) ([]domain.Payment, error) {
	query := `
		SELECT id, staff_id, student_id, spp_id, bulan_dibayar, tanggal_bayar, jumlah_bayar, created_at, updated_at
		FROM payments
		WHERE student_id = ? AND spp_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, studentID, sppID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []domain.Payment{}
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.StaffID, &p.StudentID, &p.SppID, &p.BulanDibayar, &p.TanggalBayar,
			&p.JumlahBayar, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, rows.Err()
}
