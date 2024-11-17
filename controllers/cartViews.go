package controllers

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// func AddToCart(w http.ResponseWriter, r *http.Request) {
// 	var body struct {
// 		ProductID uint
// 	}

// 	user, ok := r.Context().Value("user").(models.User)
// 	if !ok {
// 		http.Error(w, "User not found", http.StatusUnauthorized)
// 		return
// 	}
// 	var cart models.Cart
// 	initializers.DB.Model(models.Cart{UserID: user.ID}).First(&cart)

// 	json.NewDecoder(r.Body).Decode(&body)

// 	cartItem := models.CartItem{
// 		Quantity:  1,
// 		ProductID: body.ProductID,
// 		CartID:    cart.ID,
// 	}

// 	if len(cart.CartItems) > 0 {
// 		var newItems []models.CartItem
// 		initializers.DB.Model(&cart).Association("CartItems").Find(&newItems)
// 		for i := 0; i < len(newItems); i++ {
// 			if newItems[i].ProductID == body.ProductID {
// 				newItems[i].Quantity += 1
// 				break
// 			}
// 		}
// 		initializers.DB.Model(&cart).Association("CartItems").Clear()
// 		initializers.DB.Model(&cart).Association("CartItems").Append(&newItems)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(cart)

// 	} else {
// 		initializers.DB.Create(&cartItem)
// 		initializers.DB.Model(&cart).Association("CartItems").Append(cartItem)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(cart)
// 	}

// }

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
