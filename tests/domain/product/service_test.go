package product

import (
	"context"
	"purchase-cart-service/internal/domain/product"
	"purchase-cart-service/repository"
	"testing"
)

var svc *product.Service

func init() {
	svc = product.NewService(repository.NewProductRepository("InMemory"), repository.NewVatRateRepository("InMemory"))

}

func TestProduct_List(t *testing.T) {

	// Act
	list, err := svc.GetAllProducts(context.Background(), "IT")

	// Assert
	if err != nil {
		t.Fatalf("List errore inatteso: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("Catalogo prodotti vuoto, atteso almeno un prodotto precaricato")
	}
}

func TestProduct_GetByID(t *testing.T) {
	// Arrange

	id := "prod1"

	// Act
	p, err := svc.GetProductByID(context.Background(), id, "IT")

	// Assert
	if err != nil {
		t.Fatalf("GetByID errore inatteso: %v", err)
	}
	if p == nil || p.ID != id {
		t.Fatalf("Prodotto non trovato o ID non combacia, got=%v want=%s", p, id)
	}
}
