package main

//go:generate sqlc generate

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/cmd"
	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/pkg/graceful"
)

//go:generate sqlc generate
func main() {
	cnf, err := config.New()
	if err != nil {
		panic(fmt.Errorf("creating config"))
	}

	ctx := context.Background()
	ctx = zerolog.New(os.Stdout).Level(cnf.App.LogLevel).With().Timestamp().Logger().WithContext(ctx)
	ctx = graceful.Context(ctx)

	if err = cmd.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			zerolog.Ctx(ctx).Info().Msg("graceful shutdown")

			return
		}

		zerolog.Ctx(ctx).Fatal().Err(err).Msg("unexpected error")
	}

	zerolog.Ctx(ctx).Info().Msg("exit")
}
