package models

import "time"

type Order struct {
	ID         string
	Items      []Item
	TotalPrice float64
	TotalVAT   float64
	CreatedAt  time.Time
}

type Item struct {
	ProductID string
	Name      string
	Quantity  int
	UnitPrice float64
	VAT       float64
}
