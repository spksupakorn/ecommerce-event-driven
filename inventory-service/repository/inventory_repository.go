package repository

import (
	"database/sql"
	"errors"

	"time"

	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) GetProduct(productID string) (*models.Product, error) {
	product := &models.Product{}

	query := `
		SELECT id, name, stock, reserved, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	err := r.db.QueryRow(query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Stock,
		&product.Reserved,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *InventoryRepository) ReserveStock(productID string, quantity int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the row for update
	var currentStock int
	query := `
		SELECT stock - reserved
		FROM products
		WHERE id = $1
		FOR UPDATE
	`

	err = tx.QueryRow(query, productID).Scan(&currentStock)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("product not found")
		}
		return err
	}

	// Check if enough stock is available
	if currentStock < quantity {
		return errors.New("insufficient stock")
	}

	// Update reserved stock
	updateQuery := `
		UPDATE products
		SET reserved = reserved + $1, updated_at = $2
		WHERE id = $3
	`

	_, err = tx.Exec(updateQuery, quantity, time.Now(), productID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *InventoryRepository) DeductStock(productID string, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1, reserved = reserved - $1, updated_at = $2
		WHERE id = $3 AND stock >= $1
	`

	result, err := r.db.Exec(query, quantity, time.Now(), productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to deduct stock")
	}

	return nil
}
