package telegram

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/menu"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

func (t *Telegram) handlerHome(ctx context.Context, msg models.Request) error {
	u, err := t.repo.user.Get(ctx, msg.UserID())
	if err != nil {
		return fmt.Errorf("getting user from repo: %w", err)
	}

	if u == nil {
		text := "Добро пожаловать!\nДля работы необходимо зарегистрироваться"
		opts := menu.NewInline(buttons.Registration, buttons.Home, buttons.Back)

		if err = t.sender.Send(msg.UserID(), text, opts); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		return nil
	}
	//
	// if u.Type == user.TypeSupplier {
	//	return t.handlerSupplierHome(ctx, msg)
	//}
	//
	// return t.handlerCustomerHome(ctx, msg)
	//
	return nil
}

//
// func (t *Telegram) handlerSupplierHome(ctx context.Context, msg models.Request) error {
//	text := "Главное меню"
//
//	opts := inline.NewLink(
//		inline.Row(
//			inline.Text( buttons.NewLink()  "👀Показать мои товары", callback.callback{buttons.KeyEndpoint: "/show_goods"}),
//			inline.Text( buttons.NewLink() "➕Добавить товары", callback.callback{buttons.KeyEndpoint: "/post_goods"}),
//		),
//	)
//
//	if _, err := t.bot.Send(msg.chatID, text, opts); err != nil {
//		return fmt.Errorf("sending message: %w", err)
//	}
//
//	return nil
//}
//
// func (t *Telegram) handlerCustomerHome(ctx context.Context, msg models.Request) error {
//	text := "Главное меню"
//
//	opts := inline.NewLink(
//		inline.Row(
//			inline.Text("Показать товары", callback.callback{buttons.KeyEndpoint: "/show_goods"}),
//		),
//	)
//
//	if _, err := t.bot.Send(msg.chatID, text, opts); err != nil {
//		return fmt.Errorf("sending message: %w", err)
//	}
//
//	return nil
//}
