package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"receipt-processor/models"
	"receipt-processor/public/v1/receipt"

	"github.com/stretchr/testify/suite"
)

type SmokeTestSuite struct {
	suite.Suite
	baseURL   string
	receipt   models.ExtReceipt
	receiptID string
}

func (suite *SmokeTestSuite) SetupTest() {
	suite.baseURL = "http://localhost:8080"
}

func (suite *SmokeTestSuite) SetupSuite() {
	suite.receipt = models.ExtReceipt{
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

// Name the function as TestA... to make sure processReceipt run first
func (suite *SmokeTestSuite) TestAProcessReceipt() {
	suite.T().Log("Running TestProcessReceipt...")

	data, err := json.Marshal(suite.receipt)
	suite.Require().NoError(err, "Failed to marshal receipt")

	res, err := http.Post(suite.baseURL+"/receipts/process", "application/json", bytes.NewReader(data))
	suite.Require().NoError(err, "Failed to send POST request")
	defer res.Body.Close()

	suite.Require().Equal(http.StatusOK, res.StatusCode, "Unexpected status code")

	body, err := io.ReadAll(res.Body)
	suite.Require().NoError(err, "Failed to read response body")

	var response receipt.ExtProcessReceiptResponse
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err, "Failed to unmarshal response")

	suite.T().Logf("Response body: %s", string(body))

	suite.receiptID = response.ID
	suite.Require().NotEmpty(suite.receiptID, "Receipt ID should not be empty")

	suite.T().Logf("Processed receipt ID: %s", suite.receiptID)
}

// Name the function as TestB... to make sure getPoints run after
func (suite *SmokeTestSuite) TestBGetPoints() {
	suite.T().Log("Running TestGetPoints")

	suite.T().Logf("Geting point from receipt ID: %s", suite.receiptID)
	res, err := http.Get(suite.baseURL + "/receipts/" + suite.receiptID + "/points")
	suite.Require().NoError(err, "Failed to send GET request")
	defer res.Body.Close()

	suite.Require().Equal(http.StatusOK, res.StatusCode, "Unexpected status code")

	body, err := io.ReadAll(res.Body)
	suite.Require().NoError(err, "Failed to read response body")

	var response receipt.ExtGetPointsResponse
	err = json.Unmarshal(body, &response)
	suite.Require().NoError(err, "Failed to unmarshal response")

	suite.Require().Greater(response.Points, int64(0), "Points should be greater than zero")

	suite.T().Logf("Retrieved points: %d", response.Points)
}

func TestSmokeTestSuite(t *testing.T) {
	suite.Run(t, new(SmokeTestSuite))
}
