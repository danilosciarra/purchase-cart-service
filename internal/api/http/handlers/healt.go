package handlers

import (
	"github.com/gin-gonic/gin"
	httpapi "purchase-cart-service/internal/api/http"
)

type HealthCheckHandler struct {
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) GetHandlers() []httpapi.HandlersMethods {
	return []httpapi.HandlersMethods{
		{
			Method:  "GET",
			Route:   "/health",
			Handler: h.Healthcheck,
		},
	}
}

// Health check
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {string} string "ok"
// @Router /health [get]
func (h *HealthCheckHandler) Healthcheck(c *gin.Context) {
	c.String(200, "ok")
}
