package repository

import (
	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/model"

	"gorm.io/gorm"
)

type PPh21Repository struct {
	DB *gorm.DB
}

func NewPPh21Repository(db *gorm.DB) *PPh21Repository {
	return &PPh21Repository{DB: db}
}

// --- REPOSITORY PENGHASILAN BULANAN ---
func (r *PPh21Repository) GetPenghasilanBulanan(jenis int, bulan int, tahun int, tipe string) (dto.PPh21PageResponse, error) {
	var listPegawai []dto.PPh21BulananResponse
	var finalResponse dto.PPh21PageResponse

	if tipe == "" {
		tipe = "bulanan"
	}

	// 1. SIAPKAN JOIN DINAMIS
	joinGaji := "LEFT JOIN gaji_pokok ON gaji_pokok.id_pegawai = pegawai.id AND gaji_pokok.tahun = ? AND gaji_pokok.jenis_penghasilan = ?"
	joinTPP := "LEFT JOIN tpp ON tpp.id_pegawai = pegawai.id AND tpp.tahun = ? AND tpp.jenis_penghasilan = ?"

	argsGaji := []interface{}{tahun, tipe}
	argsTPP := []interface{}{tahun, tipe}

	if tipe == "bulanan" {
		joinGaji += " AND gaji_pokok.bulan = ?"
		argsGaji = append(argsGaji, bulan)

		joinTPP += " AND tpp.bulan = ?"
		argsTPP = append(argsTPP, bulan)
	}

	// 2. BANGUN QUERY UTAMA
	dbQuery := r.DB.Table("pegawai").
		Select(`pegawai.id, pegawai.nip, pegawai.nik, pegawai.nama, pegawai.jabatan, pegawai.status_ptkp, 
                
                COALESCE(gaji_pokok.bulan, tpp.bulan, 0) as bulan_pajak,
                
                (COALESCE(gaji_pokok.gaji_pokok, 0) + 
                 COALESCE(gaji_pokok.tunj_suami_istri, 0) + 
                 COALESCE(gaji_pokok.tunj_anak, 0) + 
                 COALESCE(gaji_pokok.tunj_jabatan, 0) + 
                 COALESCE(gaji_pokok.tunj_fungsional, 0) + 
                 COALESCE(gaji_pokok.tunj_fungsional_umum, 0) + 
                 COALESCE(gaji_pokok.tunj_beras, 0) + 
                 COALESCE(gaji_pokok.pembulatan, 0) + 
                 COALESCE(gaji_pokok.bpjs_kes_pemda, 0) + 
                 COALESCE(gaji_pokok.jkk_pemda, 0) + 
                 COALESCE(gaji_pokok.jkm_pemda, 0)) AS bruto_gaji,
                
                (COALESCE(tpp.tpp_beban, 0) + 
                 COALESCE(tpp.tpp_prestasi, 0) + 
                 COALESCE(tpp.tpp_kondisi, 0) + 
                 COALESCE(tpp.bpjs4, 0)) AS bruto_tpp,
                
                (COALESCE(gaji_pokok.gaji_pokok, 0) + COALESCE(gaji_pokok.tunj_suami_istri, 0) + COALESCE(gaji_pokok.tunj_anak, 0) + COALESCE(gaji_pokok.tunj_jabatan, 0) + COALESCE(gaji_pokok.tunj_fungsional, 0) + COALESCE(gaji_pokok.tunj_fungsional_umum, 0) + COALESCE(gaji_pokok.tunj_beras, 0) + COALESCE(gaji_pokok.tunj_pph, 0) + COALESCE(gaji_pokok.pembulatan, 0) + COALESCE(gaji_pokok.bpjs_kes_pemda, 0) + COALESCE(gaji_pokok.jkk_pemda, 0) + COALESCE(gaji_pokok.jkm_pemda, 0) + 
                 COALESCE(tpp.tpp_beban, 0) + COALESCE(tpp.tpp_prestasi, 0) + COALESCE(tpp.tpp_kondisi, 0) + COALESCE(tpp.bpjs4, 0)) AS total_bruto`).
		Joins(joinGaji, argsGaji...).
		Joins(joinTPP, argsTPP...)

	// 3. FILTER LENGKAP
	if tipe == "bulanan" {
		dbQuery = dbQuery.Where("COALESCE(gaji_pokok.bulan, tpp.bulan, 0) = ?", bulan)
	}

	dbQuery = dbQuery.Where("COALESCE(gaji_pokok.tahun, tpp.tahun) = ?", tahun)
	dbQuery = dbQuery.Where("COALESCE(gaji_pokok.jenis_penghasilan, tpp.jenis_penghasilan) = ?", tipe)
	dbQuery = dbQuery.Where("gaji_pokok.id IS NOT NULL OR tpp.id IS NOT NULL")

	if jenis == 1 {
		dbQuery = dbQuery.Where("pegawai.status_asn IN ?", []int{1, 3})
	} else if jenis != 0 {
		dbQuery = dbQuery.Where("pegawai.status_asn = ?", jenis)
	}

	// 4. EKSEKUSI QUERY
	err := dbQuery.Debug().Scan(&listPegawai).Error
	if err != nil {
		return finalResponse, err
	}

	// 5. HITUNG SUMMARY
	var sumBruto, sumPPh float64
	for _, p := range listPegawai {
		sumBruto += p.TotalBruto
		sumPPh += p.PPh21
	}

	finalResponse.ListPegawai = listPegawai
	finalResponse.Summary.TotalPegawai = len(listPegawai)
	finalResponse.Summary.GrandTotalBruto = sumBruto
	finalResponse.Summary.GrandTotalPPh21 = sumPPh

	return finalResponse, nil
}

