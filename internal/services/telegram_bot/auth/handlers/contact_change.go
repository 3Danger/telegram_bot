package handlers

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/menu"
)

type ChangeData struct {
	Next Handler
}

func (c ChangeData) Name() string { return "change_data" }

func (c ChangeData) Process(ctx context.Context, r Repo, data models.Request) (models.Responses, error) {
	const (
		askName       = "name"
		askFamily     = "family"
		askPhone      = "phone"
		askUserType   = "user_type"
		askAdditional = "additional"
	)

	user, err := r.Get(ctx, data.UserID())
	if err != nil {
		return nil, fmt.Errorf("getting user from cache %w", err)
	}

	if user == nil {
		user = new(models.User)
		user.ID = data.UserID()
	}

	answer := data.Value(c.Name())
	if answer == "" {
		text := fmt.Sprintf(`
			Выберите что изменить:
			Имя: %s
			Фамилия: %s
			Телефон: %s
			Тим аккаунта: %s
			Дополнительная информация: %s`,
			user.FirstName, user.LastName, user.Phone, user.UserType, user.Additional)

		return models.Responses{{
			ChatID: data.ChatID(),
			Text:   text,
			Menu: menu.NewInline(
				buttons.ConstChangeName.Inline().WithValue(c.Name(), askName),
			).Add(
				buttons.ConstChangeFamily.Inline().WithValue(c.Name(), askFamily),
			).Add(
				buttons.ConstChangePhone.Inline().WithValue(c.Name(), askPhone),
			).Add(
				buttons.ConstChangeUserType.Inline().WithValue(c.Name(), askUserType),
			).Add(
				buttons.ConstChangeAdditional.Inline().WithValue(c.Name(), askAdditional),
			),
		}}, nil
	}

	switch answer {
	case askName:
		user.FirstName = data.Message()
	case askFamily:
		user.LastName = data.Message()
	case askPhone:
		user.Phone = data.Message()
	//case askUserType:
	//	user.UserType = data.Message()
	case askAdditional:
		user.Additional = data.Message()
	default:
		return nil, fmt.Errorf("unknown answer %s", answer)
	}

	if err = r.Set(ctx, data.UserID(), user); err != nil {
		return nil, fmt.Errorf("setting user to cache: %w", err)
	}

	return nil, nil
}
