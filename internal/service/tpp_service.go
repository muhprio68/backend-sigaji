package service

import (
	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/model"
	"backend-sigaji/internal/repository"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TPPService interface {
	CreateTPP(tpp *model.TPP) error
	UpdateTPP(tpp *model.TPP) error
	DeleteTPP(id uint) error
	GetTPPByID(id uint) (*dto.DetailTppResponse, error)
	GetAllTPP() ([]model.TPP, error)
	HitungTotalTPP(tpp *model.TPP) float64
	ImportExcel(file multipart.File, bulan int, tahun int, tipe string) error
	GetTppByPeriode(jenis int, bulan int, tahun int, tipe string) ([]dto.TppResponse, error)
}

type tppService struct {
	db   *gorm.DB
	repo repository.TPPRepository
}

func NewTPPService(db *gorm.DB, repo repository.TPPRepository) TPPService {
	return &tppService{
		db:   db,
		repo: repo}
}

func (s *tppService) CreateTPP(tpp *model.TPP) error {
	return s.repo.Create(tpp)
}

func (s *tppService) UpdateTPP(tpp *model.TPP) error {
	return s.repo.Update(tpp)
}

func (s *tppService) DeleteTPP(id uint) error {
	return s.repo.Delete(id)
}

func (s *tppService) GetTPPByID(id uint) (*dto.DetailTppResponse, error) {

	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, errors.New("data tidak ditemukan")
	}

	// Kalkulasi Bruto (Beban Kerja + Prestasi + Kondisi + Tunjangan)
	bruto := entity.TPPBeban + entity.TPPPrestasi + entity.TPPKondisi + entity.BPJS4 + entity.TunjPajak

	// Kalkulasi Total Potongan
	totalPotongan := entity.PotPajak + entity.BPJS4 + entity.BPJS1

	// Kalkulasi TPP Bersih
	bersih := bruto - totalPotongan

	return &dto.DetailTppResponse{
		ID:               entity.ID,
		Nama:             entity.Pegawai.Nama,
		NIP:              entity.Pegawai.NIP,
		Bulan:            entity.Bulan,
		Tahun:            entity.Tahun,
		JenisPenghasilan: entity.JenisPenghasilan,

		BebanKerja:    entity.TPPBeban,
		PrestasiKerja: entity.TPPPrestasi,
		KondisiKerja:  entity.TPPKondisi,
		TunjBpjs4:     entity.BPJS4,
		TunjPph21:     entity.TunjPajak,
		BrutoTpp:      bruto,
		PotPph21:      entity.PotPajak,
		PotBpjs4:      entity.BPJS4,
		PotBpjs1:      entity.BPJS1,
		TppBersih:     bersih,
	}, nil
}

func (s *tppService) GetAllTPP() ([]model.TPP, error) {
	return s.repo.FindAll()
}

func (s *tppService) GetTppByPeriode(jenis int, bulan int, tahun int, tipe string) ([]dto.TppResponse, error) {

	data, err := s.repo.GetByPeriode(jenis, bulan, tahun, tipe)
	if err != nil {
		return nil, err
	}

	var result []dto.TppResponse

	for _, item := range data {

		totalBruto := item.TPPBeban + item.TPPPrestasi + item.TPPKondisi + item.TunjPajak + item.BPJS4
		totalPotongan := item.PotPajak + item.BPJS4 + item.BPJS1
		totalBersih := totalBruto - totalPotongan

		result = append(result, dto.TppResponse{
			ID:          item.ID,
			Nama:        item.Pegawai.Nama,
			NIP:         item.Pegawai.NIP,
			BrutoTPP:    totalBruto,
			PotonganTPP: totalPotongan,
			TPPBersih:   totalBersih,
		})
	}

	return result, nil
}

// 🔥 Business Logic
func (s *tppService) HitungTotalTPP(tpp *model.TPP) float64 {
	return tpp.TPPBeban +
		tpp.TPPPrestasi +
		tpp.TPPKondisi +
		tpp.TunjPajak +
		tpp.BPJS4 -
		tpp.PotPajak -
		tpp.BPJS4 -
		tpp.BPJS1
}

