package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name  string
	Price uint
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
	Total     uint
	CartItems []CartItem
}
