package repository

import (
	"backend-sigaji/internal/model"

	"gorm.io/gorm"
)

type TPPRepository interface {
	Create(tpp *model.TPP) error
	Update(tpp *model.TPP) error
	Delete(id uint) error
	FindByID(id uint) (*model.TPP, error)
	FindByPeriode(idPegawai uint, bulan int, tahun int) (*model.TPP, error)
	FindAll() ([]model.TPP, error)
	GetByPeriode(jenis int, bulan int, tahun int, tipe string) ([]model.TPP, error)
}

type tppRepository struct {
	db *gorm.DB
}

func NewTPPRepository(db *gorm.DB) TPPRepository {
	return &tppRepository{db}
}

func (r *tppRepository) Create(tpp *model.TPP) error {
	return r.db.Create(tpp).Error
}

func (r *tppRepository) Update(tpp *model.TPP) error {
	return r.db.Save(tpp).Error
}

func (r *tppRepository) Delete(id uint) error {
	return r.db.Delete(&model.TPP{}, id).Error
}

func (r *tppRepository) FindByID(id uint) (*model.TPP, error) {
	var tpp model.TPP
	err := r.db.Preload("Pegawai").First(&tpp, id).Error
	return &tpp, err
}

func (r *tppRepository) FindByPeriode(idPegawai uint, bulan int, tahun int) (*model.TPP, error) {
	var tpp model.TPP
	err := r.db.Where("id_pegawai = ? AND bulan = ? AND tahun = ?", idPegawai, bulan, tahun).
		First(&tpp).Error
	return &tpp, err
}

func (r *tppRepository) FindAll() ([]model.TPP, error) {
	var tpps []model.TPP
	err := r.db.Preload("Pegawai").Find(&tpps).Error
	return tpps, err
}
func (r *tppRepository) GetByPeriode(jenis int, bulan int, tahun int, tipe string) ([]model.TPP, error) {
	var tpp []model.TPP

	query := r.db.
		Joins("JOIN pegawai ON pegawai.id = tpp.id_pegawai").
		Preload("Pegawai").
		Where("tpp.tahun = ?", tahun) // Hanya filter tahun di awal

	// 🔥 filter bulan hanya untuk bulanan / semua
	if tipe == "" || tipe == "bulanan" {
		if bulan != 0 {
			query = query.Where("tpp.bulan = ?", bulan)
		}
	}

	// 🔥 filter tipe (ini tambahan penting)
	if tipe != "" {
		// PERHATIAN: Sesuaikan "jenis_penghasilan" dengan nama kolom di tabel tpp milikmu
		query = query.Where("tpp.jenis_penghasilan = ?", tipe)
	}

	// 🔥 filter jenis ASN (PNS / PPPK)
	if jenis == 1 {
		// Jika PNS → tampilkan PNS + CPNS
		query = query.Where("pegawai.status_asn IN ?", []int{1, 3})
	} else if jenis != 0 {
		// Selain itu tetap normal
		query = query.Where("pegawai.status_asn = ?", jenis)
	}

	err := query.Find(&tpp).Error

	return tpp, err
}
