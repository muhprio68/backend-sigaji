package config

import (
	"backend-sigaji/internal/model"
	"backend-sigaji/internal/service"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	dsn := "root:@tcp(127.0.0.1:3306)/sigaji-db?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	log.Println("Database connected")
	db.AutoMigrate(
		&model.Pegawai{},
		&model.GajiPokok{},
		&model.TPP{},
		&model.AkunBelanja{},
		&model.RealisasiBelanja{},
		&model.User{},
	)

	// SEEDER ADMIN OTOMATIS:
	// Cek apakah tabel users masih kosong, kalau iya, bikinin 1 akun admin.
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		hashedPwd, _ := service.HashPassword("admin123") // Password rahasianya: admin123
		admin := model.User{
			Username: "admin",
			Nama:     "Muhammad Prio Agustian",
			Email:    "muhammadprio93@gmail.com",
			Jabatan:  "Penata Kelola Sistem dan Teknologi Informasi",
			Password: hashedPwd,
			Role:     "admin",
		}
		db.Create(&admin)
		fmt.Println("🚀 Akun Super Admin berhasil dibuat! (User: admin, Pass: admin123)")
	}

	return db
}
