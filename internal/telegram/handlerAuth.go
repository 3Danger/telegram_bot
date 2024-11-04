package telegram

import (
	"context"
	"errors"
	"fmt"

	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/services/auth"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/telegram/buttons/reply"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

// State представляет текущее состояние пользователя в процессе авторизации
const (
	StateNone       = ""
	StateCreateUser = "auth_create user"
	StateUserType   = "auth_user_type"
	StateFirstName  = "auth_first_name"
	StateLastName   = "auth_last_name"
	StateSurname    = "auth_surname"
	StatePhone      = "auth_phone"
	StateWhatsapp   = "auth_whatsapp"
	StateTelegram   = "auth_telegram"
	StateFinish     = "auth_finish"
)

// handlerAuth обрабатывает процесс авторизации
func (t *Telegram) handlerAuth(ctx context.Context, data models.Data) error {
	isRegisteredPermanent, err := t.svc.auth.RegisteredPermanent(ctx, data.UserID)
	if err != nil {
		return fmt.Errorf("checking register permanent: %w", err)
	}

	if isRegisteredPermanent {
		replyMessage := "Вы уже зарегистрированны"
		opt := reply.SendMessageOpts(reply.Row(reply.ButtonHome))

		if _, err = t.bot.SendMessage(data.ChatID, replyMessage, opt); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		return nil
	}

	return t.authProcessing(ctx, &data)
}

func (t *Telegram) authProcessing(ctx context.Context, data *models.Data) error {
	replyMessage := "Вы ввели неверные данные, пожалуйста повторите"

	pair := models.Pair{buttons.KeyEndpoint: buttons.EndpointRegistration}

	messageChan := make(chan string, 1)
	messageChan <- data.Message
	close(messageChan)
	defer func() {
		for range messageChan {
		}
	}()

	var opt *tele.SendMessageOpts

	stateNew := func() error {
		replyMessage = "Добро пожаловать!"

		if err := t.svc.auth.AddUserID(ctx, data.UserID); err != nil {
			return fmt.Errorf("adding user id: %w", err)
		}

		return nil
	}

	stateType := func() error {
		if len(data.CallbackMap) == 0 || !user.Type(data.CallbackMap[buttons.UserType]).Valid() {
			replyMessage = "Выберите пожалуйста тип аккаунта"

			opt = inline.SendMessageOpts(
				inline.Row(
					inline.Text("Я покупатель", pair.With(buttons.UserType, user.TypeCustomer.String())),
					inline.Text("Я продавец", pair.With(buttons.UserType, user.TypeSupplier.String())),
				),
			)

			return nil
		}

		if err := t.svc.auth.AddUserType(ctx, data.UserID, data.CallbackMap[buttons.UserType]); err != nil {
			return fmt.Errorf("adding user type: %w", err)
		}

		return nil
	}

	stateFirstName := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Напишите пожалуйста ваше имя"
			return nil
		}

		if err := t.svc.auth.AddFirstName(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding first name: %w", err)
		}

		return nil
	}
	stateLastName := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Напишите пожалуйста вашу фамилию"
			return nil
		}

		if err := t.svc.auth.AddLastName(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding last name: %w", err)
		}

		return nil
	}
	stateSurname := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Напишите пожалуйста вашу отчество"
			return nil
		}

		if err := t.svc.auth.AddSurname(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding surname: %w", err)
		}

		return nil
	}
	statePhone := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Укажите пожалуйста ваш номер телефона"
			return nil
		}

		if err := t.svc.auth.AddPhone(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding phone name: %w", err)
		}

		return nil
	}
	stateWhatsapp := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Укажите пожалуйста ваш номер whatsapp"
			return nil
		}

		if err := t.svc.auth.AddWhatsapp(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding whatsapp: %w", err)
		}

		return nil
	}
	stateTelegram := func() error {
		msg := <-messageChan
		if msg == "" {
			replyMessage = "Укажите пожалуйста ваш телеграм ник"
			return nil
		}

		if err := t.svc.auth.AddTelegram(ctx, data.UserID, msg); err != nil {
			return fmt.Errorf("adding telegram: %w", err)
		}

		return nil
	}

	stateCompleted := func() error {
		u, err := t.svc.auth.GetFromSession(ctx, data.UserID)
		if err != nil {
			return fmt.Errorf("getting from session: %w", err)
		}
		if u == nil {
			return errors.New("user not found")
		}

		switch data.CallbackMap[buttons.Change] {
		case "first_name":
			u.FirstName = ""
		case "last_name":
			u.LastName = ""
		case "surname":
			u.Surname = ""
		case "user_type":
			u.Type = ""
		case "phone":
			u.Phone = ""
		case "whatsapp":
			u.Whatsapp = ""
		case "telegram":
			u.Telegram = ""
		case "":
		default:
			return nil
		}

		if data.CallbackMap[buttons.Change] != "" {
			delete(data.CallbackMap, buttons.Change)
			err = t.svc.auth.SetSomething(ctx, data.UserID, func(old *user.User) error {
				*old = *u

				return nil
			})
			if err != nil {
				return fmt.Errorf("changing user data: %w", err)
			}

			return nil
		}

		answer := data.CallbackMap[buttons.Confirm]
		switch answer {
		case "yes":
			if err = t.svc.auth.SaveToPermanent(ctx, data.UserID); err != nil {
				return fmt.Errorf("saving answer: %w", err)
			}

			replyMessage = "Сохранено!"
			opt = reply.SendMessageOpts(reply.Row(reply.ButtonHome))

			return nil
		case "no":
			replyMessage = "Что изменить?"
			opt = inline.SendMessageOpts(
				inline.Row(inline.Text("Изменить фамилию", pair.With(buttons.Change, "first_name"))),
				inline.Row(inline.Text("Изменить имя", pair.With(buttons.Change, "last_name"))),
				inline.Row(inline.Text("Изменить отчество", pair.With(buttons.Change, "surname"))),
				inline.Row(inline.Text("Изменить тип аккаунта", pair.With(buttons.Change, "user_type"))),
				inline.Row(inline.Text("Изменить телефон", pair.With(buttons.Change, "phone"))),
				inline.Row(inline.Text("Изменить whatsapp", pair.With(buttons.Change, "whatsapp"))),
				inline.Row(inline.Text("Изменить telegram", pair.With(buttons.Change, "telegram"))),
			)
		default:
			opt = inline.SendMessageOpts(
				inline.Row(inline.Text("Сохранить", pair.With(buttons.Confirm, "yes"))),
				inline.Row(inline.Text("Изменить", pair.With(buttons.Confirm, "no"))),
			)
			replyMessage = "Проверьте ваши данные\n"
			replyMessage += fmt.Sprintf("Имя: %s\n", u.FirstName)
			replyMessage += fmt.Sprintf("Фамилия: %s\n", u.LastName)
			replyMessage += fmt.Sprintf("Отчество: %s\n", u.Surname)
			replyMessage += fmt.Sprintf("Тип профиля: %s\n", u.Type)
			replyMessage += fmt.Sprintf("Телефон: %s\n", u.Phone)
			replyMessage += fmt.Sprintf("Whatsapp: %s\n", u.Whatsapp)
			replyMessage += fmt.Sprintf("Telegram: %s\n", u.Telegram)
		}

		return nil
	}

	stateMap := map[auth.State]func() error{
		auth.StateNew:       stateNew,
		auth.StateType:      stateType,
		auth.StateFirstName: stateFirstName,
		auth.StateLastName:  stateLastName,
		auth.StateSurname:   stateSurname,
		auth.StatePhone:     statePhone,
		auth.StateWhatsapp:  stateWhatsapp,
		auth.StateTelegram:  stateTelegram,
		auth.StateCompleted: stateCompleted,
	}

	for range 2 {
		state, err := t.svc.auth.GetSessionState(ctx, data.UserID)
		if err != nil {
			return fmt.Errorf("getting from session: %w", err)
		}
		if fn, ok := stateMap[state]; ok {
			if err = fn(); err != nil {
				return err
			}

			continue
		}

	}

	if _, err := t.bot.SendMessage(data.ChatID, replyMessage, opt); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
