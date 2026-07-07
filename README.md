# School App — React + Golang + MySQL (Clean Architecture)

Contoh project login multi-role (**admin, guru, siswa**) dengan JWT bearer token,
middleware role-based, dan dashboard sederhana per role.

## Struktur Project

```
project/
├── backend/                     # Golang (Gin) - Clean Architecture
│   ├── cmd/api/main.go          # entry point + dependency injection
│   ├── internal/
│   │   ├── domain/              # entity & interface (tidak bergantung layer lain)
│   │   ├── usecase/             # business logic (login, dsb)
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
        ├── routes/ProtectedRoute.jsx
        ├── components/DashboardLayout.jsx
        └── pages/               # Login, DashboardAdmin, DashboardGuru, DashboardSiswa
```

### Alur Clean Architecture (backend)

```
Handler (HTTP) → Usecase (business logic) → Repository (interface) → MySQL (implementasi)
                        ↑
                    Domain (entity + interface, inti aplikasi, tidak tahu Gin/MySQL sama sekali)
```

Dependency selalu mengarah ke `domain`. `usecase` hanya bergantung pada interface
`UserRepository`, bukan implementasi MySQL-nya — sehingga gampang diganti
(misal ke PostgreSQL atau mock untuk testing) tanpa mengubah business logic.

---

## 1. Setup Database

`golang-migrate` dipakai untuk migrasi, jadi tidak perlu redirect `<` (yang bermasalah di PowerShell) dan versi schema tercatat rapi (up/down).

### a. Buat database kosong dulu (golang-migrate tidak membuat database, hanya isi tabel di database yang sudah ada)

```powershell
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS school_app CHARACTER SET utf8mb4"
```

### b. Install CLI `migrate`

**Windows (via Scoop, direkomendasikan):**
```powershell
scoop install migrate
```

**Atau download binary langsung** dari halaman release resmi:
https://github.com/golang-migrate/migrate/releases
→ cari file `migrate.windows-amd64.zip`, ekstrak, taruh `migrate.exe` di folder yang ada di PATH (atau langsung di folder `backend/`, lalu panggil pakai `.\migrate.exe`).

Cek instalasi:
```powershell
migrate -version
```

### c. Jalankan migrasi

Dari dalam folder `backend/`:
```powershell
migrate -path migrations -database "mysql://root:PASSWORD_KAMU@tcp(127.0.0.1:3306)/school_app" up
```
Ganti `PASSWORD_KAMU` dengan password MySQL kamu (kalau tidak ada password, hapus `:PASSWORD_KAMU`).

Ini akan menjalankan file di `migrations/`:
- `000001_create_users_table.up.sql` → membuat tabel `users`
- `000002_seed_default_users.up.sql` → mengisi 3 user contoh (admin, guru, siswa)

Untuk rollback (undo migrasi terakhir):
```powershell
migrate -path migrations -database "mysql://root:PASSWORD_KAMU@tcp(127.0.0.1:3306)/school_app" down 1
```

3 user contoh yang akan ter-seed:

| Role  | Email               | Password      |
|-------|---------------------|---------------|
| admin | admin@sekolah.com   | password123   |
| guru  | guru@sekolah.com    | password123   |
| siswa | siswa@sekolah.com   | password123   |

---

## 2. Setup Backend (Golang)

```bash
cd backend
cp .env.example .env
# edit .env sesuaikan DB_USER / DB_PASSWORD / JWT_SECRET

go mod tidy      # download dependency (butuh koneksi internet)
go run cmd/api/main.go
```

Server berjalan di `http://localhost:8080`.

### Endpoint API

| Method | Endpoint                | Akses           | Keterangan                       |
|--------|--------------------------|-----------------|-----------------------------------|
| POST   | `/api/auth/login`       | Public          | Login, return JWT bearer token   |
| GET    | `/api/auth/me`          | Semua role login| Profil user yang sedang login    |
| GET    | `/api/dashboard`        | Semua role login| Dashboard generik                |
| GET    | `/api/admin/dashboard`  | admin saja      | Dilindungi `RoleMiddleware`      |
| GET    | `/api/guru/dashboard`   | guru saja       | Dilindungi `RoleMiddleware`      |
| GET    | `/api/siswa/dashboard`  | siswa saja      | Dilindungi `RoleMiddleware`      |

Contoh request login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sekolah.com","password":"password123"}'
```

Contoh memanggil endpoint terproteksi:

```bash
curl http://localhost:8080/api/admin/dashboard \
  -H "Authorization: Bearer <token_dari_login>"
```

Jika role tidak sesuai (misal siswa mengakses `/api/admin/dashboard`), middleware
`RoleMiddleware` akan mengembalikan `403 Forbidden`. Jika token tidak dikirim atau
invalid/expired, `JWTAuthMiddleware` mengembalikan `401 Unauthorized`.

---

## 3. Setup Frontend (React)

```bash
cd frontend
npm install
npm run dev
```

Buka `http://localhost:5173`. Sudah diverifikasi `npm install` dan `npm run build`
berjalan tanpa error.

Alur frontend:
1. `Login.jsx` submit ke `AuthContext.login()` → memanggil `POST /api/auth/login`.
2. Token & data user disimpan di `localStorage`, lalu redirect otomatis sesuai role
   (`/admin/dashboard`, `/guru/dashboard`, atau `/siswa/dashboard`).
3. `axiosClient.js` menyisipkan header `Authorization: Bearer <token>` di setiap
   request otomatis lewat interceptor.
4. `ProtectedRoute.jsx` mengecek: sudah login? role sesuai? kalau tidak → redirect
   ke `/login` atau `/unauthorized`.
5. Jika backend membalas `401` (token expired/invalid), interceptor otomatis
   logout dan redirect ke halaman login.

---

## Catatan Keamanan untuk Produksi

- Ganti `JWT_SECRET` di `.env` dengan random string yang panjang & kuat.
- Pertimbangkan menyimpan token di **httpOnly cookie** alih-alih `localStorage`
  untuk mitigasi XSS (di contoh ini pakai localStorage agar sederhana untuk belajar).
- Tambahkan rate limiting pada endpoint `/api/auth/login` untuk mencegah brute force.
- Aktifkan HTTPS di production, dan set `AllowOrigins` CORS sesuai domain frontend asli.
