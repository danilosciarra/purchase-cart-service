package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	httpapi "purchase-cart-service/internal/api/http"
	"purchase-cart-service/internal/domain/order"
	"strings"
)

type OrderHandler struct {
	domain *order.Service
}

type OrderRequest struct {
	Items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
	CountryCode string `json:"country_code"`
}

// OrderResponse rappresenta la risposta dopo la creazione di un ordine
type OrderResponse struct {
	OrderID    string           `json:"order_id"`
	TotalPrice float64          `json:"total_price"`
	TotalVAT   float64          `json:"total_vat"`
	Items      []orderItemReply `json:"items"`
}

type orderItemReply struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	VAT       float64 `json:"vat"`
}

// ErrorResponse rappresenta una risposta di errore
type ErrorResponse struct {
	Message string `json:"message"`
}

func NewOrderHandler(domain *order.Service) *OrderHandler {
	return &OrderHandler{domain: domain}
}

func (h *OrderHandler) GetHandlers() []httpapi.HandlersMethods {
	return []httpapi.HandlersMethods{
		{
			Method:  "PUT",
			Route:   "/orders",
			Handler: h.CreateOrder,
		},
		{
			Method:  "GET",
			Route:   "/orders/:id",
			Handler: h.GetOrder,
		},
		{
			Method:  "GET",
			Route:   "/orders",
			Handler: h.GetOrders,
		},
	}
}

// CreateOrder Orders
// @Summary Crea un nuovo ordine
// @Description Crea un nuovo ordine nel sistema
// @Tags Orders
// @Accept json
// @Produce json
// @Param order body handlers.OrderRequest true "Dati ordine"
// @Success 201 {object} handlers.OrderResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Router /api/v1/orders [put]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request"})
		return
	}
	if req.CountryCode == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Country code is required"})
		return
	}
	items := make([]order.CreateItem, 0, len(req.Items))
	for _, it := range req.Items {
		if it.Quantity == 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Item quantity must be greater than zero"})
			return
		}
		items = append(items, order.CreateItem{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})
	}
	ord, err := h.domain.CreateOrder(c.Request.Context(), strings.ToUpper(req.CountryCode), items)
	if err != nil {
		if err == order.ErrInvalidItem {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid item in order"})
			return
		}
		if err == order.ErrInvalidVATRate {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid VAT rate for country"})
			return
		}
		if err == order.ErrProductNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	resp := OrderResponse{
		OrderID:    ord.ID,
		TotalPrice: ord.TotalPrice,
		TotalVAT:   ord.TotalVAT,
	}

	for _, it := range ord.Items {
		resp.Items = append(resp.Items, orderItemReply{
			ProductID: it.ProductID,
			Name:      it.Name,
			Quantity:  it.Quantity,
			UnitPrice: it.UnitPrice,
			VAT:       it.VAT,
		})
	}
	c.JSON(http.StatusCreated, resp)

}

// GetOrder
// @Summary Ottieni un ordine per ID
// @Description Recupera i dettagli di un ordine utilizzando il suo ID
// @Tags Orders
// @Produce json
// @Param id path string true "ID Ordine"
// @Success 200 {object} handlers.OrderResponse
// @Failure 400 {object} handlers.ErrorResponse
// @Failure 404 {object} handlers.ErrorResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid order ID"})
		return
	}
	ord, err := h.domain.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	if ord == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Order not found"})
		return
	}
	resp := OrderResponse{
		OrderID:    ord.Id,
		TotalPrice: ord.TotalPrice,
		TotalVAT:   ord.TotalVAT,
	}

	for _, it := range ord.Items {
		resp.Items = append(resp.Items, orderItemReply{
			ProductID: it.ID,
			Name:      it.Name,
			Quantity:  it.Quantity,
			UnitPrice: it.Price,
			VAT:       it.VAT,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// GetOrders
// @Summary Elenca tutti gli ordini
// @Description Recupera una lista di tutti gli ordini
// @Tags Orders
// @Produce json
// @Success 200 {array} handlers.OrderResponse
// @Failure 500 {object} handlers.ErrorResponse
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	// Implementation for listing all orders can be added here
	var resp []OrderResponse
	orders, err := h.domain.GetAllOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	for _, ord := range orders {
		orderResp := OrderResponse{
			OrderID:    ord.Id,
			TotalPrice: ord.TotalPrice,
			TotalVAT:   ord.TotalVAT,
		}

		for _, it := range ord.Items {
			orderResp.Items = append(orderResp.Items, orderItemReply{
				ProductID: it.ID,
				Name:      it.Name,
				Quantity:  it.Quantity,
				UnitPrice: it.Price,
				VAT:       it.VAT,
			})
		}
		resp = append(resp, orderResp)
	}
	c.JSON(http.StatusOK, resp)
}
