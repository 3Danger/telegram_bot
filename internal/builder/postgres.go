package builder

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/3Danger/telegram_bot/internal/config"
)

func (b *Build) Postgres() *pgx.Conn {
	return b.db
}

func postgres(ctx context.Context, cfg *config.Postgres) (*pgx.Conn, error) {
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Database,
	)
	conn, err := pgx.Connect(ctx, postgresUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
