-- Seed data contoh untuk melengkapi 4 user yang sudah ada di migration 000002.
-- Memakai subquery (SELECT id FROM users WHERE username = ...) alih-alih hardcode ID.

INSERT INTO
    classes (nama_kelas, tingkat)
VALUES ('XA', 10),
    ('XB', 10),
    ('XI RPL', 11);

INSERT INTO
    staffs (user_id, nama, posisi)
SELECT id, 'Petugas Satu', 'Kasir / Tata Usaha'
FROM users
WHERE
    username = 'petugas1';

INSERT INTO
    students (
        user_id,
        class_id,
        nisn,
        nama,
        alamat,
        no_telp
    )
SELECT id, (
        SELECT id
        FROM classes
        WHERE
            nama_kelas = 'XA'
    ), '0051234567', 'Siswa Satu', 'Jl. Merdeka No. 10', '081234567890'
FROM users
WHERE
    username = 'siswa1';

INSERT INTO
    spp (tahun_ajaran, nominal)
VALUES ('2025/2026', 150000.00);

INSERT INTO
    payments (
        staff_id,
        student_id,
        spp_id,
        bulan_dibayar,
        tanggal_bayar,
        jumlah_bayar
    )
SELECT (
        SELECT id
        FROM staffs
        WHERE
            nama = 'Petugas Satu'
    ), (
        SELECT id
        FROM students
        WHERE
            nisn = '0051234567'
    ), (
        SELECT id
        FROM spp
        WHERE
            tahun_ajaran = '2025/2026'
    ), 'Juli', CURRENT_DATE(), 150000.00;