// main.go
// Application entry point.

package main

import (
	"fmt"
	_ "receipt-processor/docs"
	receipt_handler "receipt-processor/public/v1/receipt"
	receiptSvc "receipt-processor/services/receipt"

	"github.com/gin-gonic/gin"
)

// @title Receipt Processor API
// @version 1.0
// @description This is a backend service written in Go using Gin framework which processes receipt awards points.

// @host localhost:8080/
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
