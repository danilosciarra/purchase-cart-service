package order

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"math"
	"purchase-cart-service/models"
	"purchase-cart-service/repository"
	"time"
)

type Service struct {
	orderRepo repository.OrderRepository
	vatRepo   repository.VatRateRepository
}

func NewService(orderRepo repository.OrderRepository, vatRepo repository.VatRateRepository) *Service {
	return &Service{
		orderRepo: orderRepo,
		vatRepo:   vatRepo,
	}
}

var ErrInvalidItem = errors.New("invalid order item")

// CreateItem is the input DTO for creating orders
type CreateItem struct {
	ProductID string
	Quantity  int
	UnitPrice float64
}

func (s *Service) CreateOrder(ctx context.Context, countryCode string, items []CreateItem) (*models.Order, error) {
	if len(items) == 0 {
		return nil, ErrInvalidItem
	}

	order := &models.Order{
		ID:        uuid.NewString(),
		CreatedAt: time.Now(),
	}
	for _, it := range items {
		if it.Quantity <= 0 || it.UnitPrice <= 0 {
			return nil, ErrInvalidItem
		}

		linePrice := float64(it.Quantity) * it.UnitPrice
		vatRate, err := s.vatRepo.GetVATRate(countryCode)
		if err != nil {
			return nil, err
		}
		vat := round2(linePrice * vatRate)
		total := round2(linePrice + vat)

		order.Items = append(order.Items, models.Item{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			VAT:       total,
		})

		order.TotalVAT += vat
		order.TotalPrice += total
	}

	order.TotalVAT = round2(order.TotalVAT)
	order.TotalPrice = round2(order.TotalPrice)
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}
func (s *Service) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	return s.orderRepo.GetAll(ctx)
}

func round2(val float64) float64 {
	return math.Round(val*100) / 100
}
