package receipt

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"receipt-processor/models"
	"receipt-processor/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockReceiptService is a mock implementation of the ReceiptService interface
type MockReceiptService struct {
	mock.Mock
}

func (m *MockReceiptService) ProcessReceipt(extReceipt models.ExtReceipt) (string, error) {
	args := m.Called(extReceipt)
	return args.String(0), args.Error(1)
}

func (m *MockReceiptService) GetPoints(id string) (int64, error) {
	args := m.Called(id)
	return args.Get(0).(int64), args.Error(1)
}

// ReceiptHandlerTestSuite defines the suite for handler tests
type ReceiptHandlerTestSuite struct {
	suite.Suite
	mockService    *MockReceiptService
	router         *gin.Engine
	mockExtReceipt models.ExtReceipt
}

// SetupTest initializes the suite
func (suite *ReceiptHandlerTestSuite) SetupTest() {
	// Reset the storage
	storage.Receipts = make(map[string]storage.ReceiptData)

	// Initialize the mock service
	suite.mockService = new(MockReceiptService)

	// Initialize the Gin router
	suite.router = gin.Default()

	// Register the routes with the mock service
	Register(suite.router, suite.mockService)

	// Define a mock external receipt (ExtReceipt)
	suite.mockExtReceipt = models.ExtReceipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
}

func (suite *ReceiptHandlerTestSuite) TestProcessReceipt() {
	// Set up mock expectations
	mockID := "mock-receipt-id"
	suite.mockService.On("ProcessReceipt", suite.mockExtReceipt).Return(mockID, nil)

	// Create a request
	req := httptest.NewRequest("POST", "/receipts/process", generateJSONBody(suite.mockExtReceipt))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	suite.router.ServeHTTP(w, req)

	// Assertions
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), mockID)
	suite.mockService.AssertCalled(suite.T(), "ProcessReceipt", suite.mockExtReceipt)
}

func (suite *ReceiptHandlerTestSuite) TestGetPoints() {
	// Set up mock expectations
	mockID := "mock-receipt-id"
	mockPoints := int64(100)
	suite.mockService.On("GetPoints", mockID).Return(mockPoints, nil)

	// Create a request
	req := httptest.NewRequest("GET", "/receipts/"+mockID+"/points", nil)
	w := httptest.NewRecorder()

	// Serve the request
	suite.router.ServeHTTP(w, req)

	// Assertions
	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), `"points":100`)
	suite.mockService.AssertCalled(suite.T(), "GetPoints", mockID)
}

// generateJSONBody creates an io.Reader containing the JSON body for testing
func generateJSONBody(extReceipt models.ExtReceipt) io.Reader {
	body, _ := json.Marshal(extReceipt)
	return bytes.NewReader(body) // Return an io.Reader
}

// Run the test suite
func TestReceiptHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptHandlerTestSuite))
}
