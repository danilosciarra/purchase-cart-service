package repository

import (
	"context"
	"purchase-cart-service/models"
	"purchase-cart-service/repository/memory"
)

type ProductRepository interface {
	GetProduct(ctx context.Context, id string) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
}

func NewProductRepository(repoType string) ProductRepository {
	var repoProduct ProductRepository
	switch repoType {
	case "InMemory":
		repoProduct = memory.NewProductRepository()
	}
	return repoProduct
}
