package service

import (
	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/repository"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"strings"

	"gorm.io/gorm"
)

type PPh21Service struct {
	db   *gorm.DB
	Repo *repository.PPh21Repository
}

func NewPPh21Service(db *gorm.DB, repo *repository.PPh21Repository) *PPh21Service {
	return &PPh21Service{
		db:   db,
		Repo: repo}
}

// Fungsi helper Kategori (tetap sama)
func getKategoriTER(statusPTKP string) string {
	status := strings.ToUpper(strings.TrimSpace(statusPTKP))
	switch status {
	case "TK/0", "TK/1", "K/0":
		return "A"
	case "TK/2", "TK/3", "K/1", "K/2":
		return "B"
	case "K/3":
		return "C"
	default:
		return "A"
	}
}

// Fungsi utama Service
func (s *PPh21Service) GetPPh21Bulanan(jenis int, bulan int, tahun int, tipe string) (dto.PPh21PageResponse, error) {
	var finalResponse dto.PPh21PageResponse

	// 1. Ambil data mentah dari repo (Sudah dalam bentuk wrapper PageResponse)
	dataRepo, err := s.Repo.GetPenghasilanBulanan(jenis, bulan, tahun, tipe)
	if err != nil {
		return finalResponse, err
	}

	// 2. Ambil seluruh master data TER
	listTER, err := s.Repo.GetAllTarifTER()
	if err != nil {
		return finalResponse, err
	}

	var listPegawaiLengkap []dto.PPh21BulananResponse
	var totalBrutoDinas float64
	var totalPPh21Dinas float64

	// 3. Looping untuk hitung TER per orang
	// Kita iterasi dari dataRepo.ListPegawai
	for _, item := range dataRepo.ListPegawai {
		totalBruto := item.BrutoGaji + item.BrutoTPP
		kategoriTER := getKategoriTER(item.StatusPTKP)

		var tarifAktif float64 = 0.0
		for _, ter := range listTER {
			if strings.ToUpper(strings.TrimSpace(ter.Kategori)) == kategoriTER &&
				totalBruto >= ter.BatasBawah && totalBruto <= ter.BatasAtas {
				tarifAktif = ter.TarifPersentase
				break
			}
		}

		pph := totalBruto * tarifAktif
		pphBulat := math.Round(pph)

		// Simpan hasil hitungan individu
		pegawai := dto.PPh21BulananResponse{
			ID:         item.ID,
			NIP:        item.NIP,
			NIK:        item.NIK,
			Nama:       item.Nama,
			Jabatan:    item.Jabatan,
			StatusPTKP: item.StatusPTKP,
			BulanPajak: item.BulanPajak,
			BrutoGaji:  item.BrutoGaji,
			BrutoTPP:   item.BrutoTPP,
			TotalBruto: totalBruto,
			TarifTER:   tarifAktif,
			PPh21:      pphBulat,
		}

		listPegawaiLengkap = append(listPegawaiLengkap, pegawai)

		// Update akumulasi untuk Grand Total di Service
		totalBrutoDinas += totalBruto
		totalPPh21Dinas += pphBulat
	}

	// 4. Susun Final Response
	finalResponse.ListPegawai = listPegawaiLengkap
	finalResponse.Summary.TotalPegawai = len(listPegawaiLengkap)
	finalResponse.Summary.GrandTotalBruto = totalBrutoDinas
	finalResponse.Summary.GrandTotalPPh21 = totalPPh21Dinas

	return finalResponse, nil
}

func (s *PPh21Service) GenerateXMLEBupot(jenis int, bulan int, tahun int, tglPotong string, tipe string) ([]byte, error) {
	// 1. Ambil data
	dataPage, err := s.GetPPh21Bulanan(jenis, bulan, tahun, tipe)
	if err != nil {
		return nil, err
	}

	// 2. Siapkan wadah XML Root
	bulkXML := dto.MmPayrollBulk{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		TIN:      "0003326535601000",
	}

	// 3. Looping data pegawai
	for _, item := range dataPage.ListPegawai {

		ratePersen := item.TarifTER * 100
		rateStr := fmt.Sprintf("%.2f", ratePersen)
		rateStr = strings.TrimRight(rateStr, "0")
		rateStr = strings.TrimRight(rateStr, ".")

		// 🔥 Gunakan Bulan asli dari database, fallback ke parameter kalau 0
		masaPajak := item.BulanPajak
		if masaPajak == 0 {
			masaPajak = bulan
		}

		bupot := dto.MmPayroll{
			TaxPeriodMonth:            masaPajak, // Masuk ke XML
			TaxPeriodYear:             tahun,
			CounterpartOpt:            "Resident",
			CounterpartPassport:       &dto.NilElement{Nil: "true"},
			CounterpartTin:            item.NIK,
			StatusTaxExemption:        item.StatusPTKP,
			Position:                  item.Jabatan,
			TaxCertificate:            "N/A",
			TaxObjectCode:             "21-100-01",
			Gross:                     int64(item.TotalBruto),
			Rate:                      rateStr,
			IDPlaceOfBusinessActivity: "0003326535601000000000",
			WithholdingDate:           tglPotong,
		}

		bulkXML.ListOfMmPayroll.MmPayroll = append(bulkXML.ListOfMmPayroll.MmPayroll, bupot)
	}

	// 4. Generate XML
	xmlBytes, err := xml.MarshalIndent(bulkXML, "", "\t")
	if err != nil {
		return nil, err
	}

	customHeader := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
	finalXML := append(customHeader, xmlBytes...)

	return finalXML, nil
}

