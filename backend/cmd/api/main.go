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
	"github.com/mgierok/guitar-specs2/backend/internal/migrate"
	"github.com/mgierok/guitar-specs2/backend/internal/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.Env)

	ctx := context.Background()
	pool, err := db.Open(ctx, db.Config{
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

	if cfg.AutoMigrate {
		if err := migrate.Apply(ctx, pool, cfg.MigratePath); err != nil {
			logger.Error("migrate_failed", "error", err)
			return
		}
	}

	if cfg.AutoSeed {
		var count int
		if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM guitar").Scan(&count); err != nil {
			logger.Error("seed_check_failed", "error", err)
			return
		}

		if count == 0 {
			data, err := seed.Load(cfg.SeedPath)
			if err != nil {
				logger.Error("seed_load_failed", "error", err)
				return
			}
			if err := seed.Apply(ctx, queries, data); err != nil {
				logger.Error("seed_apply_failed", "error", err)
				return
			}
			logger.Info("seed_applied", "guitars", len(data.Guitars))
		} else {
			logger.Info("seed_skip", "reason", "data_exists", "guitars", count)
		}
	}

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
