package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Admin       bool
	ResetToken  string
	TokenExpiry time.Time
	Cart        Cart
}