// --- SERVICE ---
func (s *PPh21Service) GetRekapA1Service(tahun int, pegawaiID int) ([]dto.PPh21A1Response, error) {
	// 1. Tarik rekap mentah dari Repo dengan parameter ID
	listData, err := s.Repo.GetRekapA1(tahun, pegawaiID)
	if err != nil {
		return nil, err
	}

	for i := range listData {
		p := &listData[i]

		// Hitung Biaya Jabatan (5% dari Bruto, max 6jt/tahun)
		biayaJab := p.TotalBruto * 0.05
		if biayaJab > 6000000 {
			biayaJab = 6000000
		}
		p.BiayaJabatan = math.Floor(biayaJab)

		// Hitung Neto Setahun
		p.TotalNeto = p.TotalBruto - p.BiayaJabatan - p.PensionContribution
		p.NilaiPTKP = getNilaiPTKP(p.StatusPTKP)

		// Hitung PKP (Bulatkan ke bawah ke ribuan terdekat)
		pkp := p.TotalNeto - p.NilaiPTKP
		if pkp < 0 {
			pkp = 0
		}
		p.PKP = math.Floor(pkp/1000) * 1000

		// Hitung PPh 21 Setahun
		p.PPh21Setahun = hitungPajakProgresif(p.PKP)
	}

	return listData, nil
}

func (s *PPh21Service) GenerateXMLA1(tahun int, pegawaiID int, tglPotong string) ([]byte, error) {
	// 1. Ambil data yang sudah dikalkulasi pajaknya
	data, err := s.GetRekapA1Service(tahun, pegawaiID)
	if err != nil {
		return nil, err
	}

	bulk := dto.A1Bulk{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		TIN:      os.Getenv("NPWP_DINAS"),
		ListOfA1: dto.ListOfA1{A1: []dto.A1Detail{}},
	}

	for _, item := range data {
		// Logika Mutasi: Jika bulan akhir < 12, maka PartialYear
		statusWithholding := "Annualized"
		if item.BulanAkhir < 12 {
			statusWithholding = "PartialYear"
		}

		detail := dto.A1Detail{
			WorkForSecondEmployer: "No",
			TaxPeriodMonthStart:   item.BulanAwal,
			TaxPeriodMonthEnd:     item.BulanAkhir,
			TaxPeriodYear:         tahun,
			CounterpartOpt:        "Resident",
			CounterpartPassport:   &dto.NilAttribute{Nil: true},
			CounterpartTin:        item.NIK,
			TaxExemptOpt:          item.StatusPTKP,
			StatusOfWithholding:   statusWithholding,
			CounterpartPosition:   item.Jabatan,
			TaxObjectCode:         "21-100-01",

			SalaryPensionJhtTht:          int64(item.SalaryPensionJhtTht),
			IncomeTaxBenefit:             int64(item.IncomeTaxBenefit),
			OtherBenefit:                 int64(item.OtherBenefit),
			InsurancePaidByEmployer:      int64(item.InsurancePaidByEmp),
			TantiemBonusThr:              int64(item.TantiemBonusThr),
			PensionContributionJhtThtFee: int64(item.PensionContribution),

			Article21IncomeTax:        int64(item.PPh21Setahun), // 🔥 Pajak hasil hitungan progresif
			PrevWhTaxSlip:             &dto.NilAttribute{Nil: true},
			TaxCertificate:            "N/A",
			IDPlaceOfBusinessActivity: os.Getenv("NITKU_DINAS"),
			WithholdingDate:           tglPotong, // 🔥 Tanggal dinamis dari Controller
		}
		bulk.ListOfA1.A1 = append(bulk.ListOfA1.A1, detail)
	}

	output, err := xml.MarshalIndent(bulk, "", "  ")
	if err != nil {
		return nil, err
	}

	header := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	return append([]byte(header), output...), nil
}

