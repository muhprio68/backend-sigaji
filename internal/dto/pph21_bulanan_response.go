package dto

type PPh21BulananResponse struct {
	ID         uint    `json:"id"`
	NIP        string  `gorm:"column:nip" json:"nip"`
	NIK        string  `gorm:"column:nik" json:"nik"`
	Nama       string  `json:"nama"`
	Jabatan    string  `gorm:"column:jabatan" json:"jabatan"`
	StatusPTKP string  `gorm:"column:status_ptkp" json:"status_ptkp"`
	BulanPajak int     `json:"bulan_pajak" gorm:"column:bulan_pajak"`
	BrutoGaji  float64 `json:"bruto_gaji"`
	BrutoTPP   float64 `json:"bruto_tpp"`
	TotalBruto float64 `json:"total_bruto"`
	TarifTER   float64 `json:"tarif_ter"`
	PPh21      float64 `json:"pph21"`
}

// DTO Baru untuk membungkus semuanya
type PPh21PageResponse struct {
	Summary struct {
		TotalPegawai    int     `json:"total_pegawai"`
		GrandTotalBruto float64 `json:"grand_total_bruto"`
		GrandTotalPPh21 float64 `json:"grand_total_pph21"`
	} `json:"summary"`
	ListPegawai []PPh21BulananResponse `json:"list_pegawai"`
}
