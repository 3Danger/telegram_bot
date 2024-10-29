package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/internal/config"
)

type Telegram interface {
	Start(ctx context.Context) error
	Stop() error
}

type telegram struct {
	api *tgbotapi.BotAPI
	log zerolog.Logger
	cnf config.Telegram
}

func New(cnf config.Telegram, logger zerolog.Logger) (Telegram, error) {
	logger = logger.With().Str("service", "telegram").Logger()

	api, err := tgbotapi.NewBotAPI(cnf.Token)
	if err != nil {
		return nil, fmt.Errorf("creating new telegram api: %w", err)
	}
	api.Debug = cnf.Debug
	return &telegram{api: api, log: logger, cnf: cnf}, nil
}
