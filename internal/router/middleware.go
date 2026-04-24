package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Samakan dengan secret key di auth_service
var jwtSecret = []byte("rahasia_negara_sigaji_2026")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil Header "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses Ditolak. Token tidak ditemukan."})
			c.Abort()
			return
		}

		// 2. Pisahkan kata "Bearer " dan tokennya
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 3. Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signature tidak valid")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah expired"})
			c.Abort()
			return
		}

		// 4. Ekstrak isi token (NIP, Role) lalu simpan di Context
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
		}

		// Lolos dari satpam, silakan masuk ke controller!
		c.Next()
	}
}
