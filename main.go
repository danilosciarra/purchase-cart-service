package main

import (
	"fmt"
	"log"
	"net/http"
	"purchase-cart-service/repository"

	httpapi "purchase-cart-service/internal/api/http"
	"purchase-cart-service/internal/api/http/handlers"
	"purchase-cart-service/internal/config"
	"purchase-cart-service/internal/domain/order"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %v", r)
		}
	}()

	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure
	orderRepo := repository.NewOrderRepository(cfg.Database.Type)
	vatRepo := repository.NewVatRateRepository(cfg.Database.Type)
	//priceCalculator := pricing.NewCalculator(cfg.VATRate)

	// Initialize domain services
	orderService := order.NewService(orderRepo, vatRepo)

	// Initialize HTTP handlers
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup router
	router := httpapi.NewRouter(orderHandler)

	log.Println("Purchase Cart Service started on :8080")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.WebApp.HostName, cfg.WebApp.Port), router); err != nil {
		log.Fatal(err)
	}
}
