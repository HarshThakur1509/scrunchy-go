package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string
	Price int
}

type CartItem struct {
	gorm.Model
	Quantity  uint
	CartID    uint
	ProductID uint
	Product   Product
}

type Cart struct {
	gorm.Model
	UserID    uint
	CartItems []CartItem
}
