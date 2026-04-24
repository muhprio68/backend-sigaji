package main

import (
	"backend-sigaji/internal/config"
	"backend-sigaji/internal/controller"
	"backend-sigaji/internal/repository"
	"backend-sigaji/internal/router"
	"backend-sigaji/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Waduh, file .env tidak ditemukan bre!")
	}

	// 1️⃣ Init Gin
	r := gin.Default()

	// 2️⃣ Init Database
	db := config.InitDB()

	// 3️⃣ Init Repository
	pegawaiRepo := repository.NewPegawaiRepository(db)
	gajiRepo := repository.NewGajiPokokRepository(db)
	tppRepo := repository.NewTPPRepository(db)
	pph21Repo := repository.NewPPh21Repository(db)
	realisasiRepo := repository.NewRealisasiRepository(db)

	// 4️⃣ Init Service
	pegawaiService := service.NewPegawaiService(db, pegawaiRepo)
	gajiService := service.NewGajiPokokService(db, gajiRepo)
	tppService := service.NewTPPService(db, tppRepo)
	ppph21Service := service.NewPPh21Service(db, pph21Repo)
	realisasiService := service.NewRealisasiService(db, realisasiRepo)
	authService := service.NewAuthService(db)

	// 5️⃣ Init Controller
	pegawaiController := controller.NewPegawaiController(pegawaiService)
	gajiController := controller.NewGajiPokokController(gajiService)
	tppController := controller.NewTPPController(tppService)
	pph21Controller := controller.NewPPh21Controller(ppph21Service)
	realisasiController := controller.NewRealisasiController(realisasiService)
	authController := controller.NewAuthController(authService)

	// 6️⃣ Init Routes
	router.InitRoutes(r, authController, pegawaiController, gajiController, tppController, pph21Controller, realisasiController)

	r.Run(":8080")
}
