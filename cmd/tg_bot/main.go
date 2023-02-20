package main

import (
	"context"
	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/telegram"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt /*, syscall.SIGTERM*/)
	defer stop()
	erg, ctx := errgroup.WithContext(ctx)

	conf, logger, err := initiate()
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init")
	}
	tg, err := telegram.New(conf.Telegram, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't create tg api")
	}
	erg.Go(func() error {
		return tg.Start(ctx)
	})
	if err = erg.Wait(); err != nil {
		log.Error().Err(err).Msg("erg waited")
	}
}

func initiate() (*config.Config, zerolog.Logger, error) {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	var conf config.Config

	if err := envconfig.Process("", &conf); err != nil {
		return nil, logger, err
	}
	return &conf, logger, nil
}
