package psql

import (
	"database/sql"
	"fmt"
	"receipt-processor/models"
)

var ErrNoRows = fmt.Errorf("no receipt found with provided receipt ID")

type ReceiptRepo interface {
	UploadReceipt(receipt models.ExtReceipt, id string, points int64) error
	GetPoints(id string) (int64, error)
}

type receiptRepo struct {
	db *sql.DB
}

func New(db *sql.DB) ReceiptRepo {
	return &receiptRepo{db: db}
}

// upload receipt query
func (r receiptRepo) UploadReceipt(receipt models.ExtReceipt, id string, points int64) error {
	// use transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction when uploading receipt: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// insert into receipt table
	query := `INSERT INTO receipt (id, retailer, purchaseDate, purchaseTime, total, points)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.Exec(query, id, receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, points)
	if err != nil {
		return fmt.Errorf("failed to insert receipt data into receipt table: %v", err)
	}

	// insert into item table
	for description, price := range receipt.Items {
		var itemId int
		query = `INSERT INTO item (short_description, price)
		VALUES ($1, $2)
		RETURNING id
		`

		err := tx.QueryRow(query, description, price).Scan(&itemId)

		if err != nil {
			return fmt.Errorf("failed to insert items into item table: %v", err)
		}

		// insert into receipt_item table
		query = `INSERT INTO receipt_item (receipt_id, item_id)
			VALUES ($1, $2)
		`

		_, err = tx.Exec(query, id, itemId)
		if err != nil {
			return fmt.Errorf("failed to link receipt and item: %v", err)
		}
	}

	return nil
}

// get points query
func (r receiptRepo) GetPoints(id string) (int64, error) {
	var points int64
	query := `SELECT points FROM receipt WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&points)
	if err == sql.ErrNoRows {
		return 0, ErrNoRows
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get points from receipt table: %v", err)
	}

	return points, nil
}
