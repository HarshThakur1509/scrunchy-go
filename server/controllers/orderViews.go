package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"github.com/markbates/goth/gothic"
	"github.com/razorpay/razorpay-go"
)

func Pay(w http.ResponseWriter, r *http.Request) {
	PAY_ID := os.Getenv("PAY_ID")
	PAY_SECRET := os.Getenv("PAY_SECRET")

	client := razorpay.NewClient(PAY_ID, PAY_SECRET)

	// Retrieve user ID from the session
	userID, err := gothic.GetFromSession("user_id", r)
	if err != nil || userID == "" {
		// Return an empty JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{"exists": false})
		return
	}

	var user models.User

	initializers.DB.Preload("Cart").Omit("password").First(&user, userID)

	data := map[string]interface{}{
		"amount":   user.Cart.Total * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	w.Header().Set("Content-Type", "application/json")
	body, err := client.Order.Create(data, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func PayResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)

}
