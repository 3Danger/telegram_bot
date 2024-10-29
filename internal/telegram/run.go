package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *telegram) Start(ctx context.Context) error {
	var (
		log     = t.log.With().Str("method", "run").Logger()
		updates tgbotapi.UpdatesChannel
	)
	{
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates = t.api.GetUpdatesChan(u)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case u := <-updates:
			if u.Message == nil {
				continue
			}

			fmt.Printf("ID %d", u.SentFrom().ID)

			m := tgbotapi.NewMessage(u.FromChat().ID, fmt.Sprint(u.UpdateID))
			msg, err := t.api.Send(m)
			if err != nil {
				log.Error().Err(err).Msg("couldn't send message")
			}
			_ = msg
		}
	}
}

func (t *telegram) Stop() error {
	return nil
}
