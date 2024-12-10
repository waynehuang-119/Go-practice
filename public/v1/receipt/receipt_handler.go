package receipt

import (
	"errors"
	"net/http"
	"receipt-processor/models"
	repo "receipt-processor/repo/psql"
	receiptSvc "receipt-processor/services/receipt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var receiptService receiptSvc.ReceiptService

type ErrorResponse struct {
	Error string `json:"error"`
}

// Register router for the APIs
func Register(router *gin.Engine, service receiptSvc.ReceiptService) {
	receiptService = service
	router.Use(cors.Default())

	// Swagger for API docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define API routes
	router.POST("/receipts/process", ProcessReceipt)
	router.GET("/receipts/:id/points", GetPoints)

	// Custom 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not Found"})
	})
}

// ProcessReceipt godoc
// @Summary Submits a receipt for processing and returns an ID
// @Description Receives a receipt in JSON format and processes it, returning a unique ID for the receipt.
// @Tags receipts
// @Accept json
// @Produce json
// @Param receipt body models.ExtReceipt true "Receipt data"
// @Success 200 {object} ExtProcessReceiptResponse "Receipt processed successfully"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Error processing receipt"
// @Router /receipts/process [post]
func ProcessReceipt(c *gin.Context) {
	var extReceipt models.ExtReceipt

	// Parse JSON body
	if err := c.ShouldBindJSON(&extReceipt); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	id, err := receiptService.ProcessReceipt(extReceipt)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error processing receipt"})
		return
	}

	response := ExtProcessReceiptResponse{ID: id}
	c.JSON(http.StatusOK, response)
}

// GetPoints godoc
// @Summary Retrieves points associated with a receipt by ID
// @Description Fetches the points linked to a receipt using its unique ID.
// @Tags receipts
// @Accept json
// @Produce json
// @Param id path string true "Receipt ID"
// @Success 200 {object} ExtGetPointsResponse "Points retrieved successfully"
// @Failure 404 {object} ErrorResponse "Receipt not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /receipts/{id}/points [get]
func GetPoints(c *gin.Context) {
	id := c.Param("id")

	points, err := receiptService.GetPoints(id)
	if err != nil {
		// Check if the error is a "not found" error
		if errors.Is(err, repo.IdNotFound) {
			c.Error(err)
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Receipt not found"})
			return
		}
		// Handle all other errors as internal server errors
		c.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	// Build and send the successful response
	response := ExtGetPointsResponse{Points: points}
	c.JSON(http.StatusOK, response)
}
