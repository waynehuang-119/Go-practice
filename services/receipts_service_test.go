package services

import (
	model "go-practice/models"
	"go-practice/storage"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ReceiptServiceTestSuite struct {
	suite.Suite
	mockExtReceipt model.ExtReceipt // Change from model.Receipt to model.ExtReceipt
}

func (suite *ReceiptServiceTestSuite) SetupTest() {
	// Reset the storage before each test
	storage.Receipts = make(map[string]storage.ReceiptData)

	// Define a mock external receipt (ExtReceipt)
	suite.mockExtReceipt = model.ExtReceipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []model.Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
		Total: "35.35",
	}
}

// Test ProcessReceipt
func (suite *ReceiptServiceTestSuite) TestProcessReceipt() {
	id, _ := ProcessReceipt(suite.mockExtReceipt) // Pass ExtReceipt here

	// Verify that the receipt is stored with the generated ID
	storedData, exists := storage.GetReceiptData(id)
	suite.Require().True(exists, "Receipt should exist in storage")
	// Compare stored data's receipt with the mock data (after conversion to internal format)
	suite.Require().Equal(suite.mockExtReceipt.Retailer, storedData.Receipt.Retailer, "Stored receipt retailer should match mock")
	suite.Require().Equal(suite.mockExtReceipt.PurchaseDate, storedData.Receipt.PurchaseDate, "Stored receipt date should match mock")
	suite.Require().Equal(int64(0), storedData.Point, "Points should initially be 0")
}

// Test get points and point calculation
func (suite *ReceiptServiceTestSuite) TestGetPoints() {
	id, _ := ProcessReceipt(suite.mockExtReceipt) // Pass ExtReceipt here
	points, err := GetPoints(id)

	// Verify points calculation
	suite.Require().NoError(err, "GetPoints should not return an error")
	suite.Require().Greater(points, int64(0), "Points should be greater than 0")
	// suite.Require().Equal(points, int64(28), "Points of this mock receipt should be 28")

	// Verify that points are updated in storage
	storedData, exists := storage.GetReceiptData(id)
	suite.Require().True(exists, "Receipt should exist in storage")
	suite.Require().Equal(points, storedData.Point, "Stored points should match calculated points")
}

// TestGetPointsInvalidID verifies behavior with an invalid ID
func (suite *ReceiptServiceTestSuite) TestGetPointsInvalidID() {
	_, err := GetPoints(uuid.New().String())
	suite.Require().Error(err, "GetPoints should return an error for an invalid ID")
}

// Run the test suite
func TestReceiptServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptServiceTestSuite))
}
