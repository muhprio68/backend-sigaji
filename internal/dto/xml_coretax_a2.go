package dto

import "encoding/xml"

// A2Bulk adalah root element dari XML
type A2Bulk struct {
	XMLName  xml.Name `xml:"A2Bulk"`
	XmlnsXsi string   `xml:"xmlns:xsi,attr"`
	TIN      string   `xml:"TIN"` // NPWP Instansi
	ListOfA2 ListOfA2 `xml:"ListOfA2"`
}

type ListOfA2 struct {
	A2 []A2Item `xml:"A2"`
}

type A2Item struct {
	WorkForSecondEmployer        string        `xml:"WorkForSecondEmployer"`
	TaxPeriodMonthStart          int           `xml:"TaxPeriodMonthStart"`
	TaxPeriodMonthEnd            int           `xml:"TaxPeriodMonthEnd"`
	TaxPeriodYear                int           `xml:"TaxPeriodYear"`
	CounterpartTin               string        `xml:"CounterpartTin"` // NIK / NPWP Pegawai
	CounterpartNipNrp            string        `xml:"CounterpartNipNrp"`
	TaxExemptOpt                 string        `xml:"TaxExemptOpt"`        // Status PTKP
	StatusOfWithholding          string        `xml:"StatusOfWithholding"` // FullYear / PartialYear
	CounterpartPosition          string        `xml:"CounterpartPosition"`
	CounterpartRank              string        `xml:"CounterpartRank"`
	TaxObjectCode                string        `xml:"TaxObjectCode"` // 21-100-01
	NumberOfMonths               int           `xml:"NumberOfMonths"`
	SalaryPensionJhtTht          float64       `xml:"SalaryPensionJhtTht"`         // Gaji Pokok
	WifeBenefit                  float64       `xml:"WifeBenefit"`                 // Tunjangan Istri
	ChildBenefit                 float64       `xml:"ChildBenefit"`                // Tunjangan Anak
	IncomeImprovementBenefit     float64       `xml:"IncomeImprovementBenefit"`    // TPP
	StructuralFunctionalBenefit  float64       `xml:"StructuralFunctionalBenefit"` // Tunj Jabatan
	RiceBenefit                  float64       `xml:"RiceBenefit"`                 // Tunj Beras
	OtherBenefit                 float64       `xml:"OtherBenefit"`                // Tunj Lainnya
	OtherRegularIncome           float64       `xml:"OtherRegularIncome"`
	PensionContributionJhtThtFee float64       `xml:"PensionContributionJhtThtFee"` // IWP
	Zakat                        float64       `xml:"Zakat"`
	PrevWhTaxSlip                PrevWhTaxSlip `xml:"PrevWhTaxSlip"`
	Article21IncomeTax           float64       `xml:"Article21IncomeTax"` // PPh 21 Setahun
	IDPlaceOfBusinessActivity    string        `xml:"IDPlaceOfBusinessActivity"`
	WithholdingDate              string        `xml:"WithholdingDate"`
}

// Struct khusus untuk menangani <PrevWhTaxSlip xsi:nil="true"/>
type PrevWhTaxSlip struct {
	XsiNil string `xml:"xsi:nil,attr"`
}
