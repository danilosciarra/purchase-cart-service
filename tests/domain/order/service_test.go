package order

import (
	"context"
	"purchase-cart-service/internal/domain/order"
	"purchase-cart-service/repository"
	"testing"
)

var svc *order.Service

func init() {
	svc = order.NewService(repository.NewOrderRepository("InMemory"), repository.NewVatRateRepository("InMemory"), repository.NewProductRepository("InMemory"))

}

func TestCreateOrder_CalcolaTotaliEIVA(t *testing.T) {
	// Arrange
	svc := order.NewService(repository.NewOrderRepository("InMemory"), repository.NewVatRateRepository("InMemory"), repository.NewProductRepository("InMemory"))
	req := []order.CreateItem{
		{ProductID: "prod1", Quantity: 2},
		{ProductID: "prod2", Quantity: 1},
	}

	// Act
	res, err := svc.CreateOrder(context.Background(), "IT", req)

	// Assert
	if err != nil {
		t.Fatalf("CreateOrder errore inatteso: %v", err)
	}
	// Totale netto = 2*10 + 1*20 = 48.80
	// IVA 22% = 8.80
	// Totale lordo = 48.80
	const expectedTotalVat = 8.80
	const expectedTotalPrice = 48.80

	if diff := res.TotalVAT - expectedTotalVat; diff < -0.0001 || diff > 0.0001 {
		t.Errorf("TotalVAT errato, got=%.2f want=%.2f", res.TotalVAT, expectedTotalVat)
	}
	if diff := res.TotalPrice - expectedTotalPrice; diff < -0.0001 || diff > 0.0001 {
		t.Errorf("TotalPrice errato, got=%.2f want=%.2f", res.TotalPrice, expectedTotalPrice)
	}

	// Verifica IVA di riga
	// Riga prod1: 2*10 = 20; IVA = 4.40
	// Riga prod2: 1*5 = 5; IVA = 1.10
	var prod1VAT, prod2VAT float64
	for _, it := range res.Items {
		switch it.ProductID {
		case "prod1":
			prod1VAT = it.VAT - it.UnitPrice*float64(it.Quantity)
		case "prod2":
			prod2VAT = it.VAT - it.UnitPrice*float64(it.Quantity)
		}
	}
	if diff := prod1VAT - 4.40; diff < -0.0001 || diff > 0.0001 {
		t.Errorf("IVA prod1 errata, got=%.2f want=%.2f", prod1VAT, 20.0)
	}
	if diff := prod2VAT - 4.40; diff < -0.0001 || diff > 0.0001 {
		t.Errorf("IVA prod2 errata, got=%.2f want=%.2f", prod2VAT, 20.0)
	}

}
