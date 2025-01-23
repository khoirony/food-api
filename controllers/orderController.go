package controllers

import (
	"net/http"
	"strconv"

	"rony/food-api/database"
	"rony/food-api/helpers"
	"rony/food-api/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidateOrderInput struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
	FoodID      uint   `json:"food_id" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Subtotal    int    `json:"subtotal" binding:"required"`
	Discount    int    `json:"discount" binding:"required"`
	Delivery    int    `json:"delivery" binding:"required"`
	Total       int    `json:"total" binding:"required"`
}

func FindOrders(c *gin.Context) {

	// Inisialisasi query
	var orders []models.Order
	query := database.DB.Preload("Food").Preload("Food.Category").Find(&orders)

	// Eksekusi query
	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve orders with foods",
		})
		return
	}

	for i := range orders {
		if orders[i].Food.Image != "" {
			orders[i].Food.Image = GenerateS3URL(bucketName, orders[i].Food.Image)
		}

		if orders[i].Food.Category.Image != "" {
			orders[i].Food.Category.Image = GenerateS3URL(bucketName, orders[i].Food.Category.Image)
		}
	}

	// Return hasil dalam JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lists Data Orders",
		"data":    orders,
	})
}

// store a post
func StoreOrder(c *gin.Context) {
	//validate input
	var input ValidateOrderInput

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

	// Cari data food berdasarkan FoodID
	var food models.Food
	if err := database.DB.Where("id = ?", input.FoodID).First(&food).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "food_id " + strconv.Itoa(int(input.FoodID)) + " not found",
		})
		return
	}

	// Create order
	order := models.Order{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Address:     input.Address,
		FoodID:      input.FoodID,
		Amount:      input.Amount,
		Subtotal:    input.Subtotal,
		Discount:    input.Discount,
		Delivery:    input.Delivery,
		Total:       input.Total,
	}

	// Save order to database
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create order",
		})
		return
	}

	// Preload food to ensure food is included in the response
	if err := database.DB.Preload("Food").First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to preload food",
		})
		return
	}

	// Return response JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order Created Successfully",
		"data":    order,
	})
}
