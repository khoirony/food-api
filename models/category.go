package models

type Category struct {
	Id      	int    `json:"id" gorm:"primary_key"`
	Name   		string `json:"name"`
	Description string `json:"description"`
	Image       string  `json:"image"`
	Foods []Food `json:"foods" gorm:"foreignKey:CategoryID"` // Menambahkan relasi dengan Food
}