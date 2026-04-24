package dto

type TppResponse struct {
	ID          uint    `json:"id"`
	Nama        string  `json:"nama"`
	NIP         string  `json:"nip"`
	BrutoTPP    float64 `json:"bruto_tpp"`
	PotonganTPP float64 `json:"potongan_tpp"`
	TPPBersih   float64 `json:"tpp_bersih"`
}
