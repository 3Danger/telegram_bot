package telegram

import (
	"context"
	"fmt"
	"regexp"

	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/repo/user"
)

// State представляет текущее состояние пользователя в процессе авторизации
const (
	StateNone            = ""
	StateWaitingForName  = "waiting for name"
	StateWaitingForPhone = "waiting for phone"

	// Новые команды навигации
	cmdBack      = "⬅️ Назад"
	cmdConfirm   = "✅ Подтвердить"
	cmdEditName  = "📝 Изменить ФИО"
	cmdEditPhone = "📱 Изменить телефон"
)

// handlerAuth обрабатывает процесс авторизации
func (t *Telegram) handlerAuth(c tele.Context) error {
	ctx := getContext(c)

	// Обработка команды "Назад"
	if c.Text() == cmdBack {
		return t.handleBackNavigation(ctx, c)
	}

	// Обработка команд редактирования
	switch c.Text() {
	case cmdEditName:
		return t.handleEditName(ctx, c)
	case cmdEditPhone:
		return t.handleEditPhone(ctx, c)
	case cmdConfirm:
		return t.handleConfirmRegistration(ctx, c)
	}

	// Остальной код остается без изменений
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

	// Создаем кнопку для отправки контакта
	contactButton := &tele.ReplyMarkup{ResizeKeyboard: true}
	contactButton.Reply(
		contactButton.Row(
			contactButton.Contact("Отправить номер телефона"),
		),
	)

	if err = t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForPhone); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	return c.Send("Теперь, пожалуйста, поделитесь вашим номером телефона", contactButton)
}

func (t *Telegram) authStateWaitingPhone(ctx context.Context, c tele.Context) error {
	if c.Message().Contact == nil {
		return c.Send("Пожалуйста, используйте кнопку 'Отправить номер телефона'", createMenu(cmdBack))
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
		[]string{cmdConfirm},
		[]string{cmdEditName, cmdEditPhone},
	)

	return c.Send(msg, buttons)
}
func (t *Telegram) handleBackNavigation(ctx context.Context, c tele.Context) error {
	state, err := t.repo.state.Get(ctx, c.Sender().ID)
	if err != nil {
		return fmt.Errorf("get user state: %w", err)
	}

	switch state {
	case StateWaitingForPhone:
		if err := t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForName); err != nil {
			return fmt.Errorf("set user state: %w", err)
		}
		return c.Send("Пожалуйста, введите ваше ФИО заново", createMenu(cmdBack))

	case StateWaitingForName:
		if err := t.repo.state.Set(ctx, c.Sender().ID, StateNone); err != nil {
			return fmt.Errorf("set user state: %w", err)
		}
		return c.Send("Регистрация отменена", createMenu(home))
	}

	return nil
}

func (t *Telegram) handleEditName(ctx context.Context, c tele.Context) error {
	if err := t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForName); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}
	return c.Send("Пожалуйста, введите ваше ФИО заново", createMenu(cmdBack))
}

func (t *Telegram) handleEditPhone(ctx context.Context, c tele.Context) error {
	if err := t.repo.state.Set(ctx, c.Sender().ID, StateWaitingForPhone); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	contactButton := &tele.ReplyMarkup{ResizeKeyboard: true}
	contactButton.Reply(
		contactButton.Row(contactButton.Contact("Отправить номер телефона")),
		contactButton.Row(contactButton.Text(cmdBack)),
	)

	return c.Send("Пожалуйста, поделитесь вашим номером телефона заново", contactButton)
}

func (t *Telegram) handleConfirmRegistration(ctx context.Context, c tele.Context) error {
	if err := t.repo.state.Set(ctx, c.Sender().ID, StateNone); err != nil {
		return fmt.Errorf("set user state: %w", err)
	}

	return c.Send("Регистрация успешно завершена!", createMenu(home))
}
