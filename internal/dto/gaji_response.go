package dto

type GajiResponse struct {
	ID           uint    `json:"id"`
	Nama         string  `json:"nama"`
	NIP          string  `json:"nip"`
	BrutoGaji    float64 `json:"bruto_gaji"`
	PotonganGaji float64 `json:"potongan_gaji"`
	GajiBersih   float64 `json:"gaji_bersih"`
}
