package builder

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/3Danger/telegram_bot/internal/config"
)

type Build struct {
	db *pgx.Conn

	cnf *config.Config
}

func New(ctx context.Context, cnf *config.Config) (b *Build, err error) {
	b = new(Build)

	b.cnf = cnf

	if b.db, err = postgres(ctx, &cnf.Postgres); err != nil {
		return nil, fmt.Errorf("creating postgres db: %w", err)
	}

	return b, nil
}
