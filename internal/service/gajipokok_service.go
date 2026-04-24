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

type GajiPokokService interface {
	Create(gaji *model.GajiPokok) error
	Delete(id uint) error
	GetByID(id uint) (*dto.DetailGajiResponse, error)
	GetAll() ([]model.GajiPokok, error)
	HitungTotal(gaji *model.GajiPokok) float64
	ImportExcel(file multipart.File, bulan int, tahun int, jenis_gapok string) error
	GetGajiByPeriode(jenis int, bulan int, tahun int, tipe string) ([]dto.GajiResponse, error)
}

type gajiPokokService struct {
	db   *gorm.DB
	repo repository.GajiPokokRepository
}

func NewGajiPokokService(db *gorm.DB, repo repository.GajiPokokRepository) GajiPokokService {
	return &gajiPokokService{
		db:   db,
		repo: repo,
	}
}

func (s *gajiPokokService) Create(gaji *model.GajiPokok) error {
	// 🔥 cek duplicate periode
	existing, _ := s.repo.FindByPeriode(gaji.IDPegawai, gaji.Bulan, gaji.Tahun)
	if existing != nil && existing.ID != 0 {
		return errors.New("gaji pokok periode ini sudah ada")
	}

	return s.repo.Create(gaji)
}

func (s *gajiPokokService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *gajiPokokService) GetByID(id uint) (*dto.DetailGajiResponse, error) {
	entity, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("data tidak ditemukan")
	}

	// 1. Bruto Pajak (Gaji + Tunjangan + Jaminan Pemda IN)
	brutoPajak := calculateBruto(entity)

	// 2. Hitung Total Potongan Pegawai (Termasuk Jaminan Pemda OUT)
	totalPotongan := entity.Iwp8Persen +
		entity.Iwp1Persen +
		entity.PotonganTapera +
		entity.PotonganLain +
		entity.TunjPph +
		entity.BpjsKesPemda + // IN-OUT
		entity.JkkPemda + // IN-OUT
		entity.JkmPemda // IN-OUT

	// 3. Gaji Bersih (Bruto - Semua Potongan)
	bersih := brutoPajak - totalPotongan

	return &dto.DetailGajiResponse{
		ID:               entity.ID,
		Nama:             entity.Pegawai.Nama,
		NIP:              entity.Pegawai.NIP,
		Bulan:            entity.Bulan,
		Tahun:            entity.Tahun,
		JenisPenghasilan: entity.JenisPenghasilan,

		GajiPokok:          entity.GajiPokok,
		TunjSuamiIstri:     entity.TunjSuamiIstri,
		TunjAnak:           entity.TunjAnak,
		TunjJabatan:        entity.TunjJabatan,
		TunjFungsional:     entity.TunjFungsional,
		TunjFungsionalUmum: entity.TunjFungsionalUmum,
		TunjBeras:          entity.TunjBeras,
		TunjPph:            entity.TunjPph,
		Pembulatan:         entity.Pembulatan,

		BpjsKesPemda: entity.BpjsKesPemda,
		JkkPemda:     entity.JkkPemda,
		JkmPemda:     entity.JkmPemda,

		Iwp8Persen:     entity.Iwp8Persen,
		Iwp1Persen:     entity.Iwp1Persen,
		PotonganTapera: entity.PotonganTapera,
		PotonganLain:   entity.PotonganLain,

		BrutoGaji:     brutoPajak,
		TotalPotongan: totalPotongan,
		GajiBersih:    bersih,
	}, nil
}

