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
func (s *Service) GetProductByID(ctx context.Context, productID string, countryCode string) (*Detail, error) {
	product, err := s.productRepo.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}
	vatRate, err := s.vatRepo.GetVATRate(countryCode)
	if err != nil {
		return nil, err
	}
	vat := utils.Round2(product.Price * vatRate)
	return &Detail{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		PriceWithVAT: utils.Round2(product.Price + vat),
		Price:        product.Price,
		VAT:          vatRate,
	}, nil
}
