package psql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Open DB connection
func InitDB() error {
	var err error
	db, err = sql.Open("postgres", GetPsqlConfig())

	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create Receipt table
	if err = createReceiptTable(db); err != nil {
		// Close the connection if table creation fails
		db.Close()
		return fmt.Errorf("error creating Receipt table: %v", err)
	}

	// Create Item table
	if err = createItemTable(db); err != nil {
		db.Close()
		return fmt.Errorf("error creating Item table: %v", err)
	}

	// Create ReceiptItem table
	if err = createReceiptItemTable(db); err != nil {
		db.Close()
		return fmt.Errorf("error creating Receipt_Item table: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL and created Tables!")
	return nil
}

// Close the database connection
func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

func createReceiptTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Receipt (
		id UUID PRIMARY KEY,
		retailer TEXT NOT NULL,
		purchaseDate TEXT NOT NULL,
		purchaseTime TEXT NOT NULL,
		total TEXT NOT NULL,
		point INT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
    	updated_at TIMESTAMP DEFAULT NOW()
		)
	`
	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("error creating Receipt table: %v", err)
	}
	return nil
}

func createItemTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Item (
		id SERIAL PRIMARY KEY,
		shortDescription TEXT,
		price TEXT
		)
	`

	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("error creating Item table: %v", err)
	}
	return nil
}

func createReceiptItemTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS Receipt_Item (
		receipt_id UUID REFERENCES Receipt(id) ON DELETE CASCADE,
    	item_id SERIAL REFERENCES Item(id) ON DELETE CASCADE,
    	PRIMARY KEY (receipt_id, item_id)
		)
	`

	_, err := db.Exec(query)

	if err != nil {
		return fmt.Errorf("error creating Receipt_Item table: %v", err)
	}

	return nil
}
