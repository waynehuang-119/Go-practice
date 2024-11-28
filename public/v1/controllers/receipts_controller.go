package controllers

import (
	"encoding/json"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/services"

	"github.com/gorilla/mux"
)

func ProcessReceiptController(w http.ResponseWriter, r *http.Request) {
	var extReceipt models.ExtReceipt

	err := json.NewDecoder(r.Body).Decode(&extReceipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := services.ProcessReceipt(extReceipt)
	if err != nil {
		http.Error(w, "Error processing receipt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ExtProcessReceiptResponse{ID: id}
	json.NewEncoder(w).Encode(response)
}

func GetPointsController(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	points, err := services.GetPoints(id)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ExtGetPointsResponse{Points: points}
	json.NewEncoder(w).Encode(response)
}
