# School App — React + Golang + MySQL (Clean Architecture)

Contoh project login multi-role (**admin, petugas, guru, siswa**) dengan JWT bearer token,
login berbasis **username**, middleware role-based, dashboard per role, dan CRUD data Kelas.

Skema database mengikuti rancangan **V1 (SPP & Master)**.

## Struktur Project

```
project/
├── backend/                     # Golang (Gin) - Clean Architecture
│   ├── cmd/api/main.go          # entry point + dependency injection
│   ├── internal/
│   │   ├── domain/              # entity & interface (tidak bergantung layer lain)
│   │   ├── usecase/             # business logic (login, CRUD kelas, dsb)
│   │   ├── repository/mysql/    # implementasi akses data ke MySQL
│   │   ├── delivery/http/
│   │   │   ├── handler/         # HTTP handler (controller)
│   │   │   ├── middleware/      # JWT auth & role middleware
│   │   │   └── router/          # route registration
│   │   └── config/              # load .env & koneksi DB
│   ├── pkg/                     # jwt, hash, response helper (reusable, tidak spesifik domain)
│   ├── migrations/              # golang-migrate: up/down SQL files
│   └── .env.example
│
└── frontend/                    # React (Vite)
    └── src/
        ├── api/axiosClient.js   # axios instance + interceptor bearer token
        ├── context/AuthContext.jsx
        ├── routes/ProtectedRoute.jsx, GuestRoute.jsx
        ├── components/Sidebar.jsx, Topbar.jsx, DashboardLayout.jsx
        └── pages/               # Login, DashboardAdmin, DashboardPetugas, DashboardGuru, DashboardSiswa, KelasManagement
```

### Alur Clean Architecture (backend)

```
Handler (HTTP) → Usecase (business logic) → Repository (interface) → MySQL (implementasi)
                        ↑
                    Domain (entity + interface, inti aplikasi, tidak tahu Gin/MySQL sama sekali)
```

---

## 1. Skema Database (V1 — SPP & Master)

| Tabel      | Kolom                                                                                       | Keterangan                                                                   |
| ---------- | ------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `users`    | id, **username**, password, role                                                            | Login pakai username, bukan email. Role: `admin`, `petugas`, `guru`, `siswa` |
| `classes`  | id, nama_kelas, tingkat (INT)                                                               | Contoh: "XA", tingkat 10                                                     |
| `students` | id, user_id (FK), class_id (FK), nisn, nama, alamat, no_telp                                | Profil siswa, 1-1 ke `users`                                                 |
| `staffs`   | id, user_id (FK), nama, posisi                                                              | Profil petugas SPP, 1-1 ke `users`                                           |
| `spp`      | id, tahun_ajaran, nominal                                                                   | Master nominal SPP per tahun ajaran                                          |
| `payments` | id, staff_id (FK), student_id (FK), spp_id (FK), bulan_dibayar, tanggal_bayar, jumlah_bayar | Transaksi pembayaran                                                         |

Role `guru` sudah tersedia di tabel `users` untuk kebutuhan V2 (Learning Center), tapi belum
punya tabel profil atau fitur khusus di V1 — dashboard guru saat ini masih placeholder.

---

## 2. Setup Database

### a. Buat database kosong

```powershell
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS school_app CHARACTER SET utf8mb4"
```

### b. Install CLI `migrate`

**Windows (via Scoop, direkomendasikan):**

```powershell
scoop install migrate
```

Atau download binary dari https://github.com/golang-migrate/migrate/releases (`migrate.windows-amd64.zip`).

### c. Jalankan migrasi

Dari dalam folder `backend/`:

```powershell
migrate -path migrations -database "mysql://root:PASSWORD_KAMU@tcp(127.0.0.1:3306)/school_app" up
```

| File                           | Keterangan                                                     |
| ------------------------------ | -------------------------------------------------------------- |
| `000001_create_users_table`    | Tabel `users` (username, password, role)                       |
| `000002_seed_default_users`    | 4 user contoh: admin, petugas1, guru1, siswa1                  |
| `000003_create_classes_table`  | Tabel `classes`                                                |
| `000004_create_students_table` | Tabel `students`                                               |
| `000005_create_staffs_table`   | Tabel `staffs`                                                 |
| `000006_create_spp_table`      | Tabel `spp`                                                    |
| `000007_create_payments_table` | Tabel `payments`                                               |
| `000008_seed_sample_data`      | Data contoh: 3 kelas, profil staff & student, 1 SPP, 1 payment |

4 user contoh yang ter-seed (password sama untuk semua):

