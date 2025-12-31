package main

import (
	"context"

	"github.com/mgierok/guitar-specs2/backend/internal/app"
	"github.com/mgierok/guitar-specs2/backend/internal/config"
	"github.com/mgierok/guitar-specs2/backend/internal/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.Env)
	runner := app.NewRunner(cfg, logger)
	_ = runner.Run(context.Background())
}
