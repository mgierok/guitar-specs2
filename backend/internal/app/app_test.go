package app

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mgierok/guitar-specs2/backend/internal/config"
	dbsqlc "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
	"github.com/mgierok/guitar-specs2/backend/internal/seed"
)

func TestApplyBootstrap_MigrateError(t *testing.T) {
	cfg := config.Config{AutoMigrate: true}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	deps := BootstrapDeps{
		Migrate: func(ctx context.Context, _ *pgxpool.Pool, _ string) error {
			return errors.New("migrate failed")
		},
	}

	runner := NewRunnerWithDeps(cfg, logger, deps)
	err := runner.applyBootstrap(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected migrate error, got nil")
	}
}

func TestApplyBootstrap_SeedApplied(t *testing.T) {
	cfg := config.Config{AutoSeed: true}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	called := struct {
		load   bool
		apply  bool
		logged bool
	}{}

	deps := BootstrapDeps{
		CountRows: func(ctx context.Context, _ *pgxpool.Pool) (int, error) {
			return 0, nil
		},
		LoadSeed: func(_ string) (seed.Dataset, error) {
			called.load = true
			return seed.Dataset{Guitars: []seed.Guitar{{ID: "1"}}}, nil
		},
		ApplySeed: func(ctx context.Context, _ *dbsqlc.Queries, _ seed.Dataset) error {
			called.apply = true
			return nil
		},
		LogApplied: func(_ *slog.Logger, _ seed.Dataset) {
			called.logged = true
		},
	}

	runner := NewRunnerWithDeps(cfg, logger, deps)
	err := runner.applyBootstrap(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called.load || !called.apply || !called.logged {
		t.Fatalf("expected seed load/apply/log, got load=%v apply=%v logged=%v", called.load, called.apply, called.logged)
	}
}

func TestApplyBootstrap_SeedSkipped(t *testing.T) {
	cfg := config.Config{AutoSeed: true}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	called := struct {
		skip bool
	}{}

	deps := BootstrapDeps{
		CountRows: func(ctx context.Context, _ *pgxpool.Pool) (int, error) {
			return 2, nil
		},
		LoadSeed: func(_ string) (seed.Dataset, error) {
			t.Fatal("did not expect LoadSeed to be called")
			return seed.Dataset{}, nil
		},
		ApplySeed: func(ctx context.Context, _ *dbsqlc.Queries, _ seed.Dataset) error {
			t.Fatal("did not expect ApplySeed to be called")
			return nil
		},
		LogSkip: func(_ *slog.Logger, _ int) {
			called.skip = true
		},
	}

	runner := NewRunnerWithDeps(cfg, logger, deps)
	err := runner.applyBootstrap(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !called.skip {
		t.Fatal("expected seed skip to be logged")
	}
}
