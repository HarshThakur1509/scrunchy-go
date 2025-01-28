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
	CartID    uint `gorm:"foreignKey:CartID"`
	ProductID uint
	Product   Product
}

type Cart struct {
	gorm.Model
	UserID    uint `gorm:"foreignKey:UserID"`
	Total     uint
	CartItems []CartItem
}
