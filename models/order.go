package models

type Order struct {
	Id      		int		`json:"id" gorm:"primary_key"`
	Name   			string	`json:"name"`
	PhoneNumber		string	`json:"phone_number"`
	Address			string	`json:"address"`
	FoodID			uint    `json:"food_id"`
	Food 			Food	`json:"food" gorm:"foreignKey:FoodID"`
	Amount      	int		`json:"amount"`
	Subtotal      	int		`json:"subtotal"`
	Discount      	int		`json:"discount"`
	Delivery      	int		`json:"delivery"`
	Total			int		`json:"total"`
}