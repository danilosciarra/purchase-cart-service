package order

import (
	"context"
	"errors"
	"purchase-cart-service/internal/domain/product"
	"purchase-cart-service/models"
	"purchase-cart-service/repository"
	"purchase-cart-service/utils"
)

type Service struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	vatRepo     repository.VatRateRepository
}

func NewService(orderRepo repository.OrderRepository, vatRepo repository.VatRateRepository, productRepo repository.ProductRepository) *Service {
	return &Service{
		orderRepo:   orderRepo,
		vatRepo:     vatRepo,
		productRepo: productRepo,
	}
}

var ErrInvalidItem = errors.New("invalid order item")
var ErrInvalidVATRate = errors.New("invalid VAT rate")
var ErrProductNotFound = errors.New("product not found")

func (s *Service) CreateOrder(ctx context.Context, countryCode string, items []CreateItem) (*models.Order, error) {
	if len(items) == 0 {
		return nil, ErrInvalidItem
	}
	vatRate, err := s.vatRepo.GetVATRate(countryCode)
	if err != nil {
		return nil, ErrInvalidVATRate
	}
	order := &models.Order{}
	for _, it := range items {
		product, err := s.productRepo.GetProduct(ctx, it.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, ErrProductNotFound
		}
		if it.Quantity <= 0 || product.Price <= 0 {
			return nil, ErrInvalidItem
		}

		linePrice := float64(it.Quantity) * product.Price

		vat := utils.Round2(linePrice * vatRate)
		total := utils.Round2(linePrice + vat)

		order.Items = append(order.Items, models.Item{
			ProductID: it.ProductID,
			Name:      product.Name,
			Quantity:  it.Quantity,
			UnitPrice: product.Price,
			VAT:       total,
		})

		order.TotalVAT += vat
		order.TotalPrice += total
	}

	order.TotalVAT = utils.Round2(order.TotalVAT)
	order.TotalPrice = utils.Round2(order.TotalPrice)
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id string) (*Detail, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, nil
	}
	return s.GetOrderDetail(ctx, order)

}
func (s *Service) GetAllOrders(ctx context.Context) ([]*Detail, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var details []*Detail
	for _, order := range orders {
		detail, err := s.GetOrderDetail(ctx, order)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}

func (s *Service) GetOrderDetail(ctx context.Context, order *models.Order) (*Detail, error) {
	var products []ProductDetail
	for _, item := range order.Items {
		product, err := s.productRepo.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product != nil {
			product.Price = item.UnitPrice
			products = append(products, ProductDetail{
				Product:  *product,
				VAT:      item.VAT,
				Quantity: item.Quantity,
			})
		}
	}
	return &Detail{
		Id:         order.ID,
		TotalPrice: order.TotalPrice,
		TotalVAT:   order.TotalVAT,
		Items:      products,
	}, nil
}

func (s *Service) GetAllProducts(ctx context.Context) ([]product.Detail, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var details []product.Detail
	for _, p := range products {
		details = append(details, product.Detail{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return details, nil
}
