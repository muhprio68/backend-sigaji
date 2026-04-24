package dto

type DetailGajiResponse struct {
	ID               uint   `json:"id"`
	Nama             string `json:"nama"`
	NIP              string `json:"nip"`
	Bulan            int    `json:"bulan"`
	Tahun            int    `json:"tahun"`
	JenisPenghasilan string `json:"jenis_penghasilan"`

	// --- 1. KOMPONEN GAJI & TUNJANGAN ---
	GajiPokok          float64 `json:"gaji_pokok"`
	TunjSuamiIstri     float64 `json:"tunj_suami_istri"`
	TunjAnak           float64 `json:"tunj_anak"`
	TunjJabatan        float64 `json:"tunj_jabatan"`
	TunjFungsional     float64 `json:"tunj_fungsional"`
	TunjFungsionalUmum float64 `json:"tunj_fungsional_umum"`
	TunjBeras          float64 `json:"tunj_beras"`
	TunjPph            float64 `json:"tunj_pph"`
	Pembulatan         float64 `json:"pembulatan"`

	// --- 2. KOMPONEN IURAN PEMDA (Wajib untuk PPh 21) ---
	BpjsKesPemda float64 `json:"bpjs_kes_pemda"`
	JkkPemda     float64 `json:"jkk_pemda"`
	JkmPemda     float64 `json:"jkm_pemda"`

	// HASIL KALKULASI PENERIMAAN
	BrutoGaji float64 `json:"bruto_gaji"` // Total dari Poin 1 + Poin 2

	// --- 3. KOMPONEN POTONGAN PEGAWAI ---
	Iwp8Persen     float64 `json:"iwp_8_persen"` // Pensiun & THT
	Iwp1Persen     float64 `json:"iwp_1_persen"` // BPJS Pegawai
	PotonganTapera float64 `json:"potongan_tapera"`
	PotonganLain   float64 `json:"potongan_lain"`

	// HASIL KALKULASI BERSIH
	TotalPotongan float64 `json:"total_potongan"` // Total dari Poin 3
	GajiBersih    float64 `json:"gaji_bersih"`    // (Total Poin 1) - (Total Poin 3)
}
