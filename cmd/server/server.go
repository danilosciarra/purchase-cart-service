package server

import (
	"fmt"
	"net/http"
	httpapi "purchase-cart-service/internal/api/http"
	"purchase-cart-service/internal/api/http/handlers"
	"purchase-cart-service/internal/config"
	"purchase-cart-service/internal/domain/order"
	"purchase-cart-service/internal/domain/product"
	"purchase-cart-service/repository"
)

type Server struct {
	router   *httpapi.Router
	hostname string
	port     int
}

func New(cfg *config.Config) *Server {
	orderRepo := repository.NewOrderRepository(cfg.Database.Type)
	vatRepo := repository.NewVatRateRepository(cfg.Database.Type)
	productRepo := repository.NewProductRepository(cfg.Database.Type)
	srv := &Server{
		router:   httpapi.NewRouter(),
		hostname: cfg.WebApp.HostName,
		port:     cfg.WebApp.Port,
	}
	hc := handlers.NewHealthCheckHandler()
	oh := handlers.NewOrderHandler(order.NewService(orderRepo, vatRepo, productRepo))
	ph := handlers.NewProductHandler(product.NewService(productRepo, vatRepo))
	srv.router.RegisterMethods("/", hc)
	srv.router.RegisterMethods("/api/v1", oh, ph)
	return srv
}

func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.hostname, s.port), s.router.Get())
}
