package httpapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"purchase-cart-service/docs"
	_ "purchase-cart-service/docs"
	"purchase-cart-service/internal/api/http/handlers"
)

// @title Purchase Cart Service API
// @version 1.0
// @description API per la gestione degli ordini del carrello acquisti
// @host localhost:8080
// @BasePath /

// NewRouter configures and returns the HTTP router for the service
func NewRouter(orderHandler *handlers.OrderHandler) http.Handler {
	router := gin.Default()
	hCheck := &handlers.HealthCheckHandler{}
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	v1.GET("/health", hCheck.Healthcheck)

	v1.GET("/orders/:id", orderHandler.GetOrder)
	v1.GET("/orders", orderHandler.GetOrders)
	v1.POST("/orders", orderHandler.CreateOrder)

	return router
}
