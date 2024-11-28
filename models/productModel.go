package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string
	Price int
	Image string
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
