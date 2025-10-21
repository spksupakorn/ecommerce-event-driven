package repository

import (
	"database/sql"

	"time"

	"github.com/google/uuid"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(req *models.CreateOrderRequest) (*models.Order, error) {
	order := &models.Order{
		ID:        uuid.New().String(),
		ItemID:    req.ItemID,
		Quantity:  req.Quantity,
		UserEmail: req.UserEmail,
		Status:    models.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO orders (id, item_id, quantity, user_email, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		order.ID,
		order.ItemID,
		order.Quantity,
		order.UserEmail,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetByID(id string) (*models.Order, error) {
	order := &models.Order{}

	query := `
		SELECT id, item_id, quantity, user_email, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.ItemID,
		&order.Quantity,
		&order.UserEmail,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) UpdateStatus(id, status string) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}
