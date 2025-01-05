package telegrambot

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/services/telegram_bot/handlers/auth"
)

func (t *Telegram) configureRoutes() {
	t.router[buttons.Registration.Button().Url] = auth.NewAuth(t.repo.user)
}

func (t *Telegram) MessageProcessor(
	ctx context.Context, msg models.Request,
) ([]models.Response, error) {
	resp, err := t.messageProcessor(ctx, msg.Endpoint(), msg)
	if err != nil {
		return nil, fmt.Errorf("handle callback query: %w", err)
	}

	return resp, nil
}

func (t *Telegram) messageProcessor(
	ctx context.Context, endpoint string, msg models.Request,
) ([]models.Response, error) {
	switch endpoint {
	case buttons.Back.Button().Url, "":
		// STEP-BACK HANDLER
		key, err := t.repo.chain.Pop(ctx, msg.UserID())
		if err != nil {
			return nil, fmt.Errorf("back-stepping: %w", err)
		}

		if key != "" {
			return t.messageProcessor(ctx, key, msg)
		}

		fallthrough
	case buttons.Home.Button().Url:
		// HOME HANDLER
		resp, err := t.handlerHome(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("handle %s: %w", endpoint, err)
		}

		if err := t.repo.chain.Clear(ctx, msg.UserID()); err != nil {
			return nil, fmt.Errorf("cleaning states: %w", err)
		}

		return resp, nil
	default:
		// BY ROUTE HANDLER
		handler, ok := t.router[endpoint]
		if !ok {
			return t.messageProcessor(ctx, buttons.Back.Button().Url, msg)
		}

		resp, err := handler.Handle(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("handle %s: %w", endpoint, err)
		}

		if err := t.repo.chain.Push(ctx, msg.UserID(), endpoint); err != nil {
			return nil, fmt.Errorf("saving command story: %w", err)
		}

		return resp, nil
	}
}
