package router

import (
	"backend-sigaji/internal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	r *gin.Engine,
	authController *controller.AuthController, // 🔥 Tambahan Controller Login
	pegawaiController *controller.PegawaiController,
	gajiController *controller.GajiPokokController,
	tppController *controller.TPPController,
	pph21Controller *controller.PPh21Controller,
	realisasiController *controller.RealisasiController,
) {

	// ==========================================
	// 1. SETUP CORS (Wajib untuk Frontend terpisah)
	// ==========================================
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Mengizinkan semua port (termasuk 5500 Live Server)
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	api := r.Group("/api")
	{
		// ==========================================
		// 2. PUBLIC ROUTES (Tanpa Satpam JWT)
		// ==========================================
		api.POST("/login", authController.Login)

		// ==========================================
		// 3. PROTECTED ROUTES (Dijaga Satpam JWT)
		// ==========================================
		protected := api.Group("/")
		protected.Use(AuthMiddleware()) // 🔥 Panggil file middleware.go yang tadi dibikin
		{
			pegawai := protected.Group("/pegawai")
			{
				pegawai.GET("", pegawaiController.GetAll)
				pegawai.POST("", pegawaiController.Create)
				pegawai.POST("/import", pegawaiController.Import)
				pegawai.GET("/:id", pegawaiController.GetByID)
				pegawai.PUT("/:id", pegawaiController.Update)
				pegawai.DELETE("/:id", pegawaiController.Delete)
			}

			gaji := protected.Group("/gajipokok")
			{
				gaji.POST("", gajiController.Create)
				gaji.GET("", gajiController.GetGaji)
				gaji.GET("/all", gajiController.GetAll)
				gaji.GET("/:id", gajiController.GetByID)
				gaji.DELETE("/:id", gajiController.Delete)
				gaji.POST("/import", gajiController.Import)
			}

			tpp := protected.Group("/tpp")
			{
				tpp.POST("", tppController.Create)
				tpp.GET("", tppController.GetTPP)
				tpp.GET("/all", tppController.GetAll)
				tpp.GET("/:id", tppController.GetByID)
				tpp.PUT("/:id", tppController.Update)
				tpp.DELETE("/:id", tppController.Delete)
				tpp.POST("/import", tppController.Import)
			}

			pph21 := protected.Group("/pph21")
			{
				pph21.GET("", pph21Controller.GetPPh21Bulanan)
				pph21.GET("/export-xmlbulan", pph21Controller.ExportXMLEBupot)
				pph21.GET("/tahunan/a1", pph21Controller.GetRekapA1)
				pph21.GET("/tahunan-export-a1", pph21Controller.ExportXMLA1)
				pph21.GET("/tahunan/a2", pph21Controller.GetRekapA2)
				pph21.GET("/tahunan-export-a2", pph21Controller.ExportXMLA2)
			}

			realisasi := protected.Group("/realisasi")
			{
				realisasi.GET("", realisasiController.GetPenyerapan)
			}
		}
	}
}
