package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"gorm.io/gorm"
)

func ListCart(w http.ResponseWriter, r *http.Request) {
	// Get user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Ensure the cart exists
	var cart models.Cart
	initializers.DB.FirstOrCreate(&cart, models.Cart{UserID: user.ID})

	// Fetch and return the cart items
	// var cartItems []models.CartItem
	initializers.DB.Preload("CartItems.Product").Find(&cart, cart.ID)
	// initializers.DB.Preload("CartItems.Product").Find(&cartItems, "cart_id =?", cart.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)

}

func AddToCart(w http.ResponseWriter, r *http.Request) {

	productidStr := r.PathValue("id")
	productid, _ := strconv.ParseUint(productidStr, 10, 64)
	// Get user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Ensure the cart exists
	var cart models.Cart
	initializers.DB.FirstOrCreate(&cart, models.Cart{UserID: user.ID})

	// Check if the product already exists in the cart
	var cartItem models.CartItem
	err := initializers.DB.Where("cart_id = ? AND product_id = ?", cart.ID, uint(productid)).First(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Product not in cart, create a new cart item
			cartItem = models.CartItem{
				Quantity:  1,
				ProductID: uint(productid),
				CartID:    cart.ID,
			}
			initializers.DB.Create(&cartItem)

		} else {
			// Handle other errors
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		// Product already in cart, update quantity
		cartItem.Quantity++
		initializers.DB.Save(&cartItem)

	}

	// Reload the cart with updated items and preload CartItems and Product
	initializers.DB.Preload("Product").First(&cartItem, "cart_id = ?", cart.ID)
	cart.Total += cartItem.Product.Price
	initializers.DB.Save(&cart)

	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	productidStr := r.PathValue("id")
	productid, _ := strconv.ParseUint(productidStr, 10, 64)

	// Get user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var cart models.Cart
	initializers.DB.First(&cart, models.Cart{UserID: user.ID})

	var cartItem models.CartItem
	initializers.DB.Find(&cartItem, "product_id = ? AND cart_id = ?", uint(productid), cart.ID)

	initializers.DB.Preload("Product").First(&cartItem, "cart_id = ?", cart.ID)
	cart.Total -= cartItem.Product.Price * cartItem.Quantity
	initializers.DB.Delete(&cartItem)
	initializers.DB.Save(&cart)

	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)

}

func QuantityCart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Quantity uint
	}
	productidStr := r.PathValue("id")
	productid, _ := strconv.ParseUint(productidStr, 10, 64)
	// Get user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	var cart models.Cart
	initializers.DB.First(&cart, models.Cart{UserID: user.ID})

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var cartItem models.CartItem
	initializers.DB.Find(&cartItem, "product_id = ? AND cart_id = ?", uint(productid), cart.ID)

	initializers.DB.Preload("Product").First(&cartItem, "cart_id = ?", cart.ID)

	if body.Quantity == 0 {
		cart.Total -= cartItem.Product.Price * cartItem.Quantity
		initializers.DB.Delete(&cartItem)
		initializers.DB.Save(&cart)
	} else {
		cart.Total += cartItem.Product.Price * (body.Quantity - cartItem.Quantity)
		cartItem.Quantity = body.Quantity
		initializers.DB.Save(&cartItem)
		initializers.DB.Save(&cart)

	}

	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}
