package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"github.com/razorpay/razorpay-go"
)

func Pay(w http.ResponseWriter, r *http.Request) {
	client := razorpay.NewClient("rzp_test_MQUwdShJLMIpOu", "SX6q5kYO9VACTkoidnX7RQev")

	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	initializers.DB.Preload("Cart").First(&user, user.ID)

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
