package main

import (
	"rony/food-api/controllers"
	"rony/food-api/database"
	"rony/food-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	//inisialiasai Gin
	router := gin.Default()

	//panggil koneksi database
	database.ConnectDatabase()

	// Migrasi database untuk memperbarui tabel
	database.MigrateDatabase(database.DB)

	// Inisialisasi S3 client
	region := "ap-southeast-2" // Ganti sesuai kebutuhan Anda
	controllers.InitS3Client(region)

	// Panggil fungsi untuk mendaftarkan route
	routes.RegisterRoutes(router)

	//mulai server dengan port 3000
	router.Run(":3000")
}