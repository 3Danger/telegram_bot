package handlers

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/menu"
)

type ContactFromTg struct {
	Next Handler
}

func (c ContactFromTg) Name() string { return "contact_from_tg" }

func (c ContactFromTg) Process(ctx context.Context, r Repo, data models.Request) (models.Responses, error) {
	newUser, err := r.Get(ctx, data.UserID())
	if err != nil {
		return nil, fmt.Errorf("getting user from cache %w", err)
	}

	if newUser != nil &&
		newUser.ID != 0 &&
		newUser.FirstName != "" &&
		newUser.LastName != "" &&
		newUser.Phone != "" {
		return c.Next.Process(ctx, r, data)
	}

	contact := data.Contact()
	if contact == nil {
		return models.Responses{{
			data.ChatID(),
			"Поделитесь контактами что бы продолжить",
			menu.NewReply(
				buttons.ConstContact.Reply().WithContact(),
				buttons.ConstHome.Reply(),
				buttons.ConstBack.Reply(),
			).OneTime(true).Persistent(true),
		}}, nil
	}

	if newUser == nil {
		newUser = new(models.User)
	}

	newUser.ID = data.UserID()
	newUser.FirstName = contact.FirstName
	newUser.LastName = contact.LastName
	newUser.Phone = contact.PhoneNumber

	if err = r.Set(ctx, contact.UserID, newUser); err != nil {
		return nil, fmt.Errorf("setting user to cache: %w", err)
	}

	return c.Next.Process(ctx, r, data)
}