| Role    | Username   | Password      |
| ------- | ---------- | ------------- |
| admin   | `admin`    | `password123` |
| petugas | `petugas1` | `password123` |
| guru    | `guru1`    | `password123` |
| siswa   | `siswa1`   | `password123` |

Rollback:

```powershell
migrate -path migrations -database "mysql://root:PASSWORD_KAMU@tcp(127.0.0.1:3306)/school_app" down
```

---

## 3. Setup Backend (Golang)

```bash
cd backend
cp .env.example .env
# edit .env sesuaikan DB_USER / DB_PASSWORD / JWT_SECRET

go mod tidy
go run cmd/api/main.go
```

Server berjalan di `http://localhost:8080`.

> **Catatan penting soal MySQL user:** kalau kamu login MySQL sebagai `root` tanpa password
> lewat `auth_socket`/`unix_socket` plugin (default di banyak instalasi Linux/MariaDB), koneksi
> TCP dari aplikasi Go bisa ditolak (`Access denied`). Kalau ini terjadi, buat user MySQL khusus:
>
> ```sql
> CREATE USER 'appuser'@'%' IDENTIFIED BY 'apppassword';
> GRANT ALL PRIVILEGES ON school_app.* TO 'appuser'@'%';
> FLUSH PRIVILEGES;
> ```
>
> lalu pakai `appuser`/`apppassword` di `.env`.

### Endpoint API

| Method         | Endpoint                 | Akses            | Keterangan                                      |
| -------------- | ------------------------ | ---------------- | ----------------------------------------------- |
| POST           | `/api/auth/login`        | Public           | Login pakai `username` + `password`, return JWT |
| GET            | `/api/auth/me`           | Semua role login | Profil user yang sedang login                   |
| GET            | `/api/admin/dashboard`   | admin            | —                                               |
| GET            | `/api/petugas/dashboard` | petugas          | —                                               |
| GET            | `/api/guru/dashboard`    | guru             | —                                               |
| GET            | `/api/siswa/dashboard`   | siswa            | —                                               |
| GET/POST       | `/api/admin/kelas`       | admin            | List / tambah kelas                             |
| GET/PUT/DELETE | `/api/admin/kelas/:id`   | admin            | Detail / update / hapus kelas (hard delete)     |

Contoh login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

---

## 4. Setup Frontend (React)

```bash
cd frontend
npm install
npm run dev
```

Buka `http://localhost:5173`. Form login sekarang meminta **username**, bukan email.

---

## Validasi yang Sudah Dilakukan

Konfigurasi ini **sudah diuji end-to-end secara nyata** (bukan cuma dibaca sekilas):

1. **Migration**: seluruh `up`/`down` dijalankan di MariaDB sungguhan, relasi FK diverifikasi lewat query JOIN, rollback total berhasil bersih.
2. **Backend Go**: berhasil di-_compile_ jadi binary asli dan dijalankan sebagai server sungguhan (bukan cuma cek sintaks). Semua endpoint dites lewat `curl` nyata:
   - Login dengan `username` untuk 4 role → sukses, JWT valid.
   - CRUD Kelas penuh (Create, List, Update, Delete) → sukses.
   - Role-guard: siswa coba akses `/api/admin/kelas` → `403 Forbidden` seperti seharusnya.
   - Validasi duplikat: bikin kelas dengan `nama_kelas` yang sama → `409 Conflict`.
   - Dashboard ke-4 role → semua merespons data yang benar.
3. **Bug nyata ditemukan & diperbaiki dari testing ini**: endpoint `PUT /api/admin/kelas/:id` salah mengembalikan `404 Not Found` ketika data yang dikirim **identik** dengan data yang sudah tersimpan (MySQL driver Go secara default melaporkan _rows changed_, bukan _rows matched_, di `RowsAffected()`). Sudah diperbaiki dengan menambahkan parameter `clientFoundRows=true` di DSN koneksi (`internal/config/database.go`), dan sudah dites ulang untuk memastikan fix-nya benar sekaligus tidak merusak kasus 404 yang seharusnya (ID yang benar-benar tidak ada tetap mengembalikan 404).
4. **Frontend**: `npm install` + `npm run build` sukses tanpa error.

---

## Catatan Keamanan untuk Produksi

- Ganti `JWT_SECRET` di `.env` dengan random string yang panjang & kuat.
- Pertimbangkan menyimpan token di **httpOnly cookie** alih-alih `localStorage` untuk mitigasi XSS.
- Tambahkan rate limiting pada endpoint `/api/auth/login` untuk mencegah brute force.
- Aktifkan HTTPS di production, dan set `AllowOrigins` CORS sesuai domain frontend asli.
