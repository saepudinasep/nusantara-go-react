CREATE TABLE IF NOT EXISTS payments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NOT NULL,
    student_id INT NOT NULL,
    spp_id INT NOT NULL,
    bulan_dibayar VARCHAR(20) NOT NULL,
    tanggal_bayar DATE NOT NULL,
    jumlah_bayar DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_payments_staff FOREIGN KEY (staff_id) REFERENCES staffs (id) ON UPDATE CASCADE,
    CONSTRAINT fk_payments_student FOREIGN KEY (student_id) REFERENCES students (id) ON UPDATE CASCADE,
    CONSTRAINT fk_payments_spp FOREIGN KEY (spp_id) REFERENCES spp (id) ON UPDATE CASCADE
) ENGINE = InnoDB;