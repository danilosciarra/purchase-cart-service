package memory

import (
	"context"
	"purchase-cart-service/models"
)

type OrderRepository struct {
	orders map[string]*models.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (o *OrderRepository) Save(ctx context.Context, order *models.Order) error {
	return nil
}
func (o *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	return nil, nil
}
