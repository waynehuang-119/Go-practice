package receipt

import (
	"fmt"
	"math"
	"receipt-processor/models"
	"receipt-processor/repo/psql"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type ReceiptService interface {
	ProcessReceipt(extReceipt models.ExtReceipt) (string, error)
	GetPoints(id string) (int64, error)
}

type receiptServiceImpl struct {
	repo psql.ReceiptRepository
}

func NewReceiptService(repo psql.ReceiptRepository) ReceiptService {
	return &receiptServiceImpl{repo: repo}
}

// Stores a receipt, generates an ID, process points and returns the ID
func (r *receiptServiceImpl) ProcessReceipt(receipt models.ExtReceipt) (string, error) {
	// Generate unique ID
	id := uuid.New().String()

	// Calculate points
	points := calculatePoints(receipt)

	// update receipt to db
	err := r.repo.UploadReceipt(receipt, id, points)
	if err != nil {
		return "", fmt.Errorf("repo layer - failed to upload receipt: %v", err)
	}

	return id, nil
}

// Get points for a given receipt ID
func (r *receiptServiceImpl) GetPoints(id string) (int64, error) {
	points, err := r.repo.GetPoints(id)

	if err != nil {
		return 0, err
	}

	return points, nil
}

// Calculates points for a given receipt
func calculatePoints(receipt models.ExtReceipt) int64 {
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
