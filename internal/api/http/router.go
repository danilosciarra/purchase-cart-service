package httpapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"purchase-cart-service/internal/api/http/handlers"
	"purchase-cart-service/models"
)

// NewRouter configures and returns the HTTP router for the service
func NewRouter(orderHandler *handlers.OrderHandler) http.Handler {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.ApiResponse{Result: "OK"})
	})

	// Orders
	router.POST("/orders", orderHandler.CreateOrder)

	return router
}
