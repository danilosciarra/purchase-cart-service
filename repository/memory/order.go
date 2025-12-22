package memory

import (
	"context"
	"github.com/google/uuid"
	"purchase-cart-service/models"
	"sync"
	"time"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*models.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]*models.Order)}
}

func (o *OrderRepository) Save(ctx context.Context, order *models.Order) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	order.ID = uuid.NewString()
	order.CreatedAt = time.Now()
	o.orders[order.ID] = order
	return nil
}
func (o *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if o, ok := o.orders[id]; ok {
		return o, nil
	}
	return nil, nil
}
func (o *OrderRepository) GetAll(ctx context.Context) ([]*models.Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	var orders []*models.Order
	for _, order := range o.orders {
		orders = append(orders, order)
	}
	return orders, nil

}
