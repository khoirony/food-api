package controllers

import (
	"net/http"
	"rony/food-api/database"
	"rony/food-api/helpers"
	"rony/food-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Image       string `json:"image" binding:"required"`
}

// get all category
func FindCategory(c *gin.Context) {

	//get data from database using model
	var category []models.Category
	database.DB.Find(&category)

	for i := range category {
		// Generate URL untuk gambar makanan
		if category[i].Image != "" {
			category[i].Image = GenerateS3URL(bucketName, category[i].Image)
		}
	}

	//return json
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lists Data Category",
		"data":    category,
	})
}

// store a post
func StoreCategory(c *gin.Context) {
	//validate input
	var input ValidateCategoryInput

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

	//create post
	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
		Image:       input.Image,
	}
	database.DB.Create(&category)

	//return response json
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category Created Successfully",
		"data":    category,
	})
}

// get post by id
func FindCategoryById(c *gin.Context) {
	// Ambil parameter ID dari URL
	id := c.Param("id")

	// Deklarasi variabel untuk category
	var category models.Category

	// Query untuk mencari category berdasarkan ID
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		// Jika tidak ditemukan, kembalikan respon error
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Category not found with ID: " + id,
		})
		return
	}

	// Jika ditemukan, kembalikan respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category details retrieved successfully",
		"data":    category,
	})
}

func UpdateCategoryById(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Cari data category berdasarkan ID
	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "category not found with ID: " + id,
		})
		return
	}

	// Validasi input
	var input ValidateCategoryInput
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

	// Update data category (menggunakan object category yang sudah ada)
	category.Name = input.Name
	category.Description = input.Description
	category.Image = input.Image

	// Simpan perubahan ke database
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update category",
		})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category updated successfully",
		"data":    category,
	})
}

func DeleteCategoryById(c *gin.Context) {
	// Ambil ID dari parameter
	id := c.Param("id")

	// Cari data category berdasarkan ID
	var category models.Category
	if err := database.DB.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "category not found with ID: " + id,
		})
		return
	}

	// Hapus data category
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete category",
		})
		return
	}

	// Respon sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "category deleted successfully",
		"data":    category,
	})
}
