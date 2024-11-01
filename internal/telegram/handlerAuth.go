package telegram

import (
	"context"
	"fmt"
	"regexp"

	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/telegram/constants"
)

// State представляет текущее состояние пользователя в процессе авторизации
const (
	StateNone            = ""
	StateWaitingForName  = "waiting for name"
	StateWaitingForPhone = "waiting for phone"
)

// handlerAuth обрабатывает процесс авторизации
func (t *Telegram) handlerAuth(c tele.Context) error {
	ctx := getContext(c)

	state, err := t.repo.state.Get(ctx, c.Sender().ID)
	if err != nil {
		return fmt.Errorf("get user state: %w", err)
	}

	switch state {
	case StateNone:
		return t.authStateNone(ctx, c)
	case StateWaitingForName:
		return t.authStateWaitingName(ctx, c)
	case StateWaitingForPhone:
		return t.authStateWaitingPhone(ctx, c)
	}

	return nil
}

func (t *Telegram) authStateNone(ctx context.Context, c tele.Context) error {
	currentUser, err := t.repo.user.User(ctx, c.Sender().ID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if currentUser != nil {
		if err = t.repo.state.Set(ctx, c.Sender().ID, StateNone); err != nil {
			return fmt.Errorf("set user state: %w", err)
		}

		return c.Send("Вы уже зарегистрированы")
	}

	if err = t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForName); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	return c.Send("Пожалуйста, введите ваше ФИО")
}

func (t *Telegram) authStateWaitingName(ctx context.Context, c tele.Context) error {
	reg, err := regexp.Compile("[a-zA-Zа-яА-Я]*")
	if err != nil {
		return fmt.Errorf("compile regular expression: %w", err)
	}

	names := reg.FindAllString(c.Text(), -1)

	if len(names) != 3 {
		return c.Send("Вы ввели неверные данные, пожалуйста повторите")
	}

	newUser := user.User{
		ID:         c.Sender().ID,
		IsSupplier: false,
		FirstName:  names[0],
		LastName:   names[1],
		Surname:    names[2],
	}

	if err = t.repo.user.CreateUser(ctx, newUser); err != nil {
		return fmt.Errorf("update user name: %w", err)
	}
	if err = t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForPhone); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	// Создаем кнопку для отправки контакта
	contactButton := &tele.ReplyMarkup{ResizeKeyboard: true}
	contactButton.Reply(
		contactButton.Row(
			contactButton.Contact("Отправить номер телефона"),
		),
	)

	return c.Send("Теперь, пожалуйста, поделитесь вашим номером телефона", contactButton)
}

func (t *Telegram) authStateWaitingPhone(ctx context.Context, c tele.Context) error {
	if c.Message().Contact == nil {
		return c.Send("Пожалуйста, используйте кнопку 'Отправить номер телефона'", createMenu(constants.Back))
	}

	phone := c.Message().Contact.PhoneNumber
	if err := t.repo.user.UpdateUserContactPhone(ctx, c.Sender().ID, phone); err != nil {
		return fmt.Errorf("update user phone: %w", err)
	}

	u, err := t.repo.user.User(ctx, c.Sender().ID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	if u == nil {
		return user.ErrUserNotFound
	}

	msg := fmt.Sprintf(`Проверьте ваши данные:
Имя: %s
Фамилия: %s
Отчество: %s
Телефон: %s

Если данные верны, нажмите "Подтвердить".
Для исправления данных используйте соответствующие кнопки.`,
		u.FirstName, u.LastName, u.Surname, u.Phone)

	buttons := createBigMenu(
		[]string{constants.AuthConfirm},
		[]string{constants.AuthEditName, constants.AuthEditPhone},
	)

	return c.Send(msg, buttons)
}
func (t *Telegram) handlerBackNavigation(c tele.Context) error {
	ctx := getContext(c)

	state, err := t.repo.state.Get(ctx, c.Sender().ID)
	if err != nil {
		return fmt.Errorf("get user state: %w", err)
	}

	switch state {
	case StateWaitingForPhone:
		if err := t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForName); err != nil {
			return fmt.Errorf("set user state: %w", err)
		}
		return c.Send("Пожалуйста, введите ваше ФИО заново", createMenu(constants.Back))

	case StateWaitingForName:
		if err := t.repo.state.Set(ctx, c.Sender().ID, StateNone); err != nil {
			return fmt.Errorf("set user state: %w", err)
		}
		return c.Send("Регистрация отменена", createMenu(constants.Home))
	}

	return nil
}

func (t *Telegram) handlerEditName(c tele.Context) error {
	if err := t.repo.state.Set(getContext(c), c.Sender().ID, StateWaitingForName); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}
	return c.Send("Пожалуйста, введите ваше ФИО заново", createMenu(constants.Back))
}

func (t *Telegram) handlerEditPhone(c tele.Context) error {
	if err := t.repo.state.Set(getContext(c), c.Sender().ID, StateWaitingForPhone); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	contactButton := &tele.ReplyMarkup{ResizeKeyboard: true}
	contactButton.Reply(
		contactButton.Row(contactButton.Contact("Отправить номер телефона")),
		contactButton.Row(contactButton.Text(constants.Back)),
	)

	return c.Send("Пожалуйста, поделитесь вашим номером телефона заново", contactButton)
}

func (t *Telegram) handlerConfirmRegistration(c tele.Context) error {
	if err := t.repo.state.Delete(getContext(c), c.Sender().ID); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	return c.Send("Регистрация успешно завершена!", createMenu(constants.Home))
}
