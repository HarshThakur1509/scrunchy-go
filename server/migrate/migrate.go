package main

import (
	"os"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
)

func init() {
	initializers.ConnectDB()
}

func main() {
	User := &models.User{}
	Address := &models.Address{}
	Product := &models.Product{}
	CartItem := &models.CartItem{}
	Cart := &models.Cart{}

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	initializers.DB.AutoMigrate(User, Address, Product, CartItem, Cart)
}
