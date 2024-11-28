package services

import (
	"fmt"
	"math"
	"receipt-processor/models"
	"receipt-processor/storage"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Stores a receipt, generates an ID, and returns the ID.
func ProcessReceipt(extReceipt models.ExtReceipt) (string, error) {
	var id string

	// Generate unique ID
	for {
		id = uuid.New().String()

		// Check if the ID exists in storage
		if _, exists := storage.Receipts[id]; !exists {
			break // ID is unique, exit the loop
		}
	}

	// Convert external receipt to internal receipt
	var items []models.Item
	for _, extItem := range extReceipt.Items {
		items = append(items, models.Item{
			ShortDescription: extItem.ShortDescription,
			Price:            extItem.Price,
		})
	}

	internalReceipt := models.Receipt{
		ID:           id,
		Retailer:     extReceipt.Retailer,
		PurchaseDate: extReceipt.PurchaseDate,
		PurchaseTime: extReceipt.PurchaseTime,
		Items:        items,
		Total:        extReceipt.Total,
	}

	receiptData := storage.ReceiptData{Receipt: internalReceipt, Point: 0}
	storage.UpdateReceiptData(id, receiptData)
	return id, nil
}

// Get points for a given receipt ID, calculating them if 0.
func GetPoints(id string) (int64, error) {
	receiptData, exists := storage.GetReceiptData(id)
	if !exists {
		return 0, fmt.Errorf("receipt with id %s does not exist", id)
	}

	// If points are 0, calculate and update them
	if receiptData.Point == 0 {
		receiptData.Point = calculatePoints(receiptData.Receipt)
		storage.UpdateReceiptData(id, receiptData)
	}

	return receiptData.Point, nil
}

// Calculates points for a given receipt.
func calculatePoints(receipt models.Receipt) int64 {
	var points int64

	// Rule 1: One point for every alphanumeric character in the retailer name.
	retailerRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += int64(len(retailerRegex.FindAllString(receipt.Retailer, -1)))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return points
	}
	if total == float64(int(total)) {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += int64(len(receipt.Items) / 2 * 5)

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int64(math.Ceil(price * 0.2))
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	day, err := strconv.Atoi(strings.Split(receipt.PurchaseDate, "-")[2])
	if err == nil && day%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour, err := strconv.Atoi(strings.Split(receipt.PurchaseTime, ":")[0])
	if err == nil && hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}
