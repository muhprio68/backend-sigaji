package dto

import "encoding/xml"

type A1Bulk struct {
	XMLName  xml.Name `xml:"A1Bulk"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	TIN      string   `xml:"TIN"` // NPWP Dinas
	ListOfA1 ListOfA1 `xml:"ListOfA1"`
}

type ListOfA1 struct {
	A1 []A1Detail `xml:"A1"`
}

type NilAttribute struct {
	Nil bool `xml:"xsi:nil,attr"`
}

type A1Detail struct {
	WorkForSecondEmployer        string        `xml:"WorkForSecondEmployer"`
	TaxPeriodMonthStart          int           `xml:"TaxPeriodMonthStart"`
	TaxPeriodMonthEnd            int           `xml:"TaxPeriodMonthEnd"`
	TaxPeriodYear                int           `xml:"TaxPeriodYear"`
	CounterpartOpt               string        `xml:"CounterpartOpt"`
	CounterpartPassport          *NilAttribute `xml:"CounterpartPassport"`
	CounterpartTin               string        `xml:"CounterpartTin"`
	TaxExemptOpt                 string        `xml:"TaxExemptOpt"`
	StatusOfWithholding          string        `xml:"StatusOfWithholding"`
	CounterpartPosition          string        `xml:"CounterpartPosition"`
	TaxObjectCode                string        `xml:"TaxObjectCode"`
	NumberOfMonths               int           `xml:"NumberOfMonths"`
	SalaryPensionJhtTht          int64         `xml:"SalaryPensionJhtTht"`
	GrossUpOpt                   string        `xml:"GrossUpOpt"`
	IncomeTaxBenefit             int64         `xml:"IncomeTaxBenefit"`
	OtherBenefit                 int64         `xml:"OtherBenefit"`
	Honorarium                   int64         `xml:"Honorarium"`
	InsurancePaidByEmployer      int64         `xml:"InsurancePaidByEmployer"`
	Natura                       int64         `xml:"Natura"`
	TantiemBonusThr              int64         `xml:"TantiemBonusThr"`
	PensionContributionJhtThtFee int64         `xml:"PensionContributionJhtThtFee"`
	Zakat                        int64         `xml:"Zakat"`
	PrevWhTaxSlip                *NilAttribute `xml:"PrevWhTaxSlip"`
	TaxCertificate               string        `xml:"TaxCertificate"`
	Article21IncomeTax           int64         `xml:"Article21IncomeTax"`
	IDPlaceOfBusinessActivity    string        `xml:"IDPlaceOfBusinessActivity"`
	WithholdingDate              string        `xml:"WithholdingDate"`
}

// Struct untuk mapping data dari Database
type DataBupotA1DB struct {
	NPWP      string
	NIK       string
	Nama      string
	Alamat    string
	Jabatan   string
	PTKP      string
	GajiPokok int64
	TPP       int64
	Thr       int64
	Asuransi  int64
	Potongan  int64
	PphTotal  int64
}
