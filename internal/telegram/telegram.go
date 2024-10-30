package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/3Danger/telegram_bot/internal/config"
)

type Telegram interface {
	Start(ctx context.Context) error
	Stop() error
}

type telegram struct {
	api *tgbotapi.BotAPI
	cnf config.Telegram
}

func New(cnf config.Telegram) (Telegram, error) {
	api, err := tgbotapi.NewBotAPI(cnf.Token)
	if err != nil {
		return nil, fmt.Errorf("creating new telegram api: %w", err)
	}

	api.Debug = cnf.Debug

	return &telegram{api: api, cnf: cnf}, nil
}
