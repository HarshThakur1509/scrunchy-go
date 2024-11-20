package controllers

import (
	"Scrunchy/initializers"
	"Scrunchy/models"
	"encoding/json"
	"net/http"
)

func PostProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string
		Price int
	}
	json.NewDecoder(r.Body).Decode(&body)
	product := models.Product{Name: body.Name, Price: body.Price}
	result := initializers.DB.Create(&product)

	if result.Error != nil {
		http.Error(w, "Something went wrong!!", http.StatusBadRequest)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	initializers.DB.Find(&products)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
