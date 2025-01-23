package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	// "time"
	"fmt"

	"rony/food-api/database"
	"rony/food-api/helpers"
	"rony/food-api/models"

	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidateFoodInput struct {
	Name        string  `json:"name" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"` // Menambahkan validasi kustom
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`       // Wajib diisi, harus lebih besar dari 0
	Rate        float64 `json:"rate" binding:"required,gte=0,lte=5"` // Wajib diisi, antara 0-5
	Image       string  `json:"image" binding:"required"`
}

// Define S3 bucket and region
const (
	bucketName = "food-app-khoirony"
	region     = "ap-southeast-2" // Ganti dengan region S3 Anda
	s3BaseURL  = "https://%s.s3.%s.amazonaws.com/"
)

// S3 client
var s3Client *s3.Client

// Initialize S3 client
func InitS3Client(region string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load AWS SDK config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

// GenerateS3URL generates a direct URL for an S3 object
func GenerateS3URL(bucketName, objectKey string) string {
	// Format the URL
	url := fmt.Sprintf(s3BaseURL, bucketName, region) + objectKey
	return url
}

// FindFoods retrieves a list of foods from the database
func FindFoods(c *gin.Context) {
	// Ambil query parameter untuk filter dan pencarian
	categoryIDStr := c.Query("category_id")
	searchName := c.Query("name")
	id := c.Query("id")

	// Inisialisasi query
	var foods []models.Food
	query := database.DB.Preload("Category")

	// Tambahkan filter berdasarkan kategori jika category_id diberikan
	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			if(categoryID != 0){
				query = query.Where("category_id = ?", categoryID)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid category_id parameter",
			})
			return
		}
	}

	// Tambahkan filter pencarian berdasarkan nama makanan jika name diberikan
	if searchName != "" {
		query = query.Where("name LIKE ?", "%"+searchName+"%")
	}

	if id != "" {
		query = query.Where("id = ?", id)
	}

	// Eksekusi query
	if err := query.Find(&foods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve foods with categories",
		})
		return
	}

	// Generate full S3 URL untuk setiap gambar
	for i := range foods {
		// Generate URL untuk gambar makanan
		if foods[i].Image != "" {
			foods[i].Image = GenerateS3URL(bucketName, foods[i].Image)
		}

		// Generate URL untuk gambar kategori
		if foods[i].Category.Image != "" {
			foods[i].Category.Image = GenerateS3URL(bucketName, foods[i].Category.Image)
		}
	}

	// Return hasil dalam JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lists Data foods",
		"data":    foods,
	})
}


// store a post
func StoreFood(c *gin.Context) {
	//validate input
	var input ValidateFoodInput

	if err := c.ShouldBindJSON(&input); err != nil {
		// Parsing error validasi
		errors := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, helpers.GetErrorMsg(e))
		}

		// Jika validasi gagal, kembalikan respon error
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Cari data category berdasarkan CategoryID
	var category models.Category
	if err := database.DB.Where("id = ?", input.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "category_id " + strconv.Itoa(int(input.CategoryID)) + " not found",
		})
		return
	}

	// Create post (food)
	food := models.Food{
		Name:        input.Name,
		CategoryID:  input.CategoryID,
		Description: input.Description,
		Price:       input.Price,
		Rate:        input.Rate,
		Image:       input.Image,
	}

	// Save food to database
	if err := database.DB.Create(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create food",
		})
		return
	}

	// Preload category to ensure category is included in the response
	if err := database.DB.Preload("Category").First(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to preload category",
		})
		return
	}

	// Return response JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Food Created Successfully",
		"data":    food,
	})
}


// get post by id
func FindFoodById(c *gin.Context) {
	// Ambil parameter ID dari URL
	id := c.Param("id")

	// Deklarasi variabel untuk food
	var food models.Food
	query := database.DB.Preload("Category")

	// Query untuk mencari food berdasarkan ID
	if err := database.DB.Where("id = ?", id).First(&food).Error; err != nil {
		// Jika tidak ditemukan, kembalikan respon error
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Food not found with ID: " + id,
		})
		return
	}

	// Eksekusi query
	if err := query.Find(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve foods with categories",
		})
		return
	}

	food.Image = GenerateS3URL(bucketName, food.Image)
	food.Category.Image = GenerateS3URL(bucketName, food.Category.Image)
	// Jika ditemukan, kembalikan respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Food details retrieved successfully",
		"data":    food,
	})
}

func UpdateFoodById(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Cari data food berdasarkan ID
	var food models.Food
	if err := database.DB.Where("id = ?", id).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Food not found with ID: " + id,
		})
		return
	}

	// Validasi input
	var input ValidateFoodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, e.Field()+" is invalid")
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Cari data food berdasarkan ID
	var category models.Category
	if err := database.DB.Where("id = ?", input.CategoryID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "category_id " + strconv.Itoa(int(input.CategoryID)) + " not found",
		})
		return
	}

	// Update data food (menggunakan object food yang sudah ada)
	food.Name = input.Name
	food.Description = input.Description
	food.CategoryID = input.CategoryID
	food.Price = input.Price
	food.Rate = input.Rate
	food.Image = input.Image

	// Simpan perubahan ke database
	if err := database.DB.Save(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update food",
		})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Food updated successfully",
		"data":    food,
	})
}

func DeleteFoodById(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Cari data food berdasarkan ID
	var food models.Food
	if err := database.DB.Where("id = ?", id).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Food not found with ID: " + id,
		})
		return
	}

	// Hapus data food
	if err := database.DB.Delete(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete food",
		})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Food deleted successfully",
		"data":    food,
	})
}