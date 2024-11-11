//nolint:exhaustruct
package auth

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	r "github.com/3Danger/telegram_bot/internal/repo"
	session "github.com/3Danger/telegram_bot/internal/repo/session/inmemory"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/menu"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
	"github.com/3Danger/telegram_bot/internal/telegram/sender"
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

type response struct {
	msg string
	opt keyboard.Menu
}

type subHandler func(ctx context.Context, a *Auth, data models.Request) (response, bool, error)

type Auth struct {
	repo          repo
	subHandlerMap map[stateType]subHandler
	sender        *sender.Sender
}

const cacheSize = 10000

func NewAuth(users userpg.Querier, sender *sender.Sender) *Auth {
	a := &Auth{
		repo: repo{
			state: session.NewRepo[stateType](cacheSize),
			cache: session.NewRepo[*userpg.User](cacheSize),
			user:  users,
		},
		sender: sender,
		subHandlerMap: map[stateType]subHandler{
			stateWelcome:  subHandlerWelcome,
			stateContact:  subHandlerContact,
			stateUserType: subHandlerUserType,
			// stateComplete: subHandlerComplete,
		},
	}

	return a
}

func (a *Auth) Handle(ctx context.Context, data models.Request) error {
	u, err := a.repo.user.Get(ctx, data.UserID())
	if err != nil {
		return errors.Wrap(err, "getting user from temporary repo")
	}

	if u != nil {
		msg := "Вы уже зарегистрированны!"
		opt := menu.NewInline(buttons.Home, buttons.Back)

		if err := a.sender.Send(data.ChatID(), msg, opt); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		return nil
	}

	state, err := a.repo.state.Get(ctx, data.UserID())
	if err != nil {
		return fmt.Errorf("getting session stateUserType: %w", err)
	}

	var resp response
	resp.opt = keyboard.Menu(menu.NewInline(buttons.Home, buttons.Back))

	for {
		handler, ok := a.subHandlerMap[state]
		if !ok {
			break
		}

		var next bool

		if resp, next, err = handler(ctx, a, data); err != nil {
			return fmt.Errorf("handling session: %w", err)
		}

		if err = a.sender.Send(data.ChatID(), resp.msg, resp.opt); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		if !next {
			break
		}

		state++
	}

	if err = a.repo.state.Set(ctx, data.UserID(), state); err != nil {
		return fmt.Errorf("getting session stateUserType: %w", err)
	}

	return nil
}

func subHandlerWelcome(_ context.Context, _ *Auth, _ models.Request) (response, bool, error) {
	return response{
		msg: "Добро пожаловать на страницу регистрации\nНажми поделиться контактами что бы продолжить!",
		opt: menu.NewReply(buttons.Contact).OneTime(true).Persistent(true),
	}, true, nil
}

func subHandlerContact(ctx context.Context, a *Auth, data models.Request) (response, bool, error) {
	contact := data.Contact()
	if contact == nil {
		return response{
			msg: "Поделитесь контактами что бы продолжить",
			opt: menu.NewReply(buttons.Contact).OneTime(true).Persistent(true),
		}, false, nil
	}

	newUser, err := a.repo.cache.Get(ctx, contact.UserID)
	if err != nil {
		return response{}, false, errors.Wrap(err, "getting user from cache")
	}

	if newUser == nil {
		newUser = new(userpg.User)
	}

	newUser.ID = contact.UserID
	newUser.FirstName = contact.LastName
	newUser.LastName = contact.FirstName
	newUser.Phone = contact.PhoneNumber

	if err = a.repo.cache.Set(ctx, contact.UserID, newUser); err != nil {
		return response{}, false, fmt.Errorf("setting user to cache: %w", err)
	}

	return response{}, true, nil
}

func subHandlerUserType(ctx context.Context, a *Auth, data models.Request) (response, bool, error) {
	const (
		supplier = "supplier"
		customer = "customer"
		userType = "user_type"
	)

	u, err := a.repo.cache.Get(ctx, data.UserID())
	if err != nil {
		return response{}, false, fmt.Errorf("getting user from cache: %w", err)
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
		return response{
			msg: "Выберите тип аккаунта",
			opt: menu.NewInline(
				inline.New(buttons.ConstAuthChoiceImCustomer).WithValue(userType, customer),
				inline.New(buttons.ConstAuthChoiceImSupplier).WithValue(userType, supplier),
			),
		}, false, nil
	}

	if err = a.repo.cache.Set(ctx, data.UserID(), u); err != nil {
		return response{}, false, fmt.Errorf("setting user to cache: %w", err)
	}

	return response{}, true, nil
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
