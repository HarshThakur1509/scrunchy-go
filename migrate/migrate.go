package main

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	User := &models.User{}
	Product := &models.Product{}
	CartItem := &models.CartItem{}
	Cart := &models.Cart{}

	initializers.DB.AutoMigrate(User, Product, CartItem, Cart)
}
