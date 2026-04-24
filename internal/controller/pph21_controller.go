package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend-sigaji/internal/service"

	"github.com/gin-gonic/gin"
)

type PPh21Controller struct {
	Service *service.PPh21Service
}

func NewPPh21Controller(service *service.PPh21Service) *PPh21Controller {
	return &PPh21Controller{Service: service}
}

// 1. Endpoint untuk JSON (API) - Bulanan
func (c *PPh21Controller) GetPPh21Bulanan(ctx *gin.Context) {
	jenisStr := ctx.DefaultQuery("jenis", "0")
	bulanStr := ctx.Query("bulan")
	tahunStr := ctx.Query("tahun")
	tipeStr := ctx.Query("tipe")

	jenis, _ := strconv.Atoi(jenisStr)
	tahun, errTahun := strconv.Atoi(tahunStr)
	bulan, _ := strconv.Atoi(bulanStr)

	if errTahun != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Tahun wajib angka"})
		return
	}

	dataPegawai, err := c.Service.GetPPh21Bulanan(jenis, bulan, tahun, tipeStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dataPegawai})
}

// 2. Endpoint JSON Rekap A1 (PPPK) - Support Filter ID
func (c *PPh21Controller) GetRekapA1(ctx *gin.Context) {
	tahunStr := ctx.DefaultQuery("tahun", strconv.Itoa(time.Now().Year()))
	tahun, _ := strconv.Atoi(tahunStr)

	// 🔥 Tambahkan filter ID biar bisa lihat rekap per orang di UI
	idStr := ctx.Query("id")
	pegawaiID, _ := strconv.Atoi(idStr)

	results, err := c.Service.GetRekapA1Service(tahun, pegawaiID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": results})
}

// 3. Endpoint JSON Rekap A2 (PNS) - Support Filter ID
func (h *PPh21Controller) GetRekapA2(ctx *gin.Context) {
	tahunStr := ctx.Query("tahun")
	tahun, _ := strconv.Atoi(tahunStr)

	idStr := ctx.Query("id")
	pegawaiID, _ := strconv.Atoi(idStr)

	if tahun == 0 {
		tahun = time.Now().Year()
	}

	data, err := h.Service.GetRekapA2Service(tahun, pegawaiID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// 4. Export XML A1 (PPPK)
func (c *PPh21Controller) ExportXMLA1(ctx *gin.Context) {
	tahunStr := ctx.DefaultQuery("tahun", strconv.Itoa(time.Now().Year()))
	tahun, _ := strconv.Atoi(tahunStr)
	idStr := ctx.Query("id")
	pegawaiID, _ := strconv.Atoi(idStr)

	tglPotong := ctx.Query("tgl_potong")
	if tglPotong == "" {
		tglPotong = fmt.Sprintf("%d-12-31", tahun)
	}

	fileBytes, err := c.Service.GenerateXMLA1(tahun, pegawaiID, tglPotong)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fileName := fmt.Sprintf("A1_Bulk_%d.xml", tahun)
	if pegawaiID != 0 {
		fileName = fmt.Sprintf("A1_Pegawai_%d_%d.xml", pegawaiID, tahun)
	}

	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/xml")
	ctx.Data(http.StatusOK, "application/xml", fileBytes)
}

// 5. Export XML A2 (PNS)
func (h *PPh21Controller) ExportXMLA2(ctx *gin.Context) {
	tahunStr := ctx.Query("tahun")
	tahun, _ := strconv.Atoi(tahunStr)
	idStr := ctx.Query("id")
	pegawaiID, _ := strconv.Atoi(idStr)

	if tahun == 0 {
		tahun = time.Now().Year()
	}

	tglPotong := ctx.Query("tgl_potong")
	if tglPotong == "" {
		tglPotong = fmt.Sprintf("%d-12-31", tahun)
	}

	// 🔥 Panggil service dengan 2 parameter (tahun, pegawaiID)
	dataTahunan, err := h.Service.GetRekapA2Service(tahun, pegawaiID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if len(dataTahunan) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
		return
	}

	xmlBytes, err := h.Service.GenerateXMLA2(dataTahunan, tahun, tglPotong)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	filename := fmt.Sprintf("A2_Bulk_%d.xml", tahun)
	if pegawaiID != 0 {
		filename = fmt.Sprintf("A2_Pegawai_%d_%d.xml", pegawaiID, tahun)
	}

	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Type", "application/xml")
	ctx.Data(http.StatusOK, "application/xml", xmlBytes)
}

// 6. Export XML Bulanan (e-Bupot 21/26)
func (c *PPh21Controller) ExportXMLEBupot(ctx *gin.Context) {
	jenisStr := ctx.DefaultQuery("jenis", "0")
	bulanStr := ctx.Query("bulan")
	tahunStr := ctx.Query("tahun")
	tipeStr := ctx.DefaultQuery("tipe", "bulanan")
	tglPotong := ctx.Query("tgl_potong")

	jenis, _ := strconv.Atoi(jenisStr)
	bulan, _ := strconv.Atoi(bulanStr)
	tahun, _ := strconv.Atoi(tahunStr)

	if tahun == 0 || tglPotong == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tahun dan tgl_potong wajib"})
		return
	}

	xmlData, err := c.Service.GenerateXMLEBupot(jenis, bulan, tahun, tglPotong, tipeStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	golonganStr := "PNS"
	if jenis == 2 {
		golonganStr = "PPPK"
	}

	periodeStr := strings.ToUpper(tipeStr)
	if tipeStr == "bulanan" {
		periodeStr = fmt.Sprintf("%02d", bulan)
	}

	fileName := fmt.Sprintf("%s_%s_%d_EBUPOT.xml", golonganStr, periodeStr, tahun)

	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.Header("Content-Type", "application/xml")
	ctx.Data(http.StatusOK, "application/xml", xmlData)
}
