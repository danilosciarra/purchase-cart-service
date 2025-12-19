package main

import (
	"log"
	"purchase-cart-service/cmd/server"
	"purchase-cart-service/internal/config"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %v", r)
		}
	}()

	// Load configuration
	cfg := config.Load()
	srv := server.New(cfg)

	log.Println("Purchase Cart Service started on :8080")
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
