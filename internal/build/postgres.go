package build

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	migration "github.com/3Danger/telegram_bot/internal/migrations"
)

func (b *Build) Postgres() *pgx.Conn {
	return b.db
}

func (b *Build) PostgresMigration() (*migrate.Migrate, error) {
	d, err := iofs.New(migration.FS, migration.PostgresPath)
	if err != nil {
		return nil, errors.Wrap(err, "embed postgres migrations")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, b.cnf.Postgres.DSN())
	if err != nil {
		return nil, errors.Wrap(err, "apply postgres migrations")
	}

	return m, nil
}
