package handlers

import (
	"net/http"
	"net/http/httptest"
	"purchase-cart-service/internal/api/http/handlers"
	"purchase-cart-service/internal/domain/product"
	"purchase-cart-service/repository"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouterForProducts() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewProductHandler(product.NewService(repository.NewProductRepository("InMemory"), repository.NewVatRateRepository("InMemory")))
	// registra i percorsi /api/v1/products con repository InMemory precaricato
	r.GET("/api/v1/products", h.GetAllProducts)
	r.GET("/api/v1/products/:id", h.GetProductByID)
	return r
}

func TestListProductsHandler_OK(t *testing.T) {
	r := setupRouterForProducts()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?country_code=IT", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

}

func TestGetProductHandler_OK(t *testing.T) {
	r := setupRouterForProducts()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/prod1?country_code=IT", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestGetProductHandler_NotFound(t *testing.T) {
	r := setupRouterForProducts()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/UNKNOWN?country_code=IT", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
}
