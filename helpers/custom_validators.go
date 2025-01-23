package helpers

import (
	"rony/food-api/models"
	"rony/food-api/database"
	"github.com/go-playground/validator/v10"
)

// Fungsi untuk mendaftarkan validasi kustom
func RegisterCustomValidator() *validator.Validate {
	// Membuat instance validator
	validate := validator.New()

	// Mendaftarkan validasi kustom 'categoryidvalid'
	validate.RegisterValidation("categoryidvalid", CategoryIDValid)

	return validate
}

// Fungsi validasi kustom untuk CategoryID
func CategoryIDValid(fl validator.FieldLevel) bool {
	// Ambil category_id yang diterima oleh validasi
	categoryID := fl.Field().Uint()

	// Periksa apakah category_id ada di database
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		// Jika tidak ditemukan, kembalikan false
		return false
	}

	// Jika category ditemukan, validasi sukses
	return true
}
