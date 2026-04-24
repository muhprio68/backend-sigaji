package model

type RealisasiBelanja struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Relasi ke Master Akun
	IDAkunBelanja uint        `gorm:"not null;uniqueIndex:idx_realisasi_periode" json:"id_akun_belanja"`
	AkunBelanja   AkunBelanja `gorm:"foreignKey:IDAkunBelanja;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	// Periode & Jenis
	Bulan            int    `gorm:"type:tinyint;not null;check:bulan >= 1 AND bulan <= 12;uniqueIndex:idx_realisasi_periode" json:"bulan"`
	Tahun            int    `gorm:"type:year;not null;uniqueIndex:idx_realisasi_periode" json:"tahun"`
	JenisPenghasilan string `gorm:"type:varchar(25);not null;uniqueIndex:idx_realisasi_periode" json:"jenis_penghasilan"` // bulanan, thr, gaji_13

	// Sumber Data (gaji_pokok atau tpp)
	Sumber string `gorm:"type:varchar(20);not null;uniqueIndex:idx_realisasi_periode" json:"sumber"`

	// Nominal Uang yang Cair
	Realisasi float64 `gorm:"type:decimal(15,2);default:0" json:"realisasi"`

	BaseModel
}

func (RealisasiBelanja) TableName() string {
	return "realisasi_belanja"
}
