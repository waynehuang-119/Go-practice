package receipt

import (
	"errors"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/repo"
	receiptSvc "receipt-processor/services/receipt"

	"github.com/gin-gonic/gin"
)

var receiptService receiptSvc.ReceiptService

// Register router for the APIs
func Register(router *gin.Engine, service receiptSvc.ReceiptService) {
	receiptService = service

	// Define API routes
	router.POST("/receipts/process", ProcessReceipt)
	router.GET("/receipts/:id/points", GetPoints)

	// Custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
	})
}

func ProcessReceipt(c *gin.Context) {
	var extReceipt models.ExtReceipt

	// Parse JSON body
	if err := c.ShouldBindJSON(&extReceipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := receiptService.ProcessReceipt(extReceipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing receipt"})
		return
	}

	response := ExtProcessReceiptResponse{ID: id}
	c.JSON(http.StatusOK, response)
}

func GetPoints(c *gin.Context) {
	id := c.Param("id")

	points, err := receiptService.GetPoints(id)
	if err != nil {
		// Check if the error is a "not found" error
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
			return
		}
		// Handle all other errors as internal server errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Build and send the successful response
	response := ExtGetPointsResponse{Points: points}
	c.JSON(http.StatusOK, response)
}
