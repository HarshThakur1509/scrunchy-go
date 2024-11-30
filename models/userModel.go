package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Name        string
	Phone       string
	Admin       bool
	ResetToken  string
	TokenExpiry time.Time
	Cart        Cart
	Address     Address
}

type Address struct {
	gorm.Model
	UserID  uint
	City    string
	State   string
	PinCode string
	Address string
}
