package repository

import (
	"context"
	"purchase-cart-service/models"
	"purchase-cart-service/repository/memory"
)

type OrderRepository interface {
	Save(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	GetAll(ctx context.Context) ([]*models.Order, error)
}

func NewOrderRepository(repoType string) OrderRepository {
	var repoOrder OrderRepository
	switch repoType {
	case "InMemory":
		repoOrder = memory.NewOrderRepository()
	}
	return repoOrder
}
