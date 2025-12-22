package product

import (
	"context"
	"purchase-cart-service/repository"
	"purchase-cart-service/utils"
)

type Service struct {
	productRepo repository.ProductRepository
	vatRepo     repository.VatRateRepository
}

func NewService(productRepo repository.ProductRepository, vatRepo repository.VatRateRepository) *Service {
	return &Service{
		productRepo: productRepo,
		vatRepo:     vatRepo,
	}
}
func (s *Service) GetAllProducts(ctx context.Context, countryCode string) ([]Detail, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var productsDetail []Detail
	vatRate, err := s.vatRepo.GetVATRate(countryCode)
	if err != nil {
		return nil, err
	}
	for _, p := range products {

		vat := utils.Round2(p.Price * vatRate)
		productsDetail = append(productsDetail, Detail{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			PriceWithVAT: utils.Round2(p.Price + vat),
			Price:        p.Price,
			VAT:          vatRate,
		})
	}
	return productsDetail, nil
}
