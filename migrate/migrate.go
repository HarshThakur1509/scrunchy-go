package main

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"os"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	User := &models.User{}
	Product := &models.Product{}
	CartItem := &models.CartItem{}
	Cart := &models.Cart{}

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	initializers.DB.AutoMigrate(User, Product, CartItem, Cart)
}
