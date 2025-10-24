package models

import "time"

type Product struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Stock     int       `json:"stock" db:"stock"`
	Reserved  int       `json:"reserved" db:"reserved"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type StockReservation struct {
	ProductID string
	Quantity  int
	OrderID   string
}
