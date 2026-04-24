package dto

type PPh21TahunanResponse struct {
	ID         int    `json:"id"`
	NIP        string `gorm:"column:nip"`
	NIK        string `gorm:"column:nik"`
	Nama       string `gorm:"column:nama"`
	Jabatan    string `gorm:"column:jabatan"`
	Golongan   string `gorm:"column:golongan"` // Disiapkan untuk CounterpartRank
	StatusPTKP string `gorm:"column:status_ptkp"`
	BulanAwal  int    `json:"BulanAwal"`
	BulanAkhir int    `json:"BulanAkhir"`

	// Rincian Penghasilan (Hasil SUM)
	TotalGajiPokok   float64 `gorm:"column:total_gaji_pokok"`
	TotalTunjIstri   float64 `gorm:"column:total_tunj_istri"`
	TotalTunjAnak    float64 `gorm:"column:total_tunj_anak"`
	TotalTunjBeras   float64 `gorm:"column:total_tunj_beras"`
	TotalTunjJabatan float64 `gorm:"column:total_tunj_jabatan"`
	TotalTunjLain    float64 `gorm:"column:total_tunj_lain"`
	TotalTPP         float64 `gorm:"column:total_tpp"`
	TotalBPJS4TPP    float64 `json:"TotalBPJS4TPP"`
	TotalIWP         float64 `gorm:"column:total_iwp"` // Untuk Iuran Pensiun

	// Hasil Kalkulasi PPh 21 (Diisi di Service)
	TotalBruto   float64
	BiayaJabatan float64
	IuranPensiun float64
	TotalNeto    float64
	NilaiPTKP    float64
	PKP          float64
	PPh21Setahun float64
}
