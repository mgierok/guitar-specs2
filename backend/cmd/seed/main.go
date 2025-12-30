package main

import (
	"context"
	"log"

	"github.com/mgierok/guitar-specs2/backend/internal/config"
	"github.com/mgierok/guitar-specs2/backend/internal/db"
	dbsqlc "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
	"github.com/mgierok/guitar-specs2/backend/internal/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := db.Open(context.Background(), db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer pool.Close()

	data, err := seed.Load("/data/guitars.json")
	if err != nil {
		log.Fatalf("load seed: %v", err)
	}

	queries := dbsqlc.New(pool)
	if err := seed.Apply(context.Background(), queries, data); err != nil {
		log.Fatalf("apply seed: %v", err)
	}
}
