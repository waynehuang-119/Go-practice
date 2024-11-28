// storage to simulate the activity of database
package storage

import (
	"receipt-processor/models"
	"sync"
)

type ReceiptData struct {
	Receipt models.Receipt
	Point   int64
}

var (
	mu sync.Mutex
	// Receipts is the main data storage
	// id -> ReceiptData
	Receipts = make(map[string]ReceiptData)
)

// Retrieves a ReceiptData by ID.
func GetReceiptData(id string) (ReceiptData, bool) {
	mu.Lock()
	defer mu.Unlock()
	data, exists := Receipts[id]
	return data, exists
}

// Updates or inserts a ReceiptData by ID.
func UpdateReceiptData(id string, data ReceiptData) {
	mu.Lock()
	defer mu.Unlock()
	Receipts[id] = data
}
