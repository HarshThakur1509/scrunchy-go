package models

import (
	"time"

	"gorm.io/gorm"
)

// User Model
type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null" `
	Password    string `gorm:"not null" json:"-"`
	Name        string `gorm:"not null"`
	Phone       string
	Admin       bool      `gorm:"default:false"`
	ResetToken  string    `json:"-"`
	TokenExpiry time.Time `json:"-"`
	Cart        Cart      `gorm:"constraint:OnDelete:CASCADE;" `
	Address     Address   `gorm:"constraint:OnDelete:CASCADE;" `
	DeletedAt   gorm.DeletedAt
}

// Address Model
type Address struct {
	gorm.Model
	UserID   uint `gorm:"index;not null;foreignKey:UserID"`
	City     string
	State    string
	ZipCode  string
	Location string
}
