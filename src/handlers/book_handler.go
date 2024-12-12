package handlers

import (
	"book-store/src/models"
	"encoding/json"
	"net/http"
)

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var body *models.Book

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Malformed request body"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
}
