package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/3Danger/telegram_bot/internal/build"
	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/services/auth"
	"github.com/3Danger/telegram_bot/internal/telegram"
)

func botCmd(ctx context.Context, b *build.Build) *cobra.Command {
	command := &cobra.Command{ //nolint:exhaustruct
		Use:       "bot",
		Short:     "run db migrations",
		ValidArgs: []string{"bot"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBot(ctx, b) //nolint:wrapcheck
		},
	}

	return command
}

func runBot(ctx context.Context, b *build.Build) error {
	cnf, err := config.New()
	if err != nil {
		return fmt.Errorf("creating config: %w", err)
	}

	var (
		repoPermanentUser = b.RepoUserPermanent()
		repoSessionUser   = b.RepoUserSession()
		repoState         = b.RepoState()
		repoCommand       = b.RepoState()
		//repoChainStates = b.RepoChainStates()
	)

	svcAuth := auth.NewService(repoPermanentUser, repoSessionUser)

	tg, err := telegram.New(
		cnf.Telegram,
		repoPermanentUser,
		repoState,
		repoCommand,
		svcAuth,
	)
	if err != nil {
		return fmt.Errorf("creating telegram client: %w", err)
	}

	if err = tg.Start(ctx); err != nil {
		return fmt.Errorf("starting telegram: %w", err)
	}

	return nil
}
