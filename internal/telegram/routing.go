package telegram

import (
	"context"
	"fmt"
	"strings"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/internal/telegram/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/reply"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

func (t *Telegram) messageProcessor(ctx context.Context, key string, msg models.Data) error {
	defer func() {
		if strings.HasPrefix(key, "/") {
			if err := t.repo.command.Set(ctx, msg.UserID, key); err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("saving command state")
			}
		}
	}()
	switch key {
	case buttons.EndpointStart, buttons.EndpointHome:
		if err := t.handlerHome(ctx, msg); err != nil {
			return fmt.Errorf("handle home: %w", err)
		}

		return nil
	case buttons.EndpointRegistration:
		if err := t.handlerAuth(ctx, msg); err != nil {
			return fmt.Errorf("handle auth: %w", err)
		}

		return nil
	default:
		key, err := t.repo.command.Get(ctx, msg.UserID)
		if err != nil {
			return fmt.Errorf("getting last command: %w", err)
		}
		if key == nil {
			if err := t.handlerUndefined(msg); err != nil {
				return fmt.Errorf("handling undefined message: %w", err)
			}

			return nil
		}

		return t.messageProcessor(ctx, *key, msg)
	}

	return nil
}

func (t *Telegram) updateProcessor(ctx context.Context, update tele.Update) error {
	msg := models.NewMessage(update)

	switch {
	case len(msg.CallbackMap) != 0:
		if err := t.messageProcessor(ctx, msg.CallbackMap[buttons.KeyEndpoint], msg); err != nil {
			return fmt.Errorf("handle callback query: %w", err)
		}
	case msg.Message != "":
		if err := t.messageProcessor(ctx, msg.Message, msg); err != nil {
			return fmt.Errorf("processing message: %w", err)
		}
	default:
		if err := t.handlerUndefined(msg); err != nil {
			return fmt.Errorf("handling undefined message: %w", err)
		}
	}

	return nil
}

func (t *Telegram) handlerUndefined(msg models.Data) error {
	opt := reply.SendMessageOpts(
		reply.Row(reply.ButtonHome),
	)
	if _, err := t.bot.SendMessage(msg.ChatID, "Вы просите странного", opt); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
