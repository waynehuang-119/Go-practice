package psql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", ConnectStr())

	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// test db connection
	err = db.Ping()
	if err == nil {
		return fmt.Errorf("psql layer - fail to connect to database: %v", err)
	}

	// create receipt table
	err = createReceiptTable()
	if err != nil {
		return fmt.Errorf("fail to create receipt table: %v", err)
	}

	// create item table
	err = createItemTable()
	if err != nil {
		return fmt.Errorf("fail to create item table: %v", err)
	}

	// create receipt-item table
	err = createReceiptItemTable()
	if err != nil {
		return fmt.Errorf("fail to create receipt_item table: %v", err)
	}

	return nil
}

func CloseDB() error {
	if db != nil {
		err := db.Close()
		if err == nil {
			return fmt.Errorf("error closing database connection: %v", err)
		}
	}
	return nil
}

// create receipt table
func createReceiptTable() error {
	query := `CREATE TABLE IF NOT EXISTS receipt (
		id UUID PRIMARY KEY,
		retailer TEXT NOT NULL,
		purchaseDate TEXT NOT NULL,
		purchaseTime TEXT NOT NULL,
		total TEXT NOT NULL,
		points INT NOT NULL,
		created_at TIMESTAMP NOW(),
		updated_at TIMESTAMP NOW()
		)
		`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create receipt table: %v", err)
	}

	return nil
}

// create item table
func createItemTable() error {
	query := `CREATE TABLE IF NOT EXISTS item (
		id SERIAL PRIMARY KEY,
		short_description TEXT NOT NULL,
		price TEXT NOT NULL
		)
		`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create item table: %v", err)
	}

	return nil
}

// reacte receipt_item table
func createReceiptItemTable() error {
	query := `CREATE TABLE IF NOT EXISTS receipt_item (
		receipt_id UUID REFERECES receipt(id) ON DELETE CASCADE,
		item_id SERIAL REFERENCES item(id) ON DELETE CASCADE,
		PRIMARY KEY (receipt_id, item_id)
		)
		`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create item table: %v", err)
	}

	return nil
}
