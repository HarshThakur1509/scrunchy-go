package models

import "gorm.io/gorm"

// Product Model
type Product struct {
	gorm.Model
	Name  string `gorm:"not null" `
	Price uint   `gorm:"not null" `
	Image string
}

// Cart Model
type Cart struct {
	gorm.Model
	UserID    uint       `gorm:"uniqueIndex;not null" `
	Total     uint       `gorm:"default:0" `
	CartItems []CartItem `gorm:"constraint:OnDelete:CASCADE;" `
}

// CartItem Model
type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"index;not null;foreignKey:CartID" `
	ProductID uint    `gorm:"not null;foreignKey:ProductID" `
	Product   Product `gorm:"constraint:OnDelete:CASCADE;" `
	Quantity  uint    `gorm:"not null;default:1" `
}
