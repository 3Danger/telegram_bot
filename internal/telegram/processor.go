package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/samber/lo"

	"github.com/3Danger/telegram_bot/internal/telegram/handlers/auth"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/menu"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

func (t *Telegram) configureRoutes() {
	t.router[buttons.Registration.Button().Url] = auth.NewAuth(t.repo.user, t.sender)
}

func (t *Telegram) updateProcessor(ctx context.Context, update tele.Update) error {
	fmt.Println(string(lo.Must(json.MarshalIndent(update, "", "\t"))))

	msg := models.NewRequest(update)

	if err := t.messageProcessor(ctx, msg.Endpoint(), msg); err != nil {
		return fmt.Errorf("handle callback query: %w", err)
	}

	return nil
}

func (t *Telegram) messageProcessor(ctx context.Context, endpoint string, msg models.Request) error {
	switch endpoint {
	case buttons.Back.Button().Url, "":
		// STEP-BACK HANDLER
		key, err := t.repo.chain.Pop(ctx, msg.UserID())
		if err != nil {
			return fmt.Errorf("back-stepping: %w", err)
		}
		if key != "" {
			return t.messageProcessor(ctx, key, msg)
		}

		fallthrough
	case buttons.Home.Button().Url:
		// HOME HANDLER
		if err := t.handlerHome(ctx, msg); err != nil {
			return fmt.Errorf("handle %s: %w", endpoint, err)
		}

		if err := t.repo.chain.Clear(ctx, msg.UserID()); err != nil {
			return fmt.Errorf("cleaning states: %w", err)
		}

		return nil
	default:
		// BY ROUTE HANDLER
		handler, ok := t.router[endpoint]
		if !ok {
			return t.messageProcessor(ctx, buttons.Back.Button().Url, msg)
		}
		if err := handler.Handle(ctx, msg); err != nil {
			return fmt.Errorf("handle %s: %w", endpoint, err)
		}
		if err := t.repo.chain.Push(ctx, msg.UserID(), endpoint); err != nil {
			return fmt.Errorf("saving command story: %w", err)
		}

		return nil
	}

	if err := t.handlerUndefined(msg); err != nil {
		return fmt.Errorf("handling undefined message: %w", err)
	}

	return nil

}

func (t *Telegram) handlerUndefined(msg models.Request) error {
	text := "Вы просите странного"
	opt := menu.NewInline(buttons.Home)

	if err := t.sender.Send(msg.ChatID(), text, opt); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
