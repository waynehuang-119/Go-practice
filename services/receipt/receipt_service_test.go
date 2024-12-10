package receipt

import (
	"receipt-processor/models"
	"receipt-processor/storage"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ReceiptServiceTestSuite defines the suite for service tests
type ReceiptServiceTestSuite struct {
	suite.Suite
	service        ReceiptService
	mockExtReceipt models.ExtReceipt
}

// SetupTest initializes the suite
func (suite *ReceiptServiceTestSuite) SetupTest() {
	// Reset the storage before each test
	storage.Receipts = make(map[string]storage.ReceiptData)

	// Initialize the ReceiptService
	suite.service = NewReceiptService()

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

func (suite *ReceiptServiceTestSuite) TestProcessReceipt() {
	// Process the mock receipt
	id, err := suite.service.ProcessReceipt(suite.mockExtReceipt)

	// Assertions
	suite.NoError(err)
	suite.NotEmpty(id)

	// Check if receipt exists in storage
	receiptData, exists := storage.GetReceiptData(id)
	suite.True(exists)
	suite.Equal(suite.mockExtReceipt.Retailer, receiptData.Receipt.Retailer)
	suite.Equal(suite.mockExtReceipt.Total, receiptData.Receipt.Total)
}

func (suite *ReceiptServiceTestSuite) TestGetPoints() {
	// Process the mock receipt and get its ID
	id, _ := suite.service.ProcessReceipt(suite.mockExtReceipt)

	// Get points for the processed receipt
	points, err := suite.service.GetPoints(id)

	// Assertions
	suite.NoError(err)
	suite.Greater(points, int64(0))
	suite.Equal(points, int64(28), "Points of this mock receipt should be 28")

	// Verify that points were updated in storage
	receiptData, _ := storage.GetReceiptData(id)
	suite.Equal(points, receiptData.Point)
}

// Run the test suite
func TestReceiptServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptServiceTestSuite))
}
