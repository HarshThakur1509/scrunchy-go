package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"github.com/markbates/goth/gothic"
)

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	// Finalize the authentication process
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		log.Println(err)
		return
	}

	// Save user to the database
	userModel := models.User{
		Name:  user.Name,
		Email: user.Email,
	}

	result := initializers.DB.FirstOrCreate(&userModel, "email = ?", userModel.Email)
	if result.Error != nil {
		http.Error(w, "Failed to Create User", http.StatusBadRequest)
		return

	}

	// Save user ID in the session
	var id string = strconv.FormatUint(uint64(userModel.ID), 10)
	err = gothic.StoreInSession("user_id", id, r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Redirect to the secure area
	redirectSecure := os.Getenv("REDIRECT_SECURE")
	if redirectSecure == "" {
		redirectSecure = "https://scrunchy.harshthakur.site/"
	}

	http.Redirect(w, r, redirectSecure, http.StatusFound)
}

func GothLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session
	err := gothic.Logout(w, r)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from the session
	userID, err := gothic.GetFromSession("user_id", r)
	if err != nil || userID == "" {
		// Return an empty JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{"exists": false})
		return
	}

	// Return an empty JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{"exists": true, "userID": userID})
}
