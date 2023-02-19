package main

import (
	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/rs/zerolog"
	"github.com/vrischmann/envconfig"
	"os"
)

func main() {
	conf, logger, err := initiate()
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init")
	}
	_ = conf
}

func initiate() (*config.Config, zerolog.Logger, error) {
	log := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	var conf *config.Config

	if err := envconfig.Init(&conf); err != nil {
		return nil, log, err
	}
	return conf, log, nil
}
