package models

type Order struct {
	ID         string
	Items      []Item
	TotalPrice int64
	TotalVAT   int64
}

type Item struct {
	ProductID string
	Quantity  int
}
