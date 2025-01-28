package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"database"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "UP", DB: "OK"}

	// Ping the database
	sqlDB, err := initializers.DB.DB()
	if err != nil || sqlDB.Ping() != nil {
		response.DB = "DOWN"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string
		Password string
		Name     string
		Phone    string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to Read Body", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return

	}
	user := models.User{Email: body.Email, Password: string(hash), Name: body.Name, Phone: body.Phone}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Signup successful"})
}

func Login(w http.ResponseWriter, r *http.Request) {

	SECRET := os.Getenv("SECRET")

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

	tokenString, err := token.SignedString([]byte(SECRET))

	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
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
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Logged out"})
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
	b, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// email := r.FormValue("email")

	var body struct {
		Email string
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to Read Body", http.StatusBadRequest)
		return
	}
	// Fetch the user by email
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	// Generate reset token and set expiration
	token, _ := RandomToken()
	// if err != nil {
	// 	http.Error(w, "Unable to generate reset token", http.StatusInternalServerError)
	// 	return
	// }
	expires := time.Now().Add(1 * time.Hour)

	// Update the database with token and expiration
	user.ResetToken = token
	user.TokenExpiry = expires
	initializers.DB.Save(&user)

	// Simulate email by printing the reset link
	// fmt.Printf("Reset link: http://localhost:3000/reset-password?token=%s\n", token)
	link := "http://localhost:5173/reset-password?token=" + token
	SendEmail(user.Email, link)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password reset link sent"))
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse token and new password from the request body
	var requestData struct {
		Token       string `json:"token"`
		NewPassword string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the inputs
	if requestData.Token == "" || requestData.NewPassword == "" {
		http.Error(w, "Token and password are required", http.StatusBadRequest)
		return
	}

	// Fetch user by reset token
	var user models.User
	err := initializers.DB.First(&user, "reset_token = ?", requestData.Token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check if the token has expired
	if user.TokenExpiry.Before(time.Now()) {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		return
	}

	// Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), 10)

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusBadRequest)
		return

	}

	// Update the user's password and clear the reset token and expiry
	user.Password = string(hash)
	user.ResetToken = ""
	user.TokenExpiry = time.Time{} // Clear the expiry

	if err := initializers.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password successfully updated"))
}

func IsAdmin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Admin": true,
	})

}

func AdminStatus(w http.ResponseWriter, r *http.Request) {
	userid := r.PathValue("id")

	var body struct {
		Admin bool
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Failed to Read Body", http.StatusBadRequest)
		return
	}

	var user models.User
	initializers.DB.First(&user, userid)
	if body.Admin {

		user.Admin = true
	} else {
		user.Admin = false
	}
	initializers.DB.Save(&user)

	w.WriteHeader(http.StatusOK)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	// Use Preload to load the associated Cart for each User
	initializers.DB.Preload("Cart").Find(&users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var user models.User
	initializers.DB.First(&user, id)

	initializers.DB.Delete(&user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User deleted"})
}
