// main.go
// Application entry point.

package main

import (
	"fmt"
	"log"
	_ "receipt-processor/docs"
	"receipt-processor/psql"
	receipt_handler "receipt-processor/public/v1/receipt"
	receiptRepo "receipt-processor/repo/psql"
	receiptSvc "receipt-processor/services/receipt"

	"github.com/gin-gonic/gin"
)

// @title Receipt Processor API
// @version 1.0
// @description This is a backend service written in Go using Gin framework which processes receipt awards points.

// @host localhost:8080/
func main() {
	// Initialize the database
	err := psql.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer psql.CloseDB()

	// Create a Gin router
	router := gin.Default()
	router.Use(errorHandler())

	// Initialize the ReceiptRepository
	receiptRepository := receiptRepo.New(psql.GetDB())

	// Create an instance of the ReceiptService
	receiptService := receiptSvc.NewReceiptService(receiptRepository)

	// Set up routes
	receipt_handler.Register(router, receiptService)

	// Start the server
	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	if err := router.Run(port); err != nil {
		panic(err)
	}
}

func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during the request
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %v", e.Err) // Log error details
			}
		}
	}
}
