package controller

import (
	"net/http"
	"strconv"
	"time"

	"backend-sigaji/internal/model"
	"backend-sigaji/internal/service"

	"github.com/gin-gonic/gin"
)

type TPPController struct {
	service service.TPPService
}

func NewTPPController(service service.TPPService) *TPPController {
	return &TPPController{service}
}

func (c *TPPController) Create(ctx *gin.Context) {
	var tpp model.TPP

	if err := ctx.ShouldBindJSON(&tpp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.service.CreateTPP(&tpp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "TPP berhasil dibuat",
		"data":    tpp,
	})
}

func (c *TPPController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAllTPP()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *TPPController) GetTPP(ctx *gin.Context) {

	jenisStr := ctx.Query("jenis")

	bulanStr := ctx.Query("bulan")
	tahunStr := ctx.Query("tahun")
	tipeStr := ctx.Query("tipe")

	jenis, _ := strconv.Atoi(jenisStr)
	bulan, _ := strconv.Atoi(bulanStr)
	tahun, _ := strconv.Atoi(tahunStr)

	data, err := c.service.GetTppByPeriode(jenis, bulan, tahun, tipeStr)
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

func (c *TPPController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	data, err := c.service.GetTPPByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *TPPController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	var tpp model.TPP
	if err := ctx.ShouldBindJSON(&tpp); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tpp.ID = uint(id)

	if err := c.service.UpdateTPP(&tpp); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "TPP berhasil diupdate",
	})
}

func (c *TPPController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	if err := c.service.DeleteTPP(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "TPP berhasil dihapus",
	})
}

func (c *TPPController) Import(ctx *gin.Context) {
	file, err := ctx.FormFile("file_tpp")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "file tidak ditemukan"})
		return
	}

	// 🔥 ambil bulan & tahun dari form
	bulanStr := ctx.PostForm("bulan_tpp")
	tahunStr := ctx.PostForm("tahun_tpp")
	tipeStr := ctx.PostForm("jenis_tpp")

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

	err = c.service.ImportExcel(src, bulan, tahun, tipeStr)
	if err != nil {
		// Gunakan status 400 agar SweetAlert tahu ini kesalahan input (NIP tidak ditemukan, dll)
		ctx.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(), // Akan berisi pesan "NIP 'xxx' pada Baris 10 tidak ditemukan..."
		})
		return
	}

	// Jika berhasil
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Import TPP berhasil dilakukan!",
	})
}
