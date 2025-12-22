package order

import "purchase-cart-service/models"

// CreateItem is the input DTO for creating orders
type CreateItem struct {
	ProductID string
	UnitPrice float64
	Quantity  int
}
type Detail struct {
	Id         string
	TotalPrice float64
	TotalVAT   float64
	Items      []ProductDetail
}
type ProductDetail struct {
	models.Product
	Quantity int
	VAT      float64
}
