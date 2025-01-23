package models

type Food struct {
	Id      	int    `json:"id" gorm:"primary_key"`
	Name   		string `json:"name"`
	CategoryID  uint    `json:"category_id"` // Foreign key yang merujuk ke Category
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"` // Relasi dengan Category
	Description string `json:"description"`
	Price       float64 `json:"price" gorm:"type:decimal(20,2)"` // 2 angka di belakang koma
	Rate        float64 `json:"rate" gorm:"type:decimal(2,1)"`   // 1 angka di belakang koma
	Image       string  `json:"image"`
}