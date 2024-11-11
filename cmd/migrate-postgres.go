package cmd

import (
	"github.com/spf13/cobra"

	"github.com/3Danger/telegram_bot/internal/build"
)

func postgresCmd(b *build.Build) *cobra.Command {
	command := &cobra.Command{ //nolint:exhaustruct
		Use:   "postgres",
		Short: "run db migrations for postgres",
		RunE: func(cmd *cobra.Command, args []string) error {
			//nolint:wrapcheck
			return cmd.Usage()
		},
	}

	command.AddCommand(up(b.PostgresMigration))
	command.AddCommand(down(b.PostgresMigration))
	command.AddCommand(step(b.PostgresMigration))
	command.AddCommand(version(b.PostgresMigration))

	return command
}
