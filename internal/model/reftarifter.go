package model

type RefTarifTer struct {
	ID              int     `gorm:"primaryKey"`
	Kategori        string  `gorm:"column:kategori"`
	BatasBawah      float64 `gorm:"column:batas_bawah"`
	BatasAtas       float64 `gorm:"column:batas_atas"`
	TarifPersentase float64 `gorm:"column:tarif_persentase"`
}

func (RefTarifTer) TableName() string {
	return "ref_tarif_ter"
}
