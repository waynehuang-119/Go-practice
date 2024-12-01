// main.go
// Application entry point.

package main

import (
	"fmt"
	receipt_handler "receipt-processor/public/v1/receipt"
	receiptSvc "receipt-processor/services/receipt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Create an instance of the ReceiptService
	receiptService := receiptSvc.NewReceiptService()

	// Set up routes
	receipt_handler.Register(router, receiptService)

	// Start the server
	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := router.Run(port); err != nil {
		panic(err)
	}
}
