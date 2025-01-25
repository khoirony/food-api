package database

import (
	"log"
	"rony/food-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:12345678@tcp(host.docker.internal:3306)/food_app"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	MigrateDatabase(database)

	DB = database
}

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&models.Food{}, &models.Category{}, &models.Order{})
	if err != nil {
		log.Fatal("Failed to migrate database: " + err.Error())
	}
}
