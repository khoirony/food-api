package routes

import (
	"rony/food-api/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes digunakan untuk mendaftarkan semua route
func RegisterRoutes(router *gin.Engine) {
	// Route untuk general (tidak termasuk dalam kelompok tertentu)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// Kelompok route untuk "foods"
	foodRoutes := router.Group("/api/foods")
	{
		foodRoutes.GET("", controllers.FindFoods)             // GET /api/foods
		foodRoutes.POST("", controllers.StoreFood)            // POST /api/foods
		foodRoutes.GET("/:id", controllers.FindFoodById)      // GET /api/foods/:id
		foodRoutes.PUT("/:id", controllers.UpdateFoodById)    // PUT /api/foods/:id
		foodRoutes.DELETE("/:id", controllers.DeleteFoodById) // DELETE /api/foods/:id
	}

	// Kelompok route untuk "category"
	categoryRoutes := router.Group("/api/category")
	{
		categoryRoutes.GET("", controllers.FindCategory)              // GET /api/category
		categoryRoutes.POST("", controllers.StoreCategory)            // POST /api/category
		categoryRoutes.GET("/:id", controllers.FindCategoryById)      // GET /api/category/:id
		categoryRoutes.PUT("/:id", controllers.UpdateCategoryById)    // PUT /api/category/:id
		categoryRoutes.DELETE("/:id", controllers.DeleteCategoryById) // DELETE /api/category/:id
	}

	// Kelompok route untuk "order"
	orderRoutes := router.Group("/api/order")
	{
		orderRoutes.GET("", controllers.FindOrders)  // GET /api/category
		orderRoutes.POST("", controllers.StoreOrder) // POST /api/category
	}

	ocrRoutes := router.Group("/api/ocr2")
	{
		ocrRoutes.POST("", controllers.OcrToText)
	}
}
