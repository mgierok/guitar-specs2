package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mgierok/guitar-specs2/backend/internal/config"
	"github.com/mgierok/guitar-specs2/backend/internal/db"
	dbsqlc "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
	"github.com/mgierok/guitar-specs2/backend/internal/handlers"
	"github.com/mgierok/guitar-specs2/backend/internal/middleware"
	"github.com/mgierok/guitar-specs2/backend/internal/migrate"
	"github.com/mgierok/guitar-specs2/backend/internal/seed"
)

type Runner struct {
	Config config.Config
	Logger *slog.Logger
	deps   BootstrapDeps
}

func NewRunner(cfg config.Config, logger *slog.Logger) *Runner {
	return &Runner{
		Config: cfg,
		Logger: logger,
		deps:   defaultBootstrapDeps(),
	}
}

func NewRunnerWithDeps(cfg config.Config, logger *slog.Logger, deps BootstrapDeps) *Runner {
	return &Runner{
		Config: cfg,
		Logger: logger,
		deps:   mergeBootstrapDeps(deps),
	}
}

func (r *Runner) Run(ctx context.Context) error {
	pool, queries, err := r.openDatabase(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := r.applyBootstrap(ctx, pool, queries); err != nil {
		return err
	}

	server := r.buildServer(queries)
	r.Logger.Info("server_start", "port", r.Config.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		r.Logger.Error("server_error", "error", err)
		return err
	}

	return nil
}

type BootstrapDeps struct {
	Migrate    func(ctx context.Context, pool *pgxpool.Pool, path string) error
	CountRows  func(ctx context.Context, pool *pgxpool.Pool) (int, error)
	LoadSeed   func(path string) (seed.Dataset, error)
	ApplySeed  func(ctx context.Context, queries *dbsqlc.Queries, data seed.Dataset) error
	LogApplied func(logger *slog.Logger, data seed.Dataset)
	LogSkip    func(logger *slog.Logger, count int)
}

func defaultBootstrapDeps() BootstrapDeps {
	return BootstrapDeps{
		Migrate: migrate.Apply,
		CountRows: func(ctx context.Context, pool *pgxpool.Pool) (int, error) {
			var count int
			err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM guitar").Scan(&count)
			return count, err
		},
		LoadSeed:  seed.Load,
		ApplySeed: seed.Apply,
		LogApplied: func(logger *slog.Logger, data seed.Dataset) {
			logger.Info("seed_applied", "guitars", len(data.Guitars))
		},
		LogSkip: func(logger *slog.Logger, count int) {
			logger.Info("seed_skip", "reason", "data_exists", "guitars", count)
		},
	}
}

func mergeBootstrapDeps(deps BootstrapDeps) BootstrapDeps {
	defaults := defaultBootstrapDeps()
	if deps.Migrate != nil {
		defaults.Migrate = deps.Migrate
	}
	if deps.CountRows != nil {
		defaults.CountRows = deps.CountRows
	}
	if deps.LoadSeed != nil {
		defaults.LoadSeed = deps.LoadSeed
	}
	if deps.ApplySeed != nil {
		defaults.ApplySeed = deps.ApplySeed
	}
	if deps.LogApplied != nil {
		defaults.LogApplied = deps.LogApplied
	}
	if deps.LogSkip != nil {
		defaults.LogSkip = deps.LogSkip
	}
	return defaults
}

func (r *Runner) openDatabase(ctx context.Context) (*pgxpool.Pool, *dbsqlc.Queries, error) {
	pool, err := db.Open(ctx, db.Config{
		Host:     r.Config.DBHost,
		Port:     r.Config.DBPort,
		User:     r.Config.DBUser,
		Password: r.Config.DBPassword,
		Name:     r.Config.DBName,
		SSLMode:  r.Config.DBSSLMode,
	})
	if err != nil {
		r.Logger.Error("db_connect_failed", "error", err)
		return nil, nil, err
	}

	return pool, dbsqlc.New(pool), nil
}

func (r *Runner) applyBootstrap(ctx context.Context, pool *pgxpool.Pool, queries *dbsqlc.Queries) error {
	if r.Config.AutoMigrate {
		if err := r.deps.Migrate(ctx, pool, r.Config.MigratePath); err != nil {
			r.Logger.Error("migrate_failed", "error", err)
			return err
		}
	}

	if !r.Config.AutoSeed {
		return nil
	}

	count, err := r.deps.CountRows(ctx, pool)
	if err != nil {
		r.Logger.Error("seed_check_failed", "error", err)
		return err
	}

	if count == 0 {
		data, err := r.deps.LoadSeed(r.Config.SeedPath)
		if err != nil {
			r.Logger.Error("seed_load_failed", "error", err)
			return err
		}
		if err := r.deps.ApplySeed(ctx, queries, data); err != nil {
			r.Logger.Error("seed_apply_failed", "error", err)
			return err
		}
		r.deps.LogApplied(r.Logger, data)
		return nil
	}

	r.deps.LogSkip(r.Logger, count)
	return nil
}

func (r *Runner) buildServer(queries *dbsqlc.Queries) *http.Server {
	router := r.buildRouter(queries)
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", r.Config.ServerPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (r *Runner) buildRouter(queries *dbsqlc.Queries) http.Handler {
	guitarHandler := handlers.GuitarHandler{Queries: queries}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(r.Logger))

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		r.Get("/guitars", guitarHandler.List)
		r.Get("/guitars/{slug}", guitarHandler.Detail)
	})

	return router
}
