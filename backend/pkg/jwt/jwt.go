package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("token tidak valid atau sudah kadaluarsa")

// Claims custom, menyimpan user id, username, dan role di dalam JWT payload
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTService(secretKey string, tokenDurationHours int) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		tokenDuration: time.Duration(tokenDurationHours) * time.Hour,
	}
}

// GenerateToken membuat JWT bearer token berisi user_id, username, dan role
func (s *JWTService) GenerateToken(userID int64, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "school-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ParseToken memvalidasi dan mem-parsing token, mengembalikan Claims jika valid
func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
