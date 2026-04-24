package model

type AkunBelanja struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	KodeRekening string  `gorm:"type:varchar(50);not null;uniqueIndex" json:"kode_rekening"`
	Uraian       string  `gorm:"type:varchar(255);not null" json:"uraian"`
	Kelompok     string  `gorm:"type:varchar(50)" json:"kelompok"`
	PaguAnggaran float64 `gorm:"type:decimal(15,2);default:0" json:"pagu_anggaran"`

	BaseModel
}

func (AkunBelanja) TableName() string {
	return "akun_belanja"
}
