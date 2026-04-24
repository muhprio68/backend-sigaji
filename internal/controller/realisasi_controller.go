package controller

import (
	"net/http"
	"strconv"

	"time"

	"backend-sigaji/internal/service"

	"github.com/gin-gonic/gin"
)

type RealisasiController struct {
	Service *service.RealisasiService
}

func NewRealisasiController(service *service.RealisasiService) *RealisasiController {
	return &RealisasiController{Service: service}
}

func (c *RealisasiController) GetPenyerapan(ctx *gin.Context) {
	// Default ke tahun sekarang kalau gak dikirim
	tahunStr := ctx.DefaultQuery("tahun", strconv.Itoa(time.Now().Year()))

	// Default bulan = 0 (berarti narik semua bulan di tahun tersebut)
	bulanStr := ctx.DefaultQuery("bulan", "0")

	tahun, _ := strconv.Atoi(tahunStr)
	bulan, _ := strconv.Atoi(bulanStr)

	data, err := c.Service.GetPenyerapanAnggaran(tahun, bulan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data penyerapan: " + err.Error(),
		})
		return
	}
	// PANGGIL SERVICE UNTUK NGITUNG PEGAWAI
	countPNS, countPPPK, _ := c.Service.GetCountPegawaiAktif(bulan, tahun)

	// 3. 🔥 PANGGIL WAKTU UPDATE DARI SERVICE 🔥
	lastUpdate, _ := c.Service.GetLastUpdate()

	var lastUpdateStr string
	if !lastUpdate.IsZero() {
		lastUpdateStr = lastUpdate.Format(time.RFC3339) // Format khusus buat Javascript
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"last_updated": lastUpdateStr,
		"total_pns":    countPNS,  // Data dikirim ke JSON
		"total_pppk":   countPPPK, // Data dikirim ke JSON
		"data":         data,
	})
}
