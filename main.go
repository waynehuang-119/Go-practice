// main.go
// Application entry point.

package main

import (
	"fmt"
	"go-practice/public/v1/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Set up routes
	routes.Register(router)

	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
