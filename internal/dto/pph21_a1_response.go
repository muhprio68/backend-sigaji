package dto

type PPh21A1Response struct {
	ID         int    `json:"id"`
	NIP        string `json:"nip"`
	NIK        string `json:"nik"`
	Nama       string `json:"nama"`
	Jabatan    string `json:"jabatan"`
	StatusPTKP string `json:"status_ptkp"`
	BulanAwal  int    `json:"bulan_awal"`
	BulanAkhir int    `json:"bulan_akhir"`

	// Rincian Penghasilan
	SalaryPensionJhtTht float64 `json:"salary_pension_jht_tht"`
	IncomeTaxBenefit    float64 `json:"income_tax_benefit"`
	InsurancePaidByEmp  float64 `json:"insurance_paid_by_emp"`
	OtherBenefit        float64 `json:"other_benefit"`
	TantiemBonusThr     float64 `json:"tantiem_bonus_thr"`
	TotalBruto          float64 `json:"total_bruto"`

	// Pengurangan
	BiayaJabatan        float64 `json:"biaya_jabatan"`
	PensionContribution float64 `json:"pension_contribution"`
	TotalNeto           float64 `json:"total_neto"`

	// Penghitungan PPh 21
	NilaiPTKP    float64 `json:"nilai_ptkp"` // Tambahan baru
	PKP          float64 `json:"pkp"`        // Tambahan baru
	PPh21Setahun float64 `json:"pph21_setahun"`
}
