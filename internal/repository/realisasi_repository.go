package repository

import (
	"backend-sigaji/internal/dto"
	"time"

	"gorm.io/gorm"
)

type RealisasiRepository struct {
	DB *gorm.DB
}

func NewRealisasiRepository(db *gorm.DB) *RealisasiRepository {
	return &RealisasiRepository{DB: db}
}

func (r *RealisasiRepository) GetPenyerapanAnggaran(tahun int, bulan int) ([]dto.RealisasiResponse, error) {
	var results []dto.RealisasiResponse

	// Base Query
	query := `
		SELECT 
			a.id, 
			a.kode_rekening, 
			a.uraian, 
			a.kelompok, 
			a.pagu_anggaran,
			COALESCE(SUM(r.realisasi), 0) AS total_realisasi,
			(a.pagu_anggaran - COALESCE(SUM(r.realisasi), 0)) AS sisa_anggaran,
			CASE 
				WHEN a.pagu_anggaran > 0 THEN (COALESCE(SUM(r.realisasi), 0) / a.pagu_anggaran) * 100
				ELSE 0 
			END AS persentase
		FROM akun_belanja a
		LEFT JOIN realisasi_belanja r ON a.id = r.id_akun_belanja AND r.tahun = ?
	`
	args := []interface{}{tahun}

	// Kalau bulan dikirim (>0), kita hitung Akumulatif (S/D Bulan tsb)
	if bulan > 0 {
		query += ` AND r.bulan <= ? `
		args = append(args, bulan)
	}

	// Grouping dan Sorting
	query += `
		GROUP BY a.id, a.kode_rekening, a.uraian, a.kelompok, a.pagu_anggaran
		ORDER BY a.kode_rekening ASC
	`

	// Eksekusi Scan langsung ke DTO
	err := r.DB.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Tambahkan fungsi ini di realisasi_repository.go
func (r *RealisasiRepository) GetCountPegawaiAktif(bulan int, tahun int) (int64, int64, error) {
	var countPNS, countPPPK int64

	// Hitung PNS & CPNS (Status 1 & 3)
	errPNS := r.DB.Table("gaji_pokok").
		Joins("JOIN pegawai ON pegawai.id = gaji_pokok.id_pegawai").
		Where("gaji_pokok.bulan = ? AND gaji_pokok.tahun = ? AND pegawai.status_asn IN (1, 3)", bulan, tahun).
		Distinct("gaji_pokok.id_pegawai").
		Count(&countPNS).Error

	if errPNS != nil {
		return 0, 0, errPNS
	}

	// Hitung PPPK (Status 2)
	errPPPK := r.DB.Table("gaji_pokok").
		Joins("JOIN pegawai ON pegawai.id = gaji_pokok.id_pegawai").
		Where("gaji_pokok.bulan = ? AND gaji_pokok.tahun = ? AND pegawai.status_asn = 2", bulan, tahun).
		Distinct("gaji_pokok.id_pegawai").
		Count(&countPPPK).Error

	if errPPPK != nil {
		return 0, 0, errPPPK
	}

	return countPNS, countPPPK, nil
}

// Tambahkan di realisasi_repository.go
func (r *RealisasiRepository) GetLastUpdate() (time.Time, error) {
	var lastUpdate time.Time
	// Pake GORM sekalian biar seragam
	err := r.DB.Table("realisasi_belanja").Select("MAX(updated_at)").Scan(&lastUpdate).Error
	return lastUpdate, err
}
