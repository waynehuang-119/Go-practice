package psql

import (
	"database/sql"
	"fmt"
	"receipt-processor/models"
)

var IdNotFound = fmt.Errorf("no receipt found with provided ID")

type ReceiptRepository interface {
	UploadReceipt(receipt models.ExtReceipt, receiptId string, points int64) error
	GetPoints(receiptId string) (int64, error)
}

type receiptRepository struct {
	db *sql.DB
}

func New(db *sql.DB) ReceiptRepository {
	return &receiptRepository{db: db}
}

func (r *receiptRepository) UploadReceipt(receipt models.ExtReceipt, receiptId string, points int64) error {
	// start the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction when uploading receipt: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// insert receipt
	query := `INSERT INTO receipt (id, retailer, purchaseDate, purchaseTime, total, points)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.Exec(query, receiptId, receipt.Retailer, receipt.PurchaseDate, receipt.PurchaseTime, receipt.Total, points)
	if err != nil {
		return fmt.Errorf("repo layer - failed to insert receipt into receipt table: %v", err)
	}

	// insert items
	for _, item := range receipt.Items {
		var itemId int
		query := `INSERT INTO item (shortDescription, price)
			VALUES ($1, $2)
			RETURNING id`

		err := tx.QueryRow(query, item.ShortDescription, item.Price).Scan(&itemId)

		if err != nil {
			return fmt.Errorf("repo layer - failed to insert item into item table: %v", err)
		}

		// insert receipt_item
		query = `INSERT INTO receipt_item (receipt_id, item_id)
			VALUES ($1, $2)`

		_, err = tx.Exec(query, receiptId, itemId)
		if err != nil {
			return fmt.Errorf("repo layer - failed to link receipt and item: %v", err)
		}
	}

	return nil

}

func (r *receiptRepository) GetPoints(id string) (int64, error) {
	var points int64

	query := `SELECT points FROM receipt WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&points)

	if err == sql.ErrNoRows {
		return 0, IdNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("error retrieving receipt points: %v", err)
	}

	return points, nil
}
