package cmd

import (
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/3Danger/telegram_bot/internal/build"
)

func migrateCmd(b *build.Build) *cobra.Command {
	command := &cobra.Command{ //nolint:exhaustruct
		Use:       "migrate",
		Short:     "run db migrations",
		ValidArgs: []string{"postgres"},
		RunE: func(cmd *cobra.Command, args []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	command.AddCommand(
		postgresCmd(b),
	)

	return command
}

type migrationConstructFn func() (*migrate.Migrate, error)

func up(constructFn migrationConstructFn) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "up",
		Short: "up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := constructFn()
			if err != nil {
				return errors.Wrap(err, "construct migration")
			}

			err = m.Up()
			if err != nil {
				if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, migrate.ErrNilVersion) {
					return nil
				}

				return errors.Wrap(err, "up migrations")
			}

			return nil
		},
	}
}

func down(constructFn migrationConstructFn) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "down",
		Short: "rollback all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := constructFn()
			if err != nil {
				return errors.Wrap(err, "construct migration")
			}

			err = m.Down()
			if err != nil {
				return errors.Wrap(err, "down migrations")
			}

			return nil
		},
	}
}

func step(constructFn migrationConstructFn) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "step",
		Short: "run N steps migrations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stepsCount, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrap(err, "steps count must be integer")
			}

			m, err := constructFn()
			if err != nil {
				return errors.Wrap(err, "construct migration")
			}

			err = m.Steps(stepsCount)
			if err != nil {
				return errors.Wrap(err, "step count migrations")
			}

			return nil
		},
	}
}

func version(constructFn migrationConstructFn) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "version",
		Short: "display current migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := constructFn()
			if err != nil {
				return errors.Wrap(err, "construct migration")
			}

			ver, dirty, err := m.Version()
			if err != nil {
				return errors.Wrap(err, "display migration version")
			}

			fmt.Printf("current version: %d, dirty: %v", ver, dirty) //nolint:forbidigo

			return nil
		},
	}
}
