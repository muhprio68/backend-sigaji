package service

import (
	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/model"
	"backend-sigaji/internal/repository"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type PegawaiService struct {
	db   *gorm.DB
	repo *repository.PegawaiRepository
}

func NewPegawaiService(db *gorm.DB, r *repository.PegawaiRepository) *PegawaiService {
	return &PegawaiService{
		db:   db,
		repo: r}
}

// Pastikan interface ini ditambahkan:
// CreatePegawai(req dto.CreatePegawaiRequest) error

func (s *PegawaiService) CreatePegawai(req dto.CreatePegawaiRequest) error {
	// Mapping dari DTO ke Model
	pegawai := model.Pegawai{
		NIP:        req.NIP,
		NIK:        req.NIK,
		Nama:       req.Nama,
		Jabatan:    req.Jabatan,
		StatusAsn:  req.StatusAsn,
		Golongan:   req.Golongan,
		StatusPTKP: req.StatusPTKP,
		Alamat:     req.Alamat,
	}

	// Panggil repository untuk simpan ke database
	err := s.repo.Create(&pegawai)
	if err != nil {
		return err
	}

	return nil
}

func (s *PegawaiService) GetAllPegawai(filter model.PegawaiFilter) ([]model.Pegawai, error) {
	return s.repo.GetAll(filter)
}

func (s *PegawaiService) GetByID(id uint) (*model.Pegawai, error) {
	return s.repo.FindByID(id)
}

func (s *PegawaiService) UpdatePegawai(data *model.Pegawai) error {
	return s.repo.Update(data)
}

func (s *PegawaiService) DeletePegawai(id uint) error {
	return s.repo.Delete(id)
}

func (s *PegawaiService) ImportExcel(file multipart.File) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {

		for i, row := range rows {
			if i == 0 {
				continue // skip header
			}

			if len(row) < 6 {
				continue
			}

			p := "TK"

			if row[13] == "1" {
				p = "K"
			}
			statusAsnInt, err := strconv.Atoi(strings.TrimSpace(row[8]))
			if err != nil {
				continue
			}

			if statusAsnInt < 0 || statusAsnInt > 255 {
				continue
			}

			nik := row[2]

			if nik == "" {
				nik = row[3]
			}

			pegawai := model.Pegawai{
				NIP:        row[0],
				NIK:        nik,
				Nama:       row[1],
				Jabatan:    row[6],
				Golongan:   row[9],
				StatusAsn:  uint8(statusAsnInt),
				StatusPTKP: p + "/" + row[14],
				Alamat:     row[11],
			}

			// cek duplikat NIP
			exist, _ := s.repo.FindByNIP(tx, pegawai.NIP)
			if exist != nil {
				continue // skip kalau sudah ada
			}

			if err := tx.Create(&pegawai).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
