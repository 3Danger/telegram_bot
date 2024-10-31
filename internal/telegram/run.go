package telegram

import (
	"context"
	"encoding/json"

	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func (t *Telegram) Start(ctx context.Context) error {
	u := api.NewUpdate(0)
	u.Timeout = 60
	updates := t.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case u := <-updates:
			// if u.Text == nil {
			//	continue
			//}

			zerolog.Ctx(ctx).Info().Interface("sent from", u.SentFrom()).Send()

			data, _ := json.MarshalIndent(u.Message, "", "  ")
			m := api.NewMessage(u.FromChat().ID, string(data))

			registration := api.NewKeyboardButton("Регистрация")
			registration.RequestContact = true
			registration.RequestLocation = true

			location := api.NewKeyboardButtonLocation("Поделись локацией черт")
			location.RequestLocation = true

			keyboard := api.NewReplyKeyboard(
				api.NewKeyboardButtonRow(registration, location),
			)

			m.ReplyMarkup = keyboard

			msg, err := t.api.Send(m)
			if err != nil {
				log.Error().Err(err).Msg("couldn't send message")
			}

			_ = msg
		}
	}
}
