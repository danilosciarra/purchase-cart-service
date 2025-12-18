package handlers

import "github.com/gin-gonic/gin"

type HealthCheckHandler struct {
}

// Health check
// @Summary Health check
// @Description Verifica lo stato del servizio
// @Tags health
// @Success 200 {string} string "ok"
// @Router /health [get]
func (h *HealthCheckHandler) Healthcheck(c *gin.Context) {
	c.String(200, "ok")
}
