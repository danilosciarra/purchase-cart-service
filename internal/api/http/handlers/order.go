package handlers

import (
	"github.com/gin-gonic/gin"
	"purchase-cart-service/internal/domain/order"
)

type OrderHandler struct {
	domain *order.Service
}

func NewOrderHandler(domain *order.Service) *OrderHandler {
	return &OrderHandler{domain: domain}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {

}
