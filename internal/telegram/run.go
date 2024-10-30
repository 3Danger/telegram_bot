package telegram

import (
	"context"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (t *telegram) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case u := <-updates:
			// if u.Message == nil {
			//	continue
			//}

			zerolog.Ctx(ctx).Info().Interface("sent from", u.SentFrom()).Send()

			data, _ := json.MarshalIndent(u.Message, "", "  ")
			m := tgbotapi.NewMessage(u.FromChat().ID, string(data))

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
