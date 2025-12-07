package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/Rasulikus/qaservice/internal/app"
	"github.com/Rasulikus/qaservice/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	server := app.App(cfg)
	log.Printf("listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
