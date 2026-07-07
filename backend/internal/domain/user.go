package domain

import "context"

// Role merepresentasikan role user dalam sistem
type Role string

const (
	RoleAdmin Role = "admin"
	RoleGuru  Role = "guru"
	RoleSiswa Role = "siswa"
)

// User adalah entity utama (domain model), tidak boleh bergantung pada layer lain
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // hashed password, tidak pernah dikirim ke response
	Role     Role   `json:"role"`
}

// UserRepository adalah interface (port) yang harus diimplementasikan oleh layer repository
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

// AuthUsecase adalah interface (port) untuk business logic autentikasi
type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (string, *User, error)
	GetProfile(ctx context.Context, userID int64) (*User, error)
}
