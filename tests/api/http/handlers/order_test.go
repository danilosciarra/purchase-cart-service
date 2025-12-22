package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	httpapi "purchase-cart-service/internal/api/http"
	"purchase-cart-service/internal/api/http/handlers"
	"purchase-cart-service/internal/domain/order"
	"purchase-cart-service/repository"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouterForOrders() *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := handlers.NewOrderHandler(order.NewService(repository.NewOrderRepository("InMemory"), repository.NewVatRateRepository("InMemory"), repository.NewProductRepository("InMemory")))
	r := httpapi.NewRouter()
	r.RegisterMethods("/api/v1", h)
	return r.Engine()
}

// helper: crea un ordine e restituisce l'order_id
func createOrderForTest(t *testing.T, r *gin.Engine) string {
	t.Helper()
	body := map[string]any{
		"country_code": "IT",
		"items": []map[string]any{
			{"product_id": "prod1", "quantity": 2},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("creazione ordine fallita, got=%d body=%s", w.Code, w.Body.String())
	}
	var resp struct {
		OrderID string `json:"order_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response JSON non valido: %v", err)
	}
	if resp.OrderID == "" {
		t.Fatalf("OrderID mancante nella risposta di creazione")
	}
	return resp.OrderID
}

func TestCreateOrderHandler_OK(t *testing.T) {
	r := setupRouterForOrders()

	body := map[string]any{
		"country_code": "IT",
		"items": []map[string]any{
			{"product_id": "prod1", "quantity": 2},
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp struct {
		OrderID    string  `json:"order_id"`
		TotalPrice float64 `json:"total_price"`
		TotalVAT   float64 `json:"total_vat"`
		Items      []struct {
			ProductID string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			UnitPrice float64 `json:"unit_price"`
			VAT       float64 `json:"vat"`
		} `json:"items"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response JSON non valido: %v", err)
	}
	if resp.OrderID == "" {
		t.Errorf("OrderID mancante")
	}
	if len(resp.Items) == 0 {
		t.Fatalf("Items mancante/vuoto")
	}
	if resp.TotalPrice <= 0 {
		t.Errorf("TotalPrice non valido, got=%.2f", resp.TotalPrice)
	}
	if resp.TotalVAT < 0 {
		t.Errorf("TotalVAT non valido, got=%.2f", resp.TotalVAT)
	}
	for _, it := range resp.Items {
		if it.ProductID == "" {
			t.Errorf("ProductID mancante")
		}
		if it.Quantity <= 0 {
			t.Errorf("Quantity non valida, got=%d", it.Quantity)
		}
		if it.UnitPrice <= 0 {
			t.Errorf("UnitPrice non valido, got=%.2f", it.UnitPrice)
		}
		if it.VAT < 0 {
			t.Errorf("VAT non valida, got=%.2f", it.VAT)
		}
	}
}

// body privo di items → 400 Bad Request
func TestCreateOrderHandler_BadRequest_NoItems(t *testing.T) {
	r := setupRouterForOrders()

	body := map[string]any{
		"country_code": "IT",
		// items mancante
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// quantity non valida (<=0) → 400 Bad Request
func TestCreateOrderHandler_BadRequest_InvalidQuantity(t *testing.T) {
	r := setupRouterForOrders()

	body := map[string]any{
		"country_code": "IT",
		"items": []map[string]any{
			{"product_id": "prod1", "quantity": 0},
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// JSON non valido → 400 Bad Request
func TestCreateOrderHandler_BadRequest_InvalidJSON(t *testing.T) {
	r := setupRouterForOrders()

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// product non trovato → 404 Not Found
func TestCreateOrderHandler_ProductNotFound(t *testing.T) {
	r := setupRouterForOrders()

	body := map[string]any{
		"country_code": "IT",
		"items": []map[string]any{
			{"product_id": "unknown_prod", "quantity": 1},
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
}

// country_code non supportato → 400 Bad Request
func TestCreateOrderHandler_InvalidCountryCode(t *testing.T) {
	r := setupRouterForOrders()

	body := map[string]any{
		"country_code": "XX", // non supportato
		"items": []map[string]any{
			{"product_id": "prod1", "quantity": 1},
		},
	}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}

// GET /api/v1/orders → lista vuota
func TestGetOrders_Empty_OK(t *testing.T) {
	r := setupRouterForOrders()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var list []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("response JSON non valido: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("lista ordini non vuota, got=%d", len(list))
	}
}

// GET /api/v1/orders → lista con un ordine
func TestGetOrders_WithOne_OK(t *testing.T) {
	r := setupRouterForOrders()
	createdID := createOrderForTest(t, r)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var list []struct {
		OrderID string `json:"order_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("response JSON non valido: %v", err)
	}
	found := false
	for _, o := range list {
		if o.OrderID == createdID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("ordine %s non presente nella lista", createdID)
	}
}

// GET /api/v1/orders/:id → OK
func TestGetOrderByID_OK(t *testing.T) {
	r := setupRouterForOrders()
	createdID := createOrderForTest(t, r)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/orders/%s", createdID), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var ord struct {
		OrderID string `json:"order_id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &ord); err != nil {
		t.Fatalf("response JSON non valido: %v", err)
	}
	if ord.OrderID != createdID {
		t.Errorf("OrderID diverso, got=%s want=%s", ord.OrderID, createdID)
	}
}

// GET /api/v1/orders/:id → NotFound
func TestGetOrderByID_NotFound(t *testing.T) {
	r := setupRouterForOrders()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/non-existent-id", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status code errato, got=%d want=%d body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
}
