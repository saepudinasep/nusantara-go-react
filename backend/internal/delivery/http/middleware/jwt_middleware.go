package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/pkg/jwt"
	"backend/pkg/response"
)

const (
	ContextUserID   = "userID"
	ContextUsername = "username"
	ContextRole     = "role"
)

// JWTAuthMiddleware memvalidasi bearer token di header Authorization.
// Jika valid, menyimpan userID/username/role ke dalam gin.Context untuk dipakai handler & middleware berikutnya.
func JWTAuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "header Authorization tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "format token harus 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		// simpan info user ke context, bisa diambil di handler via c.GetInt64(middleware.ContextUserID) dst
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

// RoleMiddleware membatasi akses endpoint hanya untuk role tertentu.
// Wajib dipasang SETELAH JWTAuthMiddleware karena bergantung pada context "role".
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextRole)
		if !exists {
			response.Error(c, http.StatusUnauthorized, "role tidak ditemukan pada token")
			c.Abort()
			return
		}

		role := roleValue.(string)
		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "anda tidak memiliki akses ke resource ini")
		c.Abort()
	}
}
