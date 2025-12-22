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
		{
			Method:  "GET",
			Route:   "/products/:id",
			Handler: h.GetProductByID,
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

// GetProductByID gestisce la richiesta per ottenere un prodotto per ID
// @Summary Get Product by ID
// @Description Retrieve a product by its ID
// @Tags Products
// @Accept json
// @Param id path string true "Product ID"
// @Param country_code query string false "Country Code for VAT calculation"
// @Produce json
// @Success 200 {object} ProductResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router  /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	countryCode := c.Query("country_code")
	productDetail, err := h.domain.GetProductByID(c.Request.Context(), productID, strings.ToUpper(countryCode))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to retrieve product"})
		return
	}
	if productDetail == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Product not found"})
		return
	}

	response := ProductResponse{
		ID:           productDetail.ID,
		Name:         productDetail.Name,
		Description:  productDetail.Description,
		Price:        productDetail.Price,
		VAT:          productDetail.VAT,
		PriceWithVAT: productDetail.PriceWithVAT,
	}

	c.JSON(http.StatusOK, response)

}
