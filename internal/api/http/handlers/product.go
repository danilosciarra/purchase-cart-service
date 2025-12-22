package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	httpapi "purchase-cart-service/internal/api/http"
	"purchase-cart-service/internal/domain/product"
	"strings"
)

type ProductHandler struct {
	domain *product.Service
}

func NewProductHandler(domain *product.Service) *ProductHandler {
	return &ProductHandler{domain: domain}
}

type ProductResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	VAT          float64 `json:"vat"`
	PriceWithVAT float64 `json:"price_with_vat"`
}

func (h *ProductHandler) GetHandlers() []httpapi.HandlersMethods {
	return []httpapi.HandlersMethods{
		{
			Method:  "GET",
			Route:   "/products",
			Handler: h.GetAllProducts,
		},
	}
}

// GetAllProducts gestisce la richiesta per ottenere tutti i prodotti
// @Summary Get All Products
// @Description Retrieve a list of all products
// @Tags Products
// @Accept json
// @Param country_code query string false "Country Code for VAT calculation"
// @Produce json
// @Success 200 {array} ProductResponse
// @Failure 500 {object} ErrorResponse
// @Router  /api/v1/products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	countryCode := c.Query("country_code")
	products, err := h.domain.GetAllProducts(c.Request.Context(), strings.ToUpper(countryCode))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to retrieve products"})
		return
	}

	var response []ProductResponse
	for _, p := range products {
		response = append(response, ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			VAT:          p.VAT,
			PriceWithVAT: p.PriceWithVAT,
		})
	}

	c.JSON(http.StatusOK, response)
}