func (s *gajiPokokService) GetGajiByPeriode(jenis int, bulan int, tahun int, tipe string) ([]dto.GajiResponse, error) {
	data, err := s.repo.GetByPeriode(jenis, bulan, tahun, tipe)
	if err != nil {
		return nil, err
	}

	var result []dto.GajiResponse

	for _, item := range data {
		brutoGaji := calculateBruto(&item)

		// Sama seperti GetByID, Jaminan Pemda ikut dihitung sebagai Potongan
		totalPotongan := item.Iwp8Persen +
			item.Iwp1Persen +
			item.PotonganTapera +
			item.PotonganLain +
			item.TunjPph +
			item.BpjsKesPemda +
			item.JkkPemda +
			item.JkmPemda

		totalBersih := brutoGaji - totalPotongan

		result = append(result, dto.GajiResponse{
			ID:           item.ID,
			Nama:         item.Pegawai.Nama,
			NIP:          item.Pegawai.NIP,
			BrutoGaji:    brutoGaji,
			PotonganGaji: totalPotongan,
			GajiBersih:   totalBersih,
		})
	}
	return result, nil
}

func (s *gajiPokokService) GetAll() ([]model.GajiPokok, error) {
	return s.repo.FindAll()
}

func (s *gajiPokokService) ImportExcel(file multipart.File, bulan int, tahun int, tipe string) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	// 1. Ambil SEMUA pegawai beserta STATUS-nya (PNS/PPPK)
	var allPegawai []model.Pegawai
	// Asumsi nama kolom statusnya "status_asn" dan nilainya 1 = PNS, 2 = PPPK, 3 = CPNS
	if err := s.db.Select("id", "nip", "nama", "status_asn").Find(&allPegawai).Error; err != nil {
		return err
	}

	// Buat struct kecil untuk menyimpan ID dan Status di dalam Map
	type PegawaiInfo struct {
		ID     uint
		Status int
	}

	pegawaiMap := make(map[string]PegawaiInfo)
	for _, p := range allPegawai {
		pegawaiMap[strings.TrimSpace(p.NIP)] = PegawaiInfo{ID: p.ID, Status: int(p.StatusAsn)}
	}

	var listGaji []model.GajiPokok

	// 🔥 KERANJANG REALISASI PNS 🔥
	var sumGajiPNS, sumTunjKeluargaPNS, sumTunjJabatanPNS, sumTunjFungsionalPNS float64
	var sumTunjUmumPNS, sumTunjBerasPNS, sumTunjPphPNS, sumPembulatanPNS float64
	var sumBpjsPemdaPNS, sumJkkPemdaPNS, sumJkmPemdaPNS float64

	// 🔥 KERANJANG REALISASI PPPK 🔥
	var sumGajiPPPK, sumTunjKeluargaPPPK, sumTunjFungsionalPPPK float64
	var sumTunjUmumPPPK, sumTunjBerasPPPK, sumPembulatanPPPK float64
	var sumBpjsPemdaPPPK, sumJkkPemdaPPPK, sumJkmPemdaPPPK float64

	// 2. Validasi Seluruh Baris Excel
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 30 {
			continue
		}

		nip := strings.TrimSpace(row[0])
		pInfo, exists := pegawaiMap[nip]
		if !exists {
			return fmt.Errorf("NIP '%s' pada Baris %d tidak ditemukan di Master Pegawai. Import dibatalkan", nip, i+1)
		}

		// PARSING ANGKA
		gajiVal := parseFloat(row, 21)
		sistriVal := parseFloat(row, 22)
		anakVal := parseFloat(row, 23)
		jabatanVal := parseFloat(row, 25)
		fungsionalVal := parseFloat(row, 26)
		umumVal := parseFloat(row, 27)
		berasVal := parseFloat(row, 28)
		pphVal := parseFloat(row, 29)
		pembulatanVal := parseFloat(row, 30)

		bpjsPemdaVal := parseFloat(row, 31)
		jkkPemdaVal := parseFloat(row, 32)
		jkmPemdaVal := parseFloat(row, 33)

		iwp8Val := parseFloat(row, 37)
		iwp1Val := parseFloat(row, 38)
		taperaVal := parseFloat(row, 34)

		potonganLainVal := parseFloat(row, 40) + parseFloat(row, 41)

		if pInfo.Status == 1 || pInfo.Status == 3 {
			iwp8Val += parseFloat(row, 35)
		}

		if pInfo.Status == 2 {
			potonganLainVal += parseFloat(row, 35)
		}

		// Tambahkan ke slice model Gaji
		listGaji = append(listGaji, model.GajiPokok{
			IDPegawai:          pInfo.ID,
			Bulan:              bulan,
			Tahun:              tahun,
			JenisPenghasilan:   tipe,
			GajiPokok:          gajiVal,
			TunjSuamiIstri:     sistriVal,
			TunjAnak:           anakVal,
			TunjJabatan:        jabatanVal,
			TunjFungsional:     fungsionalVal,
			TunjFungsionalUmum: umumVal,
			TunjBeras:          berasVal,
			TunjPph:            pphVal,
			Pembulatan:         pembulatanVal,
			BpjsKesPemda:       bpjsPemdaVal,
			JkkPemda:           jkkPemdaVal,
			JkmPemda:           jkmPemdaVal,
			Iwp8Persen:         iwp8Val,
			Iwp1Persen:         iwp1Val,
			PotonganTapera:     taperaVal,
			PotonganLain:       potonganLainVal,
		})

		// 🔥 PISAHKAN JUMLAH KE KERANJANG PNS / PPPK 🔥
		switch pInfo.Status {
		case 1, 3: // 1 = PNS, 3 = CPNS
			sumGajiPNS += gajiVal
			sumTunjKeluargaPNS += (sistriVal + anakVal)
			sumTunjJabatanPNS += jabatanVal
			sumTunjFungsionalPNS += fungsionalVal
			sumTunjUmumPNS += umumVal
			sumTunjBerasPNS += berasVal
			sumTunjPphPNS += pphVal
			sumPembulatanPNS += pembulatanVal
			sumBpjsPemdaPNS += bpjsPemdaVal
			sumJkkPemdaPNS += jkkPemdaVal
			sumJkmPemdaPNS += jkmPemdaVal

		case 2: // 2 = PPPK
			sumGajiPPPK += gajiVal
			sumTunjKeluargaPPPK += (sistriVal + anakVal)
			sumTunjFungsionalPPPK += fungsionalVal
			sumTunjUmumPPPK += umumVal
			sumTunjBerasPPPK += berasVal
			sumPembulatanPPPK += pembulatanVal
			sumBpjsPemdaPPPK += bpjsPemdaVal
			sumJkkPemdaPPPK += jkkPemdaVal
			sumJkmPemdaPPPK += jkmPemdaVal
		}
	}

	// Cek file ini mayoritas isinya PNS atau PPPK (buat filter hapus data lama)
	isImportingPNS := sumGajiPNS > 0
	isImportingPPPK := sumGajiPPPK > 0

	// 3. Eksekusi Database
	return s.db.Transaction(func(tx *gorm.DB) error {

		// 🔥 PENGHAPUSAN AMAN: Hanya hapus yang sedang diimport (PNS saja atau PPPK saja)
		if isImportingPNS {
			// Hapus Gaji Pokok PNS saja
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND id_pegawai IN (SELECT id FROM pegawai WHERE status_asn IN (1, 3))", bulan, tahun, tipe).Delete(&model.GajiPokok{}).Error; err != nil {
				return err
			}
			// Hapus Realisasi PNS saja
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND sumber = ? AND id_akun_belanja IN (1,3,5,6,8,10,12,13,15,17,19,21)", bulan, tahun, tipe, "gaji_pokok").Delete(&model.RealisasiBelanja{}).Error; err != nil {
				return err
			}
		}

		if isImportingPPPK {
			// Hapus Gaji Pokok PPPK saja
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND id_pegawai IN (SELECT id FROM pegawai WHERE status_asn = 2)", bulan, tahun, tipe).Delete(&model.GajiPokok{}).Error; err != nil {
				return err
			}
			// Hapus Realisasi PPPK saja
			if err := tx.Where("bulan = ? AND tahun = ? AND jenis_penghasilan = ? AND sumber = ? AND id_akun_belanja IN (2,4,7,9,11,14,16,18,20,22)", bulan, tahun, tipe, "gaji_pokok").Delete(&model.RealisasiBelanja{}).Error; err != nil {
				return err
			}
		}

		// Insert Gaji per Pegawai
		for _, gaji := range listGaji {
			if err := tx.Create(&gaji).Error; err != nil {
				return err
			}
		}

		// 🔥 INSERT DATA KE TABEL REALISASI SESUAI ID 28 BARIS KAMU 🔥
		realisasiData := []model.RealisasiBelanja{
			// GAJI POKOK
			{IDAkunBelanja: 1, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumGajiPNS},
			{IDAkunBelanja: 2, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumGajiPPPK},

			// TUNJANGAN KELUARGA
			{IDAkunBelanja: 3, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjKeluargaPNS},
			{IDAkunBelanja: 4, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjKeluargaPPPK},

			// TUNJANGAN JABATAN (Hanya PNS)
			{IDAkunBelanja: 5, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjJabatanPNS},

			// TUNJANGAN FUNGSIONAL
			{IDAkunBelanja: 6, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjFungsionalPNS},
			{IDAkunBelanja: 7, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjFungsionalPPPK},

			// TUNJANGAN UMUM
			{IDAkunBelanja: 8, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjUmumPNS},
			{IDAkunBelanja: 9, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjUmumPPPK},

			// TUNJANGAN BERAS
			{IDAkunBelanja: 10, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjBerasPNS},
			{IDAkunBelanja: 11, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjBerasPPPK},

			// TUNJANGAN PPH (Hanya PNS)
			{IDAkunBelanja: 12, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumTunjPphPNS},

			// PEMBULATAN
			{IDAkunBelanja: 13, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumPembulatanPNS},
			{IDAkunBelanja: 14, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumPembulatanPPPK},

			// BPJS KESEHATAN (Pemda)
			{IDAkunBelanja: 15, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumBpjsPemdaPNS},
			{IDAkunBelanja: 16, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumBpjsPemdaPPPK},

			// JKK
			{IDAkunBelanja: 17, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumJkkPemdaPNS},
			{IDAkunBelanja: 18, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumJkkPemdaPPPK},

			// JKM
			{IDAkunBelanja: 19, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumJkmPemdaPNS},
			{IDAkunBelanja: 20, Bulan: bulan, Tahun: tahun, JenisPenghasilan: tipe, Sumber: "gaji_pokok", Realisasi: sumJkmPemdaPPPK},
		}

		for _, r := range realisasiData {
			if r.Realisasi > 0 { // Hanya simpan kalau nominalnya ada (tidak nol)
				if err := tx.Create(&r).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// --- FUNGSI HELPER ---

// parseFloat mencegah error panic kalau kolom di excel ternyata kosong ("")
func parseFloat(row []string, index int) float64 {
	if index >= len(row) {
		return 0
	}
	valStr := strings.TrimSpace(row[index])
	if valStr == "" || valStr == "-" {
		return 0
	}
	// Buang koma atau titik ribuan kalau ada di Excel
	valStr = strings.ReplaceAll(valStr, ",", "")
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0
	}
	return val
}

func calculateBruto(gaji *model.GajiPokok) float64 {
	penghasilan := gaji.GajiPokok +
		gaji.TunjSuamiIstri +
		gaji.TunjAnak +
		gaji.TunjJabatan +
		gaji.TunjFungsional +
		gaji.TunjFungsionalUmum +
		gaji.TunjBeras +
		gaji.TunjPph +
		gaji.Pembulatan

	jaminanPemda := gaji.BpjsKesPemda + gaji.JkkPemda + gaji.JkmPemda

	return penghasilan + jaminanPemda
}

func (s *gajiPokokService) HitungTotal(gaji *model.GajiPokok) float64 {
	return calculateBruto(gaji)
}
