package httpapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"purchase-cart-service/docs"
)

type Router struct {
	router *gin.Engine
}
type HandlersMethods struct {
	Method  string
	Route   string
	Handler gin.HandlerFunc
}
type IHandler interface {
	GetHandlers() []HandlersMethods
}

// @title Purchase Cart Service API
// @version 1.0
// @description API per la gestione degli ordini del carrello acquisti
// @host localhost:8080
// @BasePath /
// @schemes http

// NewRouter configures and returns the HTTP router for the service
func NewRouter() *Router {
	router := gin.Default()

	// Config Swagger runtime: metadati corretti per includere tutte le route
	docs.SwaggerInfo.Title = "Purchase Cart Service API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Description = "API per la gestione degli ordini del carrello acquisti"

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Router{router: router}
}

func (r *Router) RegisterMethods(group string, handlers ...IHandler) {
	routes := r.router.Group(group)
	for _, h := range handlers {
		for _, handler := range h.GetHandlers() {
			switch handler.Method {
			case "GET":
				routes.GET(handler.Route, handler.Handler)
			case "POST":
				routes.POST(handler.Route, handler.Handler)
			case "PUT":
				routes.PUT(handler.Route, handler.Handler)
			case "DELETE":
				routes.DELETE(handler.Route, handler.Handler)
			}
		}
	}
}

func (r *Router) Get() http.Handler {
	return r.router
}
