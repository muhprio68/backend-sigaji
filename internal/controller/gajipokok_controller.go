package controller

import (
	"net/http"
	"strconv"
	"time"

	"backend-sigaji/internal/model"
	"backend-sigaji/internal/service"

	"github.com/gin-gonic/gin"
)

type GajiPokokController struct {
	service service.GajiPokokService
}

func NewGajiPokokController(service service.GajiPokokService) *GajiPokokController {
	return &GajiPokokController{service}
}
func (c *GajiPokokController) Create(ctx *gin.Context) {
	var gaji model.GajiPokok

	if err := ctx.ShouldBindJSON(&gaji); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&gaji); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gaji)
}
func (c *GajiPokokController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *GajiPokokController) GetGaji(ctx *gin.Context) {

	jenisStr := ctx.Query("jenis")

	bulanStr := ctx.Query("bulan")
	tahunStr := ctx.Query("tahun")
	tipeStr := ctx.Query("tipe")

	jenis, _ := strconv.Atoi(jenisStr)
	bulan, _ := strconv.Atoi(bulanStr)
	tahun, _ := strconv.Atoi(tahunStr)

	data, err := c.service.GetGajiByPeriode(jenis, bulan, tahun, tipeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (c *GajiPokokController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *GajiPokokController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "berhasil hapus"})
}

func (c *GajiPokokController) Import(ctx *gin.Context) {
	file, err := ctx.FormFile("file_gaji")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "file tidak ditemukan"})
		return
	}

	// 🔥 ambil bulan & tahun dari form
	bulanStr := ctx.PostForm("bulan_gaji")
	tahunStr := ctx.PostForm("tahun_gaji")
	tipeStr := ctx.PostForm("jenis_gaji")

	bulan, err := strconv.Atoi(bulanStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "bulan tidak valid"})
		return
	}

	tahun, err := strconv.Atoi(tahunStr)
	currentYear := time.Now().Year()

	if tahun > currentYear {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Tahun tidak boleh melebihi tahun sekarang",
		})
		return
	}
	if err != nil {
		ctx.JSON(400, gin.H{"error": "tahun tidak valid"})
		return
	}

	src, _ := file.Open()
	defer src.Close()

	// err = c.service.ImportExcel(src, bulan, tahun, tipeStr)
	// if err != nil {
	// 	ctx.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// ctx.JSON(200, gin.H{"message": "Import gaji berhasil"})
	err = c.service.ImportExcel(src, bulan, tahun, tipeStr)
	if err != nil {
		// DI SINI KUNCINYA:
		// Kirim status 400 (Bad Request) atau 500 dengan pesan error asli dari Service
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(), // Ini akan berisi "NIP 'xxx' pada Baris 5 tidak ditemukan..."
		})
		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Data Gaji berhasil diimport!",
	})
}
