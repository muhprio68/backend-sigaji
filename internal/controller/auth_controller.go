package controller

import (
	"backend-sigaji/internal/dto"     // Sesuaikan
	"backend-sigaji/internal/service" // Sesuaikan
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Tangkap token DAN data user dari Service
	token, user, err := c.authService.LoginAttempt(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Kirim response lengkap
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"username": user.Username,
			"nama":     user.Nama,
			"email":    user.Email,
			"jabatan":  user.Jabatan,
			"role":     user.Role,
		},
	})
}
