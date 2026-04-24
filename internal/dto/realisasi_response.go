package dto

type RealisasiResponse struct {
	ID             uint    `json:"id"`
	KodeRekening   string  `json:"kode_rekening"`
	Uraian         string  `json:"uraian"`
	Kelompok       string  `json:"kelompok"`
	PaguAnggaran   float64 `json:"pagu_anggaran"`
	TotalRealisasi float64 `json:"total_realisasi"`
	SisaAnggaran   float64 `json:"sisa_anggaran"`
	Persentase     float64 `json:"persentase"`
}
