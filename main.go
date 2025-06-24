package main

import (
	"flag"
	"log"
	"net/http"

	"backend-challenge/api"
	"backend-challenge/impl"

	"github.com/go-chi/chi/v5"
)

func main() {
	flag.Parse()

	promoPath := "./promos"

	if err := impl.LoadPromoCodes(promoPath); err != nil {
		log.Fatalf("Failed to load promo codes: %v", err)
	}

	server := impl.NewServer()
	r := chi.NewRouter()

	api.HandlerWithOptions(server, api.ChiServerOptions{
		BaseURL:    "/api",
		BaseRouter: r,
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
