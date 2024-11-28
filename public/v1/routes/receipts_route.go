package routes

import (
	"go-practice/public/v1/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// Register router for the APIs
func Register(router *mux.Router) {
	// Skip cleaning the URL path (enabling empty {id} requests and return 404 instead of 301 redirect)
	// router.SkipClean(true)

	router.HandleFunc("/receipts/process", controllers.ProcessReceiptController).Methods(http.MethodPost)
	router.HandleFunc("/receipts/{id}/points", controllers.GetPointsController).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}
