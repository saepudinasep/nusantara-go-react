package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend/internal/domain"
)

type dashboardRepository struct {
	db *sql.DB
}

// NewDashboardRepository mengembalikan implementasi domain.DashboardRepository yang berbasis MySQL
func NewDashboardRepository(db *sql.DB) domain.DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) CountStaffs(ctx context.Context) (int64, error) {
	return r.scanCount(ctx, `SELECT COUNT(*) FROM staffs`)
}

func (r *dashboardRepository) CountGuru(ctx context.Context) (int64, error) {
	return r.scanCount(ctx, `SELECT COUNT(*) FROM users WHERE role = 'guru'`)
}

func (r *dashboardRepository) CountStudents(ctx context.Context) (int64, error) {
	return r.scanCount(ctx, `SELECT COUNT(*) FROM students`)
}

func (r *dashboardRepository) CountClasses(ctx context.Context) (int64, error) {
	return r.scanCount(ctx, `SELECT COUNT(*) FROM classes`)
}

func (r *dashboardRepository) CountPaymentsToday(ctx context.Context) (int64, error) {
	return r.scanCount(ctx, `SELECT COUNT(*) FROM payments WHERE tanggal_bayar = CURDATE()`)
}

func (r *dashboardRepository) CountPaidStudentsForMonth(ctx context.Context, monthName string, year int) (int64, error) {
	query := `SELECT COUNT(DISTINCT student_id) FROM payments
	          WHERE bulan_dibayar = ? AND YEAR(tanggal_bayar) = ?`
	return r.scanCount(ctx, query, monthName, year)
}

func (r *dashboardRepository) SumPaymentsInMonth(ctx context.Context, month int, year int) (float64, error) {
	query := `SELECT COALESCE(SUM(jumlah_bayar), 0) FROM payments
	          WHERE MONTH(tanggal_bayar) = ? AND YEAR(tanggal_bayar) = ?`

	var total float64
	err := r.db.QueryRowContext(ctx, query, month, year).Scan(&total)
	return total, err
}

// RecentActivitiesAdmin menggabungkan 3 sumber aktivitas (pembayaran, siswa baru, kelas baru)
// diurutkan berdasarkan waktu terbaru, dipakai khusus untuk dashboard admin.
func (r *dashboardRepository) RecentActivitiesAdmin(ctx context.Context, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT label, detail, created_at FROM (
			SELECT 'Pembayaran SPP diterima' AS label,
			       CONCAT(st.nama, ' · Rp', FORMAT(p.jumlah_bayar, 0, 'de_DE')) AS detail,
			       p.created_at
			FROM payments p
			JOIN students st ON st.id = p.student_id

			UNION ALL

			SELECT 'Siswa baru terdaftar', s.nama, s.created_at
			FROM students s

			UNION ALL

			SELECT 'Kelas baru ditambahkan', CONCAT(c.nama_kelas, ' · Tingkat ', c.tingkat), c.created_at
			FROM classes c
		) combined
		ORDER BY created_at DESC
		LIMIT ?
	`
	return r.scanActivities(ctx, query, limit)
}

// RecentPayments mengambil transaksi pembayaran terbaru (dipakai untuk dashboard petugas)
func (r *dashboardRepository) RecentPayments(ctx context.Context, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT 'Pembayaran SPP diterima' AS label,
		       CONCAT(st.nama, ' · Rp', FORMAT(p.jumlah_bayar, 0, 'de_DE')) AS detail,
		       p.created_at
		FROM payments p
		JOIN students st ON st.id = p.student_id
		ORDER BY p.created_at DESC
		LIMIT ?
	`
	return r.scanActivities(ctx, query, limit)
}

// RecentPaymentsByStudent mengambil riwayat pembayaran milik satu siswa saja (dipakai dashboard siswa)
func (r *dashboardRepository) RecentPaymentsByStudent(ctx context.Context, studentID int64, limit int) ([]domain.ActivityItem, error) {
	query := `
		SELECT 'Pembayaran SPP berhasil' AS label,
		       CONCAT('Rp', FORMAT(p.jumlah_bayar, 0, 'de_DE'), ' · ', p.bulan_dibayar) AS detail,
		       p.created_at
		FROM payments p
		WHERE p.student_id = ?
		ORDER BY p.created_at DESC
		LIMIT ?
	`
	return r.scanActivities(ctx, query, studentID, limit)
}

func (r *dashboardRepository) FindStudentByUserID(ctx context.Context, userID int64) (int64, string, int, error) {
	query := `SELECT s.id, c.nama_kelas, c.tingkat
	          FROM students s
	          JOIN classes c ON c.id = s.class_id
	          WHERE s.user_id = ? LIMIT 1`

	var studentID int64
	var namaKelas string
	var tingkat int

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&studentID, &namaKelas, &tingkat)
	if errors.Is(err, sql.ErrNoRows) {
		// User login sebagai siswa tapi belum punya profil students — bukan error fatal,
		// biarkan usecase yang memutuskan nilai default untuk kasus ini.
		return 0, "", 0, nil
	}
	if err != nil {
		return 0, "", 0, err
	}

	return studentID, namaKelas, tingkat, nil
}

func (r *dashboardRepository) HasStudentPaidForMonth(ctx context.Context, studentID int64, monthName string, year int) (bool, error) {
	query := `SELECT COUNT(*) FROM payments
	          WHERE student_id = ? AND bulan_dibayar = ? AND YEAR(tanggal_bayar) = ?`

	count, err := r.scanCount(ctx, query, studentID, monthName, year)
	return count > 0, err
}

func (r *dashboardRepository) SumPaymentsByStudent(ctx context.Context, studentID int64) (float64, error) {
	query := `SELECT COALESCE(SUM(jumlah_bayar), 0) FROM payments WHERE student_id = ?`

	var total float64
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&total)
	return total, err
}

// ---- helper internal ----

func (r *dashboardRepository) scanCount(ctx context.Context, query string, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scanCount: %w", err)
	}
	return count, nil
}

func (r *dashboardRepository) scanActivities(ctx context.Context, query string, args ...interface{}) ([]domain.ActivityItem, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []domain.ActivityItem{}
	for rows.Next() {
		var label, detail string
		var createdAt sql.NullTime
		if err := rows.Scan(&label, &detail, &createdAt); err != nil {
			return nil, err
		}
		sub := detail
		if createdAt.Valid {
			sub = detail + " · " + formatRelativeTime(createdAt.Time)
		}
		list = append(list, domain.ActivityItem{Label: label, Sub: sub})
	}

	return list, rows.Err()
}

// formatRelativeTime mengubah timestamp menjadi teks relatif ala "2 jam lalu" / "kemarin" / "3 hari lalu"
func formatRelativeTime(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return "baru saja"
	case d < time.Hour:
		mins := int(d.Minutes())
		return fmt.Sprintf("%d menit lalu", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		return fmt.Sprintf("%d jam lalu", hours)
	case d < 48*time.Hour:
		return "kemarin"
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%d hari lalu", days)
	default:
		return t.Format("2 Jan 2006")
	}
}
