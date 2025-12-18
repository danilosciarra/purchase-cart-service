package repository

import "purchase-cart-service/repository/memory"

type VatRateRepository interface {
	GetVATRate(countryCode string) (float64, error)
}

func NewVatRateRepository(repoType string) VatRateRepository {
	var repo VatRateRepository
	switch repoType {
	case "InMemory":
		repo = memory.NewVatRateRepository()
	}
	return repo
}
