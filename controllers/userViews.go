package controllers

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to Read body", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return

	}
	user := models.User{Email: body.Email, Password: string(hash)}
	if initializers.DB.First(&user) != nil {
		Login(w, r)
		return
	}
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Failed to Create User", http.StatusBadRequest)
		return

	}

	cart := models.Cart{UserID: user.ID}
	result = initializers.DB.Create(&cart)
	if result.Error != nil {
		http.Error(w, "Failed to Create Cart", http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to Read body", http.StatusBadRequest)
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("l41^*&vjah4#%4565c4vty%#8b84"))

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString, // replace this with `tokenString`
		Path:     "/",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	// Set the cookie in the response header
	http.SetCookie(w, cookie)

	// Return an empty JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func Validate(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	initializers.DB.Preload("Cart").First(&user, user.ID)
	// Respond with the user information as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"User": user,
	})
}
