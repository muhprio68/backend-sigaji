package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"` // NIP
	Password  string    `gorm:"not null" json:"-"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	Jabatan   string    `json:"jabatan"`
	Role      string    `json:"role"` // Isinya: "admin", "pegawai_pns", atau "pegawai_pppk"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
