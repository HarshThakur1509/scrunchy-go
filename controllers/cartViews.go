package controllers

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"encoding/json"
	"errors"
	"net/http"

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
	var cartItems []models.CartItem
	initializers.DB.Preload("Product").Find(&cartItems, "cart_id =?", cart.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cartItems)

}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProductID uint
	}

	// Get user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Ensure the cart exists
	var cart models.Cart
	initializers.DB.FirstOrCreate(&cart, models.Cart{UserID: user.ID})

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the product already exists in the cart
	var cartItem models.CartItem
	err := initializers.DB.Where("cart_id = ? AND product_id = ?", cart.ID, body.ProductID).First(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Product not in cart, create a new cart item
			cartItem = models.CartItem{
				Quantity:  1,
				ProductID: body.ProductID,
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
	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProductID uint
	}

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
	initializers.DB.Delete(&cartItem, "product_id = ? AND cart_id = ?", body.ProductID, cart.ID)

	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)

}

func QuantityCart(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProductID uint
		Quantity  uint
	}
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
	initializers.DB.Find(&cartItem, "product_id = ? AND cart_id = ?", body.ProductID, cart.ID)

	cartItem.Quantity = body.Quantity
	initializers.DB.Save(&cartItem)

	initializers.DB.Preload("CartItems.Product").First(&cart, cart.ID)

	// Respond with the updated cart
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}
