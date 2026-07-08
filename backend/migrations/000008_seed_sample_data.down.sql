DELETE FROM payments
WHERE
    student_id = (
        SELECT id
        FROM students
        WHERE
            nisn = '0051234567'
    );

DELETE FROM spp WHERE tahun_ajaran = '2025/2026';

DELETE FROM students WHERE nisn = '0051234567';

DELETE FROM staffs WHERE nama = 'Petugas Satu';

DELETE FROM classes WHERE nama_kelas IN ('XA', 'XB', 'XI RPL');