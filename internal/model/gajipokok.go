package model

type GajiPokok struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	IDPegawai uint    `gorm:"not null;uniqueIndex:idx_pegawai_periode" json:"id_pegawai"`
	Pegawai   Pegawai `gorm:"foreignKey:IDPegawai;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Bulan            int    `gorm:"type:tinyint;not null;check:bulan >= 1 AND bulan <= 12;uniqueIndex:idx_pegawai_periode" json:"bulan"`
	Tahun            int    `gorm:"type:year;not null;uniqueIndex:idx_pegawai_periode" json:"tahun"`
	JenisPenghasilan string `gorm:"type:varchar(25);not null;uniqueIndex:idx_pegawai_periode" json:"jenis_penghasilan"`

	// --- 1. KOMPONEN GAJI & TUNJANGAN (Sesuai Kolom Excel) ---
	GajiPokok          float64 `gorm:"type:decimal(15,2);default:0" json:"gaji_pokok"`
	TunjSuamiIstri     float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_suami_istri"` // Akan digabung dgn Anak saat masuk Realisasi
	TunjAnak           float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_anak"`
	TunjJabatan        float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_jabatan"`
	TunjFungsional     float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_fungsional"`
	TunjFungsionalUmum float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_fungsional_umum"`
	TunjBeras          float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_beras"`
	TunjPph            float64 `gorm:"type:decimal(15,2);default:0" json:"tunj_pph"`
	Pembulatan         float64 `gorm:"type:decimal(15,2);default:0" json:"pembulatan"`

	// --- 2. KOMPONEN IURAN PEMDA (Penambah Bruto PPh 21) ---
	BpjsKesPemda float64 `gorm:"type:decimal(15,2);default:0" json:"bpjs_kes_pemda"`
	JkkPemda     float64 `gorm:"type:decimal(15,2);default:0" json:"jkk_pemda"`
	JkmPemda     float64 `gorm:"type:decimal(15,2);default:0" json:"jkm_pemda"`

	// --- 3. KOMPONEN POTONGAN PEGAWAI (Pengurang Pajak & Take Home Pay) ---
	Iwp8Persen     float64 `gorm:"type:decimal(15,2);default:0" json:"iwp_8_persen"` // Pensiun & JHT (Masuk ke XML A2)
	Iwp1Persen     float64 `gorm:"type:decimal(15,2);default:0" json:"iwp_1_persen"` // BPJS Kesehatan Pegawai
	PotonganTapera float64 `gorm:"type:decimal(15,2);default:0" json:"potongan_tapera"`
	PotonganLain   float64 `gorm:"type:decimal(15,2);default:0" json:"potongan_lain"` // Wadah untuk potongan Bank/Koperasi (Non-Pajak)

	BaseModel
}

func (GajiPokok) TableName() string {
	return "gaji_pokok"
}
