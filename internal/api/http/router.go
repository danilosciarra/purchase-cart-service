package httpapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"purchase-cart-service/internal/api/http/handlers"
)

// NewRouter configures and returns the HTTP router for the service
func NewRouter(orderHandler *handlers.OrderHandler) http.Handler {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Orders
	router.POST("/orders", orderHandler.CreateOrder)

	return router
}
