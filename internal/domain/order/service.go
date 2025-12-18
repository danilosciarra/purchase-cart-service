package order

import "purchase-cart-service/repository"

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