// --- REPOSITORY TARIF TER ---
func (r *PPh21Repository) GetAllTarifTER() ([]model.RefTarifTer, error) {
	var tarifList []model.RefTarifTer
	err := r.DB.Find(&tarifList).Error
	return tarifList, err
}

// --- REPOSITORY A1 (PPPK) ---
func (r *PPh21Repository) GetRekapA1(tahun int, pegawaiID int) ([]dto.PPh21A1Response, error) {
	var results []dto.PPh21A1Response

	query := `
		SELECT 
			p.id, p.nip, p.nik, p.nama, p.jabatan, p.status_ptkp,
			COALESCE(MIN(g.bulan), 0) AS bulan_awal,
			COALESCE(MAX(g.bulan), 0) AS bulan_akhir,
			
			-- Komponen Gaji (Salary)
			COALESCE(SUM(CASE WHEN g.jenis_penghasilan = 'bulanan' THEN g.gaji_pokok + g.tunj_suami_istri + g.tunj_anak + g.tunj_beras ELSE 0 END), 0) AS salary_pension_jht_tht,
			COALESCE(SUM(CASE WHEN g.jenis_penghasilan = 'bulanan' THEN g.tunj_pph ELSE 0 END), 0) AS income_tax_benefit,
			COALESCE(SUM(CASE WHEN g.jenis_penghasilan = 'bulanan' THEN g.bpjs_kes_pemda + g.jkk_pemda + g.jkm_pemda ELSE 0 END), 0) AS insurance_paid_by_emp,
			
			-- Other Benefit (Tunj Jabatan + TPP Bulanan)
			COALESCE(SUM(CASE WHEN g.jenis_penghasilan = 'bulanan' THEN g.tunj_jabatan + g.tunj_fungsional + g.tunj_fungsional_umum ELSE 0 END), 0) 
			+ COALESCE(t.tpp_rutin, 0) AS other_benefit,
			
			-- THR & Gaji 13 (Gaji + TPP)
			COALESCE(SUM(CASE WHEN g.jenis_penghasilan IN ('thr', 'gaji_13') THEN 
				(g.gaji_pokok + g.tunj_suami_istri + g.tunj_anak + g.tunj_jabatan + g.tunj_fungsional + g.tunj_fungsional_umum + g.tunj_beras + g.tunj_pph + g.pembulatan + g.bpjs_kes_pemda + g.jkk_pemda + g.jkm_pemda) 
			ELSE 0 END), 0) 
			+ COALESCE(t.tpp_thr_13, 0) AS tantiem_bonus_thr,
			
			COALESCE(SUM(g.iwp8_persen), 0) AS pension_contribution,
			COALESCE(SUM(g.gaji_pokok + g.tunj_suami_istri + g.tunj_anak + g.tunj_jabatan + g.tunj_fungsional + g.tunj_fungsional_umum + g.tunj_beras + g.tunj_pph + g.pembulatan + g.bpjs_kes_pemda + g.jkk_pemda + g.jkm_pemda), 0) 
			+ COALESCE(t.total_tpp_all, 0) AS total_bruto

		FROM pegawai p
		LEFT JOIN (
			-- Subquery Gaji dengan bulan_akhir
			SELECT id_pegawai, bulan, tahun, jenis_penghasilan,
				   gaji_pokok, tunj_suami_istri, tunj_anak, tunj_beras, tunj_pph, 
				   bpjs_kes_pemda, jkk_pemda, jkm_pemda, tunj_jabatan, tunj_fungsional, 
				   tunj_fungsional_umum, pembulatan, iwp8_persen,
				   MAX(bulan) OVER(PARTITION BY id_pegawai) as bulan_akhir_pegawai
			FROM gaji_pokok WHERE tahun = ?
		) g ON p.id = g.id_pegawai
		LEFT JOIN (
			SELECT id_pegawai, 
				SUM(CASE WHEN jenis_penghasilan = 'bulanan' THEN tpp_beban + tpp_prestasi + tpp_kondisi + bpjs4 ELSE 0 END) as tpp_rutin,
				SUM(CASE WHEN jenis_penghasilan IN ('thr', 'gaji_13') THEN tpp_beban + tpp_prestasi + tpp_kondisi + bpjs4 ELSE 0 END) as tpp_thr_13,
				SUM(tpp_beban + tpp_prestasi + tpp_kondisi + bpjs4) as total_tpp_all
			FROM tpp WHERE tahun = ? GROUP BY id_pegawai
		) t ON p.id = t.id_pegawai
		
		WHERE p.status_asn = 2 
		AND (g.id_pegawai IS NOT NULL OR t.total_tpp_all > 0)
		
		-- 🔥 LOGIKA PAGAR BULAN TERAKHIR (DINAMIS) 🔥
		AND (
			(? != 0 AND p.id = ?) -- Jika ID diisi (cetak per orang), lupakan filter bulan
			OR 
			(? = 0 AND g.bulan_akhir_pegawai = (SELECT MAX(bulan) FROM gaji_pokok WHERE tahun = ?)) -- Jika ID 0 (cetak masal)
		)
		
		GROUP BY p.id
		ORDER BY p.nama ASC
	`

	err := r.DB.Raw(query,
		tahun,     // Untuk subquery gaji
		tahun,     // Untuk subquery tpp
		pegawaiID, // Untuk ? != 0
		pegawaiID, // Untuk p.id = ?
		pegawaiID, // Untuk ? = 0
		tahun,     // Untuk SELECT MAX(bulan)
	).Scan(&results).Error

	return results, err
}

