package controller

import (
	"net/http"
	"strconv"

	"backend-sigaji/internal/dto"
	"backend-sigaji/internal/model"
	"backend-sigaji/internal/service"

	"github.com/gin-gonic/gin"
)

type PegawaiController struct {
	service *service.PegawaiService
}

func NewPegawaiController(s *service.PegawaiService) *PegawaiController {
	return &PegawaiController{service: s}
}

func (h *PegawaiController) Create(c *gin.Context) {
	var req dto.CreatePegawaiRequest

	// Validasi JSON Body sesuai struktur DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Data yang dikirim tidak valid atau kurang lengkap",
			"error":   err.Error(),
		})
		return
	}

	// Teruskan data ke Service
	err := h.service.CreatePegawai(req)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data pegawai ke database",
			"error":   err.Error(),
		})
		return
	}

	// Berhasil
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Data pegawai berhasil ditambahkan",
	})
}

func (c *PegawaiController) GetAll(ctx *gin.Context) {
	var filter model.PegawaiFilter

	statusStr := ctx.Query("jenis")
	if statusStr != "" {
		statusInt, err := strconv.ParseUint(statusStr, 10, 8)
		if err == nil {
			status := uint8(statusInt)
			filter.StatusAsn = &status
		}
	}

	data, err := c.service.GetAllPegawai(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *PegawaiController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.Atoi(idParam)

	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (c *PegawaiController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input model.Pegawai
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id)

	err = c.service.UpdatePegawai(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil diupdate",
	})
}

func (c *PegawaiController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	err = c.service.DeletePegawai(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data berhasil dihapus",
	})
}

func (c *PegawaiController) Import(ctx *gin.Context) {
	file, err := ctx.FormFile("file_pegawai")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "file tidak ditemukan"})
		return
	}

	src, err := file.Open()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	err = c.service.ImportExcel(src)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Import berhasil"})
}
