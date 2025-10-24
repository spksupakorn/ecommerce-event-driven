package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Inventory database connected successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		stock INTEGER NOT NULL DEFAULT 0,
		reserved INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_products_id ON products(id);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Insert sample products
	insertSampleData(db)

	log.Println("Products table created successfully")
	return nil
}

func insertSampleData(db *sql.DB) {
	query := `
	INSERT INTO products (id, name, stock, reserved)
	VALUES 
		('product-001', 'Laptop', 100, 0),
		('product-002', 'Mouse', 500, 0),
		('product-003', 'Keyboard', 300, 0)
	ON CONFLICT (id) DO NOTHING;
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Failed to insert sample data: %v", err)
	} else {
		log.Println("Sample products inserted")
	}
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}
