package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *telegram) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		m := tgbotapi.NewMessage(0, "nil")
		t.api.Send(m)
		fmt.Printf("%+v", update)
	}
	return nil
}

func (t *telegram) Stop() error {
	return nil
}
