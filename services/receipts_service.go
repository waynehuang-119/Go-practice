package services

import (
	"fmt"
	model "go-practice/models"
	"go-practice/storage"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Stores a receipt, generates an ID, and returns the ID.
func ProcessReceipt(extReceipt model.ExtReceipt) (string, error) {
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
	var items []model.Item
	for _, extItem := range extReceipt.Items {
		items = append(items, model.Item{
			ShortDescription: extItem.ShortDescription,
			Price:            extItem.Price,
		})
	}

	internalReceipt := model.Receipt{
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
func calculatePoints(receipt model.Receipt) int64 {
	var points int64

	// Rule 1: 1 point for every alphanumeric character in the retailer name.
	retailerRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += int64(len(retailerRegex.FindAllString(receipt.Retailer, -1)))

	// Rule 2: 50 points if the total is a round dollar amount.
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		// Handle error, e.g., return 0 points if total is invalid
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

	// Rule 5: If the trimmed description of an item has a length divisible by 3, multiply the price by 0.2 and round up.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int64((price * 0.2) + 0.5) // Round up
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	day, err := strconv.Atoi(strings.Split(receipt.PurchaseDate, "-")[2])
	if err == nil && day%2 != 0 {
		points += 6
	}

	// Rule 7: 10 points if the purchase time is between 2:00 PM and 4:00 PM.
	hour, err := strconv.Atoi(strings.Split(receipt.PurchaseTime, ":")[0])
	if err == nil && hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}