// --- REPOSITORY A2 (PNS & CPNS) ---
func (r *PPh21Repository) GetRekapA2(tahun int, pegawaiID int) ([]dto.PPh21TahunanResponse, error) {
	var results []dto.PPh21TahunanResponse

	query := `
		SELECT 
			p.id, p.nip, p.nik, p.nama, p.jabatan, p.golongan, p.status_ptkp,
			
			-- Ambil data bulan dari subquery gaji (g)
			COALESCE(g.bulan_awal, 0) AS bulan_awal,
			COALESCE(g.bulan_akhir, 0) AS bulan_akhir,
			
			COALESCE(g.total_gaji_pokok, 0) AS total_gaji_pokok,
			COALESCE(g.total_tunj_istri, 0) AS total_tunj_istri,
			COALESCE(g.total_tunj_anak, 0) AS total_tunj_anak,
			COALESCE(g.total_tunj_beras, 0) AS total_tunj_beras,
			COALESCE(g.total_tunj_jabatan, 0) AS total_tunj_jabatan,
			COALESCE(g.total_tunj_lain, 0) AS total_tunj_lain,
			COALESCE(g.total_iwp, 0) AS total_iwp,
			COALESCE(t.total_tpp, 0) AS total_tpp,
			COALESCE(t.total_bpjs4_tpp, 0) AS total_bpjs4_tpp
		FROM pegawai p
		LEFT JOIN (
			SELECT id_pegawai, 
				MIN(bulan) as bulan_awal,
				MAX(bulan) as bulan_akhir,
				SUM(gaji_pokok) as total_gaji_pokok,
				SUM(tunj_suami_istri) as total_tunj_istri,
				SUM(tunj_anak) as total_tunj_anak,
				SUM(tunj_beras) as total_tunj_beras,
				SUM(tunj_jabatan + tunj_fungsional + tunj_fungsional_umum) as total_tunj_jabatan,
				SUM(pembulatan + bpjs_kes_pemda + jkk_pemda + jkm_pemda) as total_tunj_lain,
				SUM(iwp8_persen) as total_iwp 
			FROM gaji_pokok 
			WHERE tahun = ? 
			GROUP BY id_pegawai
		) g ON p.id = g.id_pegawai
		LEFT JOIN (
			SELECT id_pegawai, 
				SUM(tpp_beban + tpp_prestasi + tpp_kondisi) as total_tpp,
				SUM(bpjs4) as total_bpjs4_tpp
			FROM tpp 
			WHERE tahun = ? 
			GROUP BY id_pegawai
		) t ON p.id = t.id_pegawai
		
		WHERE p.status_asn IN (1, 3) 
		AND (g.total_gaji_pokok > 0 OR t.total_tpp > 0 OR t.total_bpjs4_tpp > 0) 
		
		-- 🔥 LOGIKA PAGAR BULAN TERAKHIR (DINAMIS) 🔥
		AND (
			(? != 0 AND p.id = ?) -- Jika ID diisi (cetak per orang), lupakan filter bulan
			OR 
			(? = 0 AND g.bulan_akhir = (SELECT MAX(bulan) FROM gaji_pokok WHERE tahun = ?)) -- Jika ID 0 (cetak masal)
		)

		ORDER BY p.nama ASC
	`

	err := r.DB.Raw(query,
		tahun,     // Untuk subquery gaji
		tahun,     // Untuk subquery tpp
		pegawaiID, // Untuk ? != 0
		pegawaiID, // Untuk p.id = ?
		pegawaiID, // Untuk ? = 0
		tahun,     // Untuk SELECT MAX(bulan)
	).Scan(&results).Error

	return results, err
}
