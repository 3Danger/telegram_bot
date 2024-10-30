package build

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

func New(ctx context.Context) (*Build, error) {
	cnf, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("creating config: %w", err)
	}

	builder, err := createBuilder(ctx, cnf)
	if err != nil {
		return nil, fmt.Errorf("creating builder: %w", err)
	}

	return builder, nil
}

func createBuilder(ctx context.Context, cnf *config.Config) (b *Build, err error) {
	b = new(Build)

	if b.db, err = pgx.Connect(ctx, cnf.Postgres.DSN()); err != nil {
		return nil, fmt.Errorf("creating postgres db: %w", err)
	}

	b.cnf = cnf

	return b, nil
}
