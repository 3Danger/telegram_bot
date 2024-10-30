package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/3Danger/telegram_bot/internal/build"
)

func Run(ctx context.Context) error {
	b, err := build.New(ctx)
	if err != nil {
		return fmt.Errorf("creating builder: %w", err)
	}

	root := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	root.AddCommand(
		botCmd(ctx, b),
		migrateCmd(b),
	)

	return errors.Wrap(root.ExecuteContext(ctx), "run application")
}
