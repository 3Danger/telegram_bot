//nolint:exhaustruct
package auth

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/3Danger/telegram_bot/internal/models"
	r "github.com/3Danger/telegram_bot/internal/repo"
	session "github.com/3Danger/telegram_bot/internal/repo/session/inmemory"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/menu"
)

type stateType int

const (
	stateWelcome stateType = iota
	stateContact
	stateUserType
	stateComplete
)

type repo struct {
	state r.Repo[stateType]
	cache r.Repo[*userpg.User]
	user  userpg.Querier
}

type subHandler func(ctx context.Context, a *Auth, data models.Request) (models.Responses, bool, error)

type Auth struct {
	repo          repo
	subHandlerMap map[stateType]subHandler
}

const cacheSize = 10000

func NewAuth(users userpg.Querier) *Auth {
	return &Auth{
		repo: repo{
			state: session.NewRepo[stateType](cacheSize),
			cache: session.NewRepo[*userpg.User](cacheSize),
			user:  users,
		},
		subHandlerMap: map[stateType]subHandler{
			stateWelcome:  subHandlerWelcome,
			stateContact:  subHandlerContact,
			stateUserType: subHandlerUserType,
			// stateComplete: subHandlerComplete,
		},
	}
}

func (a *Auth) Handle(ctx context.Context, data models.Request) (models.Responses, error) {
	u, err := a.repo.user.Get(ctx, data.UserID())
	if err != nil {
		return nil, errors.Wrap(err, "getting user from temporary repo")
	}

	if u != nil {
		return models.NewResponses(
			data.ChatID(),
			"Вы уже зарегистрированны!",
			menu.NewInline(buttons.Home, buttons.Back),
		), nil
	}

	state, err := a.repo.state.Get(ctx, data.UserID())
	if err != nil {
		return nil, fmt.Errorf("getting session stateUserType: %w", err)
	}

	var response models.Responses

	for {
		handler, ok := a.subHandlerMap[state]
		if !ok {
			break
		}

		var next bool

		var resp models.Responses

		resp, next, err = handler(ctx, a, data)
		if err != nil {
			return nil, fmt.Errorf("handling session: %w", err)
		}

		response = append(response, resp...)

		if !next {
			break
		}

		state++
	}

	if err = a.repo.state.Set(ctx, data.UserID(), state); err != nil {
		return nil, fmt.Errorf("getting session stateUserType: %w", err)
	}

	return response, nil
}

func subHandlerWelcome(_ context.Context, _ *Auth, data models.Request) (models.Responses, bool, error) {
	return models.NewResponses(data.ChatID(),
		"Добро пожаловать на страницу регистрации\nНажми поделиться контактами что бы продолжить!",
		menu.NewReply(buttons.Contact).OneTime(true).Persistent(true),
	), true, nil
}

func subHandlerContact(ctx context.Context, a *Auth, data models.Request) (models.Responses, bool, error) {
	contact := data.Contact()
	if contact == nil {
		return models.NewResponses(data.ChatID(),
			"Поделитесь контактами что бы продолжить",
			menu.NewReply(buttons.Contact).OneTime(true).Persistent(true),
		), false, nil
	}

	newUser, err := a.repo.cache.Get(ctx, contact.UserID)
	if err != nil {
		return nil, false, errors.Wrap(err, "getting user from cache")
	}

	if newUser == nil {
		newUser = new(userpg.User)
	}

	newUser.ID = contact.UserID
	newUser.FirstName = contact.LastName
	newUser.LastName = contact.FirstName
	newUser.Phone = contact.PhoneNumber

	if err = a.repo.cache.Set(ctx, contact.UserID, newUser); err != nil {
		return nil, false, fmt.Errorf("setting user to cache: %w", err)
	}

	return nil, true, nil
}

func subHandlerUserType(ctx context.Context, a *Auth, data models.Request) (models.Responses, bool, error) {
	const (
		supplier = "supplier"
		customer = "customer"
		userType = "user_type"
	)

	u, err := a.repo.cache.Get(ctx, data.UserID())
	if err != nil {
		return nil, false, fmt.Errorf("getting user from cache: %w", err)
	}

	if u == nil {
		u = &userpg.User{ID: data.UserID()}
	}

	switch data.Value(userType) {
	case supplier:
		u.UserType = userpg.UserTypeSupplier
	case customer:
		u.UserType = userpg.UserTypeCustomer
	default:
		return models.NewResponses(
			data.ChatID(),
			"Выберите тип аккаунта",
			menu.NewInline(
				inline.New(buttons.ConstAuthChoiceImCustomer).WithValue(userType, customer),
				inline.New(buttons.ConstAuthChoiceImSupplier).WithValue(userType, supplier),
			),
		), false, nil
	}

	if err = a.repo.cache.Set(ctx, data.UserID(), u); err != nil {
		return nil, false, fmt.Errorf("setting user to cache: %w", err)
	}

	return nil, true, nil
}

//nolint:gocritic,dupword
//func subHandlerComplete(_ context.Context, _ *Auth, _ models.Request) (response, error) {
//	const (
//		ask     = "ask"
//		confirm = "confirm"
//		edit    = "edit"
//	)
//	resp := response{
//		msg: "Проверьте свои свои данные",
//		opt: menu.NewInline(
//			inline.New(buttons.ConstAuthSave).WithValue(ask, confirm),
//			inline.New(buttons.ConstAuthEdit).WithValue(ask, edit),
//		),
//	}
// ...
//return resp, nil
//}
