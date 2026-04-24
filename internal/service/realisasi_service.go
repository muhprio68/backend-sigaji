package service

import (
	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/repository"
	"time"

	"gorm.io/gorm"
)

type RealisasiService struct {
	db   *gorm.DB
	repo *repository.RealisasiRepository
}

func NewRealisasiService(db *gorm.DB, repo *repository.RealisasiRepository) *RealisasiService {
	return &RealisasiService{
		db:   db,
		repo: repo}
}

func (s *RealisasiService) GetPenyerapanAnggaran(tahun int, bulan int) ([]dto.RealisasiResponse, error) {
	return s.repo.GetPenyerapanAnggaran(tahun, bulan)
}

// Tambahkan fungsi ini di realisasi_service.go
func (s *RealisasiService) GetCountPegawaiAktif(bulan int, tahun int) (int64, int64, error) {
	return s.repo.GetCountPegawaiAktif(bulan, tahun)
}

// Tambahkan di realisasi_service.go
func (s *RealisasiService) GetLastUpdate() (time.Time, error) {
	return s.repo.GetLastUpdate()
}
