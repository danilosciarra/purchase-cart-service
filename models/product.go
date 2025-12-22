package models

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	VAT         float64
	CreatedAt   time.Time
}
