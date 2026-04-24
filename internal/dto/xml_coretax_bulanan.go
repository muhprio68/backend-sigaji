package dto

import "encoding/xml"

// Struct untuk menghasilkan xsi:nil="true"
type NilElement struct {
	Nil string `xml:"xsi:nil,attr"`
}

// Root XML (MmPayrollBulk)
type MmPayrollBulk struct {
	XMLName         xml.Name        `xml:"MmPayrollBulk"`
	XmlnsXsi        string          `xml:"xmlns:xsi,attr"` // Untuk namespace w3.org
	TIN             string          `xml:"TIN"`            // NPWP Instansi
	ListOfMmPayroll ListOfMmPayroll `xml:"ListOfMmPayroll"`
}

// Wrapper untuk list pegawai
type ListOfMmPayroll struct {
	MmPayroll []MmPayroll `xml:"MmPayroll"`
}

// Data per pegawai
type MmPayroll struct {
	TaxPeriodMonth            int         `xml:"TaxPeriodMonth"`
	TaxPeriodYear             int         `xml:"TaxPeriodYear"`
	CounterpartOpt            string      `xml:"CounterpartOpt"`
	CounterpartPassport       *NilElement `xml:"CounterpartPassport"`
	CounterpartTin            string      `xml:"CounterpartTin"` // NIK Pegawai / NPWP 16 Digit
	StatusTaxExemption        string      `xml:"StatusTaxExemption"`
	Position                  string      `xml:"Position"`
	TaxCertificate            string      `xml:"TaxCertificate"`
	TaxObjectCode             string      `xml:"TaxObjectCode"`
	Gross                     int64       `xml:"Gross"` // Dibulatkan jadi integer
	Rate                      string      `xml:"Rate"`  // Persentase TER (misal: 1 atau 1.25)
	IDPlaceOfBusinessActivity string      `xml:"IDPlaceOfBusinessActivity"`
	WithholdingDate           string      `xml:"WithholdingDate"` // Format: YYYY-MM-DD
}
