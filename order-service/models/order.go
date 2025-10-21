package models

import (
	"time"
)

type Order struct {
	ID        string    `json:"id" db:"id"`
	ItemID    string    `json:"item_id" db:"item_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	UserEmail string    `json:"user_email" db:"user_email"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateOrderRequest struct {
	ItemID    string `json:"item_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	UserEmail string `json:"user_email" binding:"required,email"`
}

const (
	OrderStatusPending   = "PENDING"
	OrderStatusProcessed = "PROCESSED"
	OrderStatusFailed    = "FAILED"
)
