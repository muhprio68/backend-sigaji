package model

type Pegawai struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	NIP        string `gorm:"column:nip;type:varchar(20);uniqueIndex;not null" json:"nip"`
	NIK        string `gorm:"column:nik;type:varchar(20);uniqueIndex;not null" json:"nik"`
	Nama       string `gorm:"type:varchar(100);not null" json:"nama"`
	Jabatan    string `gorm:"type:varchar(100);not null" json:"jabatan"`
	StatusAsn  uint8  `gorm:"type:smallint;not null" json:"status_asn"`
	Golongan   string `gorm:"type:varchar(5);not null" json:"golongan"`
	StatusPTKP string `gorm:"column:status_ptkp;type:varchar(10);not null" json:"status_ptkp"`
	Alamat     string `gorm:"type:text" json:"alamat"`

	BaseModel
}

const (
	StatusPNS  uint8 = 1
	StatusPPPK uint8 = 2
	StatusCPNS uint8 = 3
)

func (Pegawai) TableName() string {
	return "pegawai"
}
