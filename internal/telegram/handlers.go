package telegram

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/inline"
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
			inline.Row(inline.Text(buttons.EndpointRegistration,
				models.Pair{buttons.KeyEndpoint: buttons.EndpointRegistration},
			)),
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
	text := "Главное меню"

	opts := inline.SendMessageOpts(
		inline.Row(
			inline.Text(buttons.ButtonSupplierShowItems, models.Pair{buttons.KeyEndpoint: "/show_goods"}),
			inline.Text(buttons.ButtonSupplierPostItems, models.Pair{buttons.KeyEndpoint: "/post_goods"}),
		),
	)

	if _, err := t.bot.SendMessage(msg.ChatID, text, opts); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (t *Telegram) handlerCustomerHome(ctx context.Context, msg models.Data) error {
	text := "Главное меню"

	opts := inline.SendMessageOpts(
		inline.Row(
			inline.Text("Показать товары", models.Pair{buttons.KeyEndpoint: "/show_goods"}),
		),
	)

	if _, err := t.bot.SendMessage(msg.ChatID, text, opts); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
