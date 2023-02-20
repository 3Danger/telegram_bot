package store

import (
	"context"
	"fmt"
	"github.com/3Danger/telegram_bot/internal/config"
	pgx "github.com/jackc/pgx/v5"
	"time"
)

type store struct {
	conn             *pgx.Conn
	operationTimeout time.Duration
}

func New(ctx context.Context, cfg *config.Postgres) (*store, error) {
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Database,
	)
	conn, err := pgx.Connect(ctx, postgresUrl)
	if err != nil {
		return nil, err
	}
	return &store{conn: conn, operationTimeout: cfg.OperationTimeout}, nil
}
