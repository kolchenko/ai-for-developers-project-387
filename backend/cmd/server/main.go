package main

import (
	"log"
	"net/http"
	"time"

	"callcalendar/backend/internal/api"
	"callcalendar/backend/internal/config"
	"callcalendar/backend/internal/store"
)

func main() {
	cfg := config.Load()

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	handler := api.NewServerWithWeb(st, cfg.WebDir)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
