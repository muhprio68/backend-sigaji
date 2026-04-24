package repository

import (
	"backend-sigaji/internal/model"

	"gorm.io/gorm"
)

type PegawaiRepository struct {
	DB *gorm.DB
}

func NewPegawaiRepository(db *gorm.DB) *PegawaiRepository {
	return &PegawaiRepository{DB: db}
}

// Pastikan kamu menambahkan interface ini jika menggunakan interface
// Create(pegawai *model.Pegawai) error

func (r *PegawaiRepository) Create(pegawai *model.Pegawai) error {
	// GORM otomatis menjalankan query INSERT INTO pegawai ...
	err := r.DB.Create(pegawai).Error
	return err
}

func (r *PegawaiRepository) GetAll(filter model.PegawaiFilter) ([]model.Pegawai, error) {

	var pegawai []model.Pegawai
	query := r.DB.Model(&model.Pegawai{})

	if filter.StatusAsn != nil {

		if *filter.StatusAsn == model.StatusPNS {
			query = query.Where("status_asn IN (?)", []uint8{
				model.StatusPNS,
				model.StatusCPNS,
			})
		} else {
			query = query.Where("status_asn = ?", *filter.StatusAsn)
		}
	}

	err := query.Find(&pegawai).Error
	return pegawai, err
}

func (r *PegawaiRepository) FindByID(id uint) (*model.Pegawai, error) {
	var pegawai model.Pegawai
	err := r.DB.Where("id = ?", id).First(&pegawai).Error
	return &pegawai, err
}

func (r *PegawaiRepository) FindByNIP(tx *gorm.DB, nip string) (*model.Pegawai, error) {
	var pegawai model.Pegawai
	err := tx.Where("nip = ?", nip).First(&pegawai).Error
	if err != nil {
		return nil, err
	}
	return &pegawai, nil
}

func (r *PegawaiRepository) Update(data *model.Pegawai) error {

	return r.DB.Model(&model.Pegawai{}).
		Where("id = ?", data.ID).
		Updates(map[string]interface{}{
			"nip":         data.NIP,
			"nik":         data.NIK,
			"nama":        data.Nama,
			"jabatan":     data.Jabatan,
			"status_asn":  data.StatusAsn,
			"golongan":    data.Golongan,
			"status_ptkp": data.StatusPTKP,
			"alamat":      data.Alamat,
		}).Error
}

func (r *PegawaiRepository) Delete(id uint) error {
	var pegawai model.Pegawai

	if err := r.DB.First(&pegawai, id).Error; err != nil {
		return err
	}

	return r.DB.Delete(&pegawai).Error
}
