package memory

import "errors"

type VatRateRepository struct {
	vatRates map[string]float64
}

func NewVatRateRepository() *VatRateRepository {
	return &VatRateRepository{
		vatRates: map[string]float64{
			"US": 0.0,
			"UK": 0.2,
			"DE": 0.19,
			"FR": 0.2,
			"IT": 0.22,
		},
	}
}
func (v *VatRateRepository) GetVATRate(countryCode string) (float64, error) {
	rate, exists := v.vatRates[countryCode]
	if !exists {
		return 0, errors.New("VatRate not found")
	}
	return rate, nil
}
