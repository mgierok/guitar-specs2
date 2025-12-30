package migrate

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Apply(ctx context.Context, pool *pgxpool.Pool, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, string(contents))
	return err
}
