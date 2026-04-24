package dto

// CreatePegawaiRequest digunakan untuk menangkap payload JSON dari frontend
type CreatePegawaiRequest struct {
	NIP        string `json:"nip" binding:"required"`
	NIK        string `json:"nik" binding:"required"`
	Nama       string `json:"nama" binding:"required"`
	Jabatan    string `json:"jabatan" binding:"required"`
	StatusAsn  uint8  `json:"status_asn" binding:"required"`
	Golongan   string `json:"golongan" binding:"required"`
	StatusPTKP string `json:"status_ptkp" binding:"required"`
	Alamat     string `json:"alamat"`
}
