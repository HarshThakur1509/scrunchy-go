package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/HarshThakur1509/scrunchy-go/initializers"
	"github.com/HarshThakur1509/scrunchy-go/models"
)

func PostProduct(w http.ResponseWriter, r *http.Request) {

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Extract form fields
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	// Extract the image
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the image locally
	imagePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(imagePath)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}

	// Save to database
	product := models.Product{Name: name, Price: uint(price), Image: imagePath}
	result := initializers.DB.Create(&product)
	if result.Error != nil {
		http.Error(w, "Unable to save data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Data saved successfully"})

}

func GetProducts(w http.ResponseWriter, r *http.Request) {

	var products []models.Product

	// Fetch all products from the database
	result := initializers.DB.Find(&products)
	if result.Error != nil {
		http.Error(w, "Unable to fetch data", http.StatusInternalServerError)
		return
	}

	// Respond with JSON data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetProductIndex(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var product models.Product
	initializers.DB.First(&product, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}
