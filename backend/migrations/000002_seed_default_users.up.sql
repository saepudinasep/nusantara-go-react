-- Seed 3 user contoh, 1 untuk masing-masing role.
-- Password plain untuk ketiganya: "password123"
-- Hash bcrypt di bawah ini SUDAH diverifikasi valid untuk password tersebut.
INSERT INTO users (name, email, password, role) VALUES
('Admin Sekolah', 'admin@sekolah.com', '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK', 'admin'),
('Budi Guru', 'guru@sekolah.com', '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK', 'guru'),
('Siti Siswa', 'siswa@sekolah.com', '$2b$10$VjIkvf.t.bkW1LX1AaWPOeipgi2rG5gVWOMEx1IG3nisl80rOO0bK', 'siswa');
