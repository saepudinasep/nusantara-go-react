-- Seed 4 user contoh, 1 untuk masing-masing role sesuai skema V1.
-- Password plain untuk semuanya: "password123"
-- Hash bcrypt di bawah ini SUDAH diverifikasi valid untuk password tersebut.
INSERT INTO
    users (username, password, role)
VALUES (
        'admin',
        '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK',
        'admin'
    ),
    (
        'petugas1',
        '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK',
        'petugas'
    ),
    (
        'guru1',
        '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK',
        'guru'
    ),
    (
        'siswa1',
        '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK',
        'siswa'
    );