package repository

import (
	"backend-sigaji/internal/model"

	"gorm.io/gorm"
)

type GajiPokokRepository interface {
	Create(gaji *model.GajiPokok) error
	Delete(id uint) error
	FindByID(id uint) (*model.GajiPokok, error)
	FindAll() ([]model.GajiPokok, error)
	FindByPeriode(idPegawai uint, bulan int, tahun int) (*model.GajiPokok, error)
	GetByPeriode(jenis int, bulan int, tahun int, tipe string) ([]model.GajiPokok, error)
}

type gajiPokokRepository struct {
	db *gorm.DB
}

func NewGajiPokokRepository(db *gorm.DB) GajiPokokRepository {
	return &gajiPokokRepository{db}
}

func (r *gajiPokokRepository) Create(gaji *model.GajiPokok) error {
	return r.db.Create(gaji).Error
}

func (r *gajiPokokRepository) Delete(id uint) error {
	return r.db.Delete(&model.GajiPokok{}, id).Error
}

func (r *gajiPokokRepository) FindByID(id uint) (*model.GajiPokok, error) {
	var gaji model.GajiPokok
	err := r.db.Preload("Pegawai").First(&gaji, id).Error
	return &gaji, err
}

func (r *gajiPokokRepository) FindAll() ([]model.GajiPokok, error) {
	var gaji []model.GajiPokok
	err := r.db.Preload("Pegawai").Find(&gaji).Error
	return gaji, err
}

func (r *gajiPokokRepository) GetByPeriode(jenis int, bulan int, tahun int, tipe string) ([]model.GajiPokok, error) {
	var gaji []model.GajiPokok

	query := r.db.
		Joins("JOIN pegawai ON pegawai.id = gaji_pokok.id_pegawai").
		Preload("Pegawai").
		Where("gaji_pokok.tahun = ?", tahun)

	// 🔥 filter bulan hanya untuk bulanan / semua
	if tipe == "" || tipe == "bulanan" {
		if bulan != 0 {
			query = query.Where("gaji_pokok.bulan = ?", bulan)
		}
	}

	// 🔥 filter tipe (ini tambahan penting)
	if tipe != "" {
		query = query.Where("gaji_pokok.jenis_penghasilan = ?", tipe)
	}

	// 🔥 filter jenis ASN (PNS / PPPK)
	if jenis == 1 {
		// PNS → include CPNS
		query = query.Where("pegawai.status_asn IN ?", []int{1, 3})
	} else if jenis != 0 {
		query = query.Where("pegawai.status_asn = ?", jenis)
	}

	err := query.Find(&gaji).Error

	return gaji, err
}

func (r *gajiPokokRepository) FindByPeriode(idPegawai uint, bulan int, tahun int) (*model.GajiPokok, error) {
	var gaji model.GajiPokok
	err := r.db.Where("id_pegawai = ? AND bulan = ? AND tahun = ?", idPegawai, bulan, tahun).
		First(&gaji).Error
	return &gaji, err
}
