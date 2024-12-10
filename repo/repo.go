// repo to simulate the activity of database
package repo

import (
	"errors"
	"receipt-processor/models"
	"sync"
)

type ReceiptData struct {
	Receipt models.ExtReceipt
	Point   int64
}

var (
	mu sync.Mutex
	// Receipts is the main data storage
	// id -> ReceiptData
	Receipts    = make(map[string]ReceiptData)
	ErrNotFound = errors.New("receipt not found")
)

// Retrieves a ReceiptData by ID.
func GetReceiptData(id string) (ReceiptData, error) {
	mu.Lock()
	defer mu.Unlock()
	data, exists := Receipts[id]
	if !exists {
		return ReceiptData{}, ErrNotFound
	}
	return data, nil
}

// Updates or inserts a ReceiptData by ID.
func UpdateReceiptData(id string, data ReceiptData) {
	mu.Lock()
	defer mu.Unlock()
	Receipts[id] = data
}
