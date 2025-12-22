package memory

import (
	"context"
	"purchase-cart-service/models"
	"sync"
)

type ProductRepository struct {
	products map[string]models.Product
	mu       sync.RWMutex
}

func NewProductRepository() *ProductRepository {
	products := make(map[string]models.Product)
	products["prod1"] = models.Product{ID: "prod1", Name: "Product 1", Description: "Description of Product 1", Price: 10.0}
	products["prod2"] = models.Product{ID: "prod2", Name: "Product 2", Description: "Description of Product 2", Price: 20.0}
	products["prod3"] = models.Product{ID: "prod2", Name: "Product 3", Description: "Description of Product 3", Price: 20.0}
	products["prod4"] = models.Product{ID: "prod2", Name: "Product 4", Description: "Description of Product 4", Price: 20.0}
	products["prod5"] = models.Product{ID: "prod2", Name: "Product 5", Description: "Description of Product 5", Price: 20.0}

	return &ProductRepository{products: products}
}

func (p *ProductRepository) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if product, ok := p.products[id]; ok {
		return &product, nil
	}
	return nil, nil
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var products []models.Product
	for _, product := range p.products {
		products = append(products, product)
	}
	return products, nil
}
