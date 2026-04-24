package dto

type DetailTppResponse struct {
	ID               uint   `json:"id"`
	Nama             string `json:"nama"`
	NIP              string `json:"nip"`
	Bulan            int    `json:"bulan"`
	Tahun            int    `json:"tahun"`
	JenisPenghasilan string `json:"jenis_penghasilan"`

	BebanKerja    float64 `json:"beban_kerja"`
	PrestasiKerja float64 `json:"prestasi_kerja"`
	KondisiKerja  float64 `json:"kondisi_kerja"`
	TunjBpjs4     float64 `json:"tunj_bpjs4"`
	TunjPph21     float64 `json:"tunj_pph21"`
	BrutoTpp      float64 `json:"bruto_tpp"`
	PotPph21      float64 `json:"pot_pph21"`
	PotBpjs4      float64 `json:"pot_bpjs4"`
	PotBpjs1      float64 `json:"pot_bpjs1"`
	TppBersih     float64 `json:"tpp_bersih"`
}