func (s *tppService) ImportExcel(file multipart.File, bulan int, tahun int, tipe string) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Form_A") // Pastikan nama sheet-nya benar
	if err != nil {
		return err
	}

	// 1. Ambil SEMUA pegawai beserta STATUS-nya (PNS/PPPK/CPNS)
	var allPegawai []model.Pegawai
	if err := s.db.Select("id", "nip", "status_asn").Find(&allPegawai).Error; err != nil {
		return err
	}

	type PegawaiInfo struct {
		ID     uint
		Status int
	}

	pegawaiMap := make(map[string]PegawaiInfo)
	for _, p := range allPegawai {
		pegawaiMap[strings.TrimSpace(p.NIP)] = PegawaiInfo{ID: p.ID, Status: int(p.StatusAsn)}
	}

	var listTPP []model.TPP

	// Flag deteksi file Excel yang di-upload
	isImportingPNS := false
	isImportingPPPK := false

	// 🔥 KERANJANG REALISASI PNS 🔥
	var sumBebanPNS, sumPrestasiPNS, sumKondisiPNS, sumTunjPajakPNS, sumBpjs4PNS float64

	// 🔥 KERANJANG REALISASI PPPK 🔥
	var sumBebanPPPK, sumPrestasiPPPK, sumKondisiPPPK, sumBpjs4PPPK float64

	// Fungsi helper parsing angka
	parse := func(val string) float64 {
		val = strings.ReplaceAll(val, ",", "")
		val = strings.ReplaceAll(val, ".", "")
		val = strings.TrimSpace(val)
		if val == "" {
			return 0
		}
		valFloat, _ := strconv.ParseFloat(val, 64)
		return valFloat
	}

	// 2. Validasi Seluruh Baris Excel & Deteksi Golongan
	for i, row := range rows {
		if i < 5 {
			continue // Skip header (sesuaikan dengan format Excel TPP-mu)
		}
		if len(row) < 10 {
			continue
		}

		nip := strings.TrimSpace(row[0]) // Asumsi NIP di kolom A
		if nip == "" {
			continue
		}

		pInfo, exists := pegawaiMap[nip]
		if !exists {
			return fmt.Errorf("NIP '%s' pada Baris %d tidak ditemukan di Master Pegawai. Import dibatalkan", nip, i+1)
		}

		// PARSING ANGKA TPP
		bebanVal := parse(row[5])
		prestasiVal := parse(row[6])
		kondisiVal := parse(row[7])
		pajakVal := parse(row[8])
		bpjs4Val := parse(row[9])
		potPajakVal := parse(row[11])

		// Handle kolom ke-14 (index 13) dengan aman
		bpjs1Val := float64(0)
		if len(row) > 13 {
			bpjs1Val = parse(row[13])
		}

		// DETEKSI GOLONGAN & MASUKKAN KE KERANJANG REALISASI
		switch pInfo.Status {
		case 1, 3: // PNS & CPNS
			isImportingPNS = true
			sumBebanPNS += bebanVal
			sumPrestasiPNS += prestasiVal
			sumKondisiPNS += kondisiVal
			sumTunjPajakPNS += pajakVal // Tunjangan PPh PNS
			sumBpjs4PNS += bpjs4Val     // Iuran BPJS 4% Pemda PNS
		case 2: // PPPK
			isImportingPPPK = true
			sumBebanPPPK += bebanVal
			sumPrestasiPPPK += prestasiVal
			sumKondisiPPPK += kondisiVal
			sumBpjs4PPPK += bpjs4Val // Iuran BPJS 4% Pemda PPPK
		}

		// Tambahkan ke slice model TPP
		listTPP = append(listTPP, model.TPP{
			IDPegawai:        pInfo.ID,
			Bulan:            bulan,
			Tahun:            tahun,
			JenisPenghasilan: tipe,
			TPPBeban:         bebanVal,
			TPPPrestasi:      prestasiVal,
			TPPKondisi:       kondisiVal,
			TunjPajak:        pajakVal,
			PotPajak:         potPajakVal,
			BPJS4:            bpjs4Val,
			BPJS1:            bpjs1Val,
		})
	}

	// 3. Eksekusi Database (Hapus yang relevan, lalu Insert)
	return s.db.Transaction(func(tx *gorm.DB) error {

		// 🔥 PENGHAPUSAN AMAN
		if isImportingPNS {
			// Hapus Rincian TPP PNS
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND id_pegawai IN (SELECT id FROM pegawai WHERE status_asn IN (1, 3))", bulan, tahun, tipe).Delete(&model.TPP{}).Error; err != nil {
				return err
			}
			// Hapus Realisasi PNS (ID: 12, 15, 23, 25, 27)
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND sumber = ? AND id_akun_belanja IN (12, 15, 23, 25, 27)", bulan, tahun, tipe, "tpp").Delete(&model.RealisasiBelanja{}).Error; err != nil {
				return err
			}
		}

		if isImportingPPPK {
			// Hapus Rincian TPP PPPK
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND id_pegawai IN (SELECT id FROM pegawai WHERE status_asn = 2)", bulan, tahun, tipe).Delete(&model.TPP{}).Error; err != nil {
				return err
			}
			// Hapus Realisasi PPPK (ID: 16, 24, 26, 28)
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND sumber = ? AND id_akun_belanja IN (16, 24, 26, 28)", bulan, tahun, tipe, "tpp").Delete(&model.RealisasiBelanja{}).Error; err != nil {
				return err
			}
		}

		// Insert TPP per Pegawai
		for _, tpp := range listTPP {
			if err := tx.Create(&tpp).Error; err != nil {
				return err
			}
		}

		// 🔥 INSERT DATA KE TABEL REALISASI BELANJA 🔥
		realisasiData := []model.RealisasiBelanja{
			// TPP BEBAN KERJA
			{IDAkunBelanja: 23, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumBebanPNS},
			{IDAkunBelanja: 24, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumBebanPPPK},

			// TPP KONDISI KERJA
			{IDAkunBelanja: 25, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumKondisiPNS},
			{IDAkunBelanja: 26, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumKondisiPPPK}, // (Asumsi ID 26)

			// TPP PRESTASI KERJA
			{IDAkunBelanja: 27, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumPrestasiPNS},  // (Asumsi ID 27)
			{IDAkunBelanja: 28, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumPrestasiPPPK}, // (Asumsi ID 28)

			// BPJS KESEHATAN 4% (Masuk ke Iuran Pemda)
			{IDAkunBelanja: 15, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumBpjs4PNS},
			{IDAkunBelanja: 16, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumBpjs4PPPK},

			// TUNJANGAN PPH 21 (Hanya PNS)
			{IDAkunBelanja: 12, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "tpp", Realisasi: sumTunjPajakPNS},
		}

		// Insert hanya yang nominalnya ada / lebih dari 0
		for _, r := range realisasiData {
			if r.Realisasi > 0 {
				if err := tx.Create(&r).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
