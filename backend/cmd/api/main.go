package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/mgierok/guitar-specs2/backend/internal/config"
	"github.com/mgierok/guitar-specs2/backend/internal/db"
	dbsqlc "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
	"github.com/mgierok/guitar-specs2/backend/internal/handlers"
	"github.com/mgierok/guitar-specs2/backend/internal/log"
	"github.com/mgierok/guitar-specs2/backend/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.Env)

	pool, err := db.Open(context.Background(), db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		logger.Error("db_connect_failed", "error", err)
		return
	}
	defer pool.Close()

	queries := dbsqlc.New(pool)
	guitarHandler := handlers.GuitarHandler{Queries: queries}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger(logger))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		r.Get("/guitars", guitarHandler.List)
		r.Get("/guitars/{slug}", guitarHandler.Detail)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("server_start", "port", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server_error", "error", err)
	}
}
