package cmd

import (
	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"

	"github.com/3Danger/telegram_bot/internal/build"
)

func postgresCmd(b *build.Build) *cobra.Command {
	command := &cobra.Command{
		Use:   "postgres",
		Short: "run db migrations for postgres",
		RunE: func(cmd *cobra.Command, args []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	postgres := func() (*migrate.Migrate, error) {
		return b.PostgresMigration()
	}

	command.AddCommand(up(postgres))
	command.AddCommand(down(postgres))
	command.AddCommand(step(postgres))
	command.AddCommand(version(postgres))

	return command
}
