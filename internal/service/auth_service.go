package service

import (
	"backend-sigaji/internal/model" // Sesuaikan import kalau beda nama modul
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Kunci rahasia buat JWT (IDEALNYA TARUH DI .env!)
var jwtSecret = []byte("rahasia_negara_sigaji_2026")

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Fungsi untuk mencocokkan NIP dan Password
// Ubah signature-nya: balikin (string, model.User, error)
func (s *AuthService) LoginAttempt(username string, plainPassword string) (string, model.User, error) {
	var user model.User

	// 1. Cari user berdasarkan Username (NIP)
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", model.User{}, errors.New("username atau password salah")
	}

	// 2. Cocokkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	if err != nil {
		return "", model.User{}, errors.New("username atau password salah")
	}

	// 3. Buatkan Token JWT (Payload ditambahin sesuai kebutuhan)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", model.User{}, err
	}

	// Balikin token dan data user utuh
	return tokenString, user, nil
}

// Fungsi Bantuan buat Bikin Akun Admin (dipakai di main.go)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
