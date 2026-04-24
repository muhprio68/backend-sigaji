package model

type TPP struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	IDPegawai uint    `gorm:"not null;uniqueIndex:idx_tpp_periode" json:"id_pegawai"`
	Pegawai   Pegawai `gorm:"foreignKey:IDPegawai;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Bulan            int    `gorm:"type:tinyint;not null;check:bulan >= 1 AND bulan <= 12;uniqueIndex:idx_tpp_periode" json:"bulan"`
	Tahun            int    `gorm:"type:year;not null;uniqueIndex:idx_tpp_periode" json:"tahun"`
	JenisPenghasilan string `gorm:"type:varchar(25);not null;uniqueIndex:idx_tpp_periode" json:"jenis_penghasilan"`

	TPPBeban    float64 `gorm:"type:decimal(15,2);default:0" json:"tpp_beban"`
	TPPPrestasi float64 `gorm:"type:decimal(15,2);default:0" json:"tpp_prestasi"`
	TPPKondisi  float64 `gorm:"type:decimal(15,2);default:0" json:"tpp_kondisi"`
	TunjPajak   float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_pajak"`
	PotPajak    float64 `gorm:"type:decimal(15,2);default:0" json:"pot_pajak"`

	BPJS4 float64 `gorm:"type:decimal(15,2);default:0" json:"bpjs_4_persen"`
	BPJS1 float64 `gorm:"type:decimal(15,2);default:0" json:"bpjs_1_persen"`

	BaseModel
}

func (TPP) TableName() string {
	return "tpp"
}