// --- SERVICE A2 (PNS) ---
func (s *PPh21Service) GetRekapA2Service(tahun int, pegawaiID int) ([]dto.PPh21TahunanResponse, error) {
	listData, err := s.Repo.GetRekapA2(tahun, pegawaiID)
	if err != nil {
		return nil, err
	}

	for i := range listData {
		p := &listData[i]

		// Hitung Bruto Full (Gaji + TPP)
		p.TotalBruto = p.TotalGajiPokok + p.TotalTunjIstri + p.TotalTunjAnak +
			p.TotalTunjBeras + p.TotalTunjJabatan + p.TotalTunjLain + p.TotalTPP + p.TotalBPJS4TPP

		biayaJab := p.TotalBruto * 0.05
		if biayaJab > 6000000 {
			biayaJab = 6000000
		}
		p.BiayaJabatan = math.Floor(biayaJab)
		p.IuranPensiun = math.Floor(p.TotalIWP)
		p.TotalNeto = p.TotalBruto - p.BiayaJabatan - p.IuranPensiun
		p.NilaiPTKP = getNilaiPTKP(p.StatusPTKP)

		pkp := p.TotalNeto - p.NilaiPTKP
		if pkp < 0 {
			pkp = 0
		}
		p.PKP = math.Floor(pkp/1000) * 1000
		p.PPh21Setahun = hitungPajakProgresif(p.PKP)
	}

	return listData, nil
}

func (s *PPh21Service) GenerateXMLA2(data []dto.PPh21TahunanResponse, tahun int, tglPotong string) ([]byte, error) {
	var listA2 []dto.A2Item

	for _, p := range data {
		statusWithholding := "FullYear"
		if p.BulanAkhir < 12 {
			statusWithholding = "PartialYear"
		}

		a2 := dto.A2Item{
			WorkForSecondEmployer: "No",
			TaxPeriodMonthStart:   p.BulanAwal,
			TaxPeriodMonthEnd:     p.BulanAkhir,
			TaxPeriodYear:         tahun,
			CounterpartTin:        p.NIK,
			CounterpartNipNrp:     p.NIP,
			TaxExemptOpt:          p.StatusPTKP,
			StatusOfWithholding:   statusWithholding,
			CounterpartPosition:   p.Jabatan,
			CounterpartRank:       p.Golongan,
			TaxObjectCode:         "21-100-02", // 🔥 Kode Objek Pajak PNS

			SalaryPensionJhtTht:          p.TotalGajiPokok,
			WifeBenefit:                  p.TotalTunjIstri,
			ChildBenefit:                 p.TotalTunjAnak,
			IncomeImprovementBenefit:     p.TotalTPP,
			StructuralFunctionalBenefit:  p.TotalTunjJabatan,
			RiceBenefit:                  p.TotalTunjBeras,
			OtherBenefit:                 p.TotalTunjLain,
			PensionContributionJhtThtFee: p.TotalIWP,

			PrevWhTaxSlip:             dto.PrevWhTaxSlip{XsiNil: "true"},
			Article21IncomeTax:        p.PPh21Setahun,
			IDPlaceOfBusinessActivity: os.Getenv("NITKU_DINAS"),
			WithholdingDate:           tglPotong,
		}
		listA2 = append(listA2, a2)
	}

	xmlData := dto.A2Bulk{
		XmlnsXsi: "http://www.w3.org/2001/XMLSchema-instance",
		TIN:      os.Getenv("NPWP_DINAS"),
		ListOfA2: dto.ListOfA2{A2: listA2},
	}

	output, err := xml.MarshalIndent(xmlData, "", "    ")
	if err != nil {
		return nil, err
	}

	header := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"
	return []byte(header + string(output)), nil
}

// --- FUNGSI HELPER (Taruh di bawah atau di file utils) ---

func getNilaiPTKP(status string) float64 {
	status = strings.ToUpper(strings.TrimSpace(status))
	switch status {
	case "TK/0":
		return 54000000
	case "TK/1", "K/0":
		return 58500000
	case "TK/2", "K/1":
		return 63000000
	case "TK/3", "K/2":
		return 67500000
	case "K/3":
		return 72000000
	default:
		return 54000000 // Default paling aman
	}
}

func hitungPajakProgresif(pkp float64) float64 {
	var pajak float64 = 0

	if pkp <= 0 {
		return 0
	}

	// 1. Sampai 60 Juta (5%)
	lapis1 := math.Min(pkp, 60000000)
	pajak += lapis1 * 0.05
	pkp -= lapis1

	// 2. 60 Juta - 250 Juta (15%)
	if pkp > 0 {
		lapis2 := math.Min(pkp, 190000000) // 250m - 60m
		pajak += lapis2 * 0.15
		pkp -= lapis2
	}

	// 3. 250 Juta - 500 Juta (25%)
	if pkp > 0 {
		lapis3 := math.Min(pkp, 250000000) // 500m - 250m
		pajak += lapis3 * 0.25
		pkp -= lapis3
	}

	// 4. 500 Juta - 5 Miliar (30%)
	if pkp > 0 {
		lapis4 := math.Min(pkp, 4500000000) // 5M - 500m
		pajak += lapis4 * 0.30
		pkp -= lapis4
	}

	// 5. Di Atas 5 Miliar (35%)
	if pkp > 0 {
		pajak += pkp * 0.35
	}

	return math.Floor(pajak) // PPh 21 selalu dibulatkan ke bawah
}
