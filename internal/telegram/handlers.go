package telegram

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/reply"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

func (t *Telegram) handlerHome(ctx context.Context, msg models.Data) error {
	u, err := t.repo.user.Get(ctx, msg.UserID)
	if err != nil {
		return fmt.Errorf("getting user from repo: %w", err)
	}

	if u == nil {
		text := "Добро пожаловать!\nДля работы необходимо зарегистрироваться"

		opts := inline.SendMessageOpts(
			inline.Text(buttons.EndpointRegistration, models.PairKeyValues{
				Key:   buttons.KeyEndpoint,
				Value: buttons.EndpointRegistration,
			}),
		)

		if _, err = t.bot.SendMessage(msg.ChatID, text, opts); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		return nil
	}

	if u.Type == user.TypeSupplier {
		return t.handlerSupplierHome(ctx, msg)
	}

	return t.handlerCustomerHome(ctx, msg)
}

func (t *Telegram) handlerSupplierHome(ctx context.Context, msg models.Data) error {
	text := ""

	opts := reply.SendMessageOpts(
		reply.ButtonSupplierPostItems,
		reply.ButtonSupplierShowItems,
	)

	if _, err := t.bot.SendMessage(msg.ChatID, text, opts); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (t *Telegram) handlerCustomerHome(ctx context.Context, msg models.Data) error {
	text := ""

	opts := reply.SendMessageOpts(
		reply.ButtonCustomerShowItems,
	)

	if _, err := t.bot.SendMessage(msg.ChatID, text, opts); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
