package main

import (
	"fmt"
	"github.com/3Danger/telegram_bot/internal/config"
	tg_bot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/vrischmann/envconfig"
	"os"
)

func main() {
	conf, logger, err := initiate()
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init")
	}
	api, err := tg_bot.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't create tg api")
	}
	api.Debug = true

	u := tg_bot.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		fmt.Printf("%+v", update)
	}
}

func initiate() (*config.Config, zerolog.Logger, error) {
	log := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)
	var conf *config.Config

	if err := envconfig.Init(&conf); err != nil {
		return nil, log, err
	}
	return conf, log, nil
}
