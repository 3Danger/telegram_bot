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
		opt := reply.SendMessageOpts(reply.ButtonHome)

		if _, err = t.bot.SendMessage(data.ChatID, replyMessage, opt); err != nil {
			return fmt.Errorf("sending message: %w", err)
		}

		return nil
	}

	var replyMessage string
	var opt *tele.SendMessageOpts
	for range 2 {
		if replyMessage, opt, err = t.authProcessing(ctx, &data); err != nil {
			return fmt.Errorf("processing auth: %w", err)
		}
	}

	if _, err = t.bot.SendMessage(data.ChatID, replyMessage, opt); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (t *Telegram) authProcessing(ctx context.Context, data *models.Data) (replyMessage string, opt *tele.SendMessageOpts, err error) {
	replyMessage = "Вы ввели неверные данные, пожалуйста повторите"
	callbackEndpoint := models.PairKeyValues{
		Key:   buttons.KeyEndpoint,
		Value: buttons.EndpointRegistration,
	}

	state, err := t.svc.auth.GetSessionState(ctx, data.UserID)
	if err != nil {
		err = fmt.Errorf("getting from session: %w", err)
		return
	}

state:
	switch state {
	case auth.StateNew:
		replyMessage = "Добро пожаловать!"

		if err = t.svc.auth.AddUserID(ctx, data.UserID); err != nil {
			err = fmt.Errorf("adding user id: %w", err)
			return
		}
	case auth.StateType:
		if len(data.CallbackMap) == 0 || !user.Type(data.CallbackMap[buttons.UserType]).Valid() {
			replyMessage = "Выберите пожалуйста тип аккаунта"

			opt = inline.SendMessageOpts(
				inline.Text(
					"Я покупатель",
					models.Pair(buttons.UserType, user.TypeCustomer.String()),
					callbackEndpoint,
				),
				inline.Text(
					"Я продавец",
					models.Pair(buttons.UserType, user.TypeSupplier.String()),
					callbackEndpoint,
				),
			)
			break
		}

		if err = t.svc.auth.AddUserType(ctx, data.UserID, data.CallbackMap[buttons.UserType]); err != nil {
			err = fmt.Errorf("adding user type: %w", err)
			return
		}

	case auth.StateFirstName:
		if data.Message == "" {
			replyMessage = "Напишите пожалуйста ваше имя"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddFirstName(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StateLastName:
		if data.Message == "" {
			replyMessage = "Напишите пожалуйста вашу фамилию"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddLastName(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StateSurname:
		if data.Message == "" {
			replyMessage = "Напишите пожалуйста вашу отчество"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddSurname(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StatePhone:
		if data.Message == "" {
			replyMessage = "Укажите пожалуйста ваш номер телефона"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddPhone(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StateWhatsapp:
		if data.Message == "" {
			replyMessage = "Укажите пожалуйста ваш номер whatsapp"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddWhatsapp(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StateTelegram:
		if data.Message == "" {
			replyMessage = "Укажите пожалуйста ваш телеграм ник"
			opt = reply.SendMessageOpts(reply.ButtonHome)
			break
		}

		if err = t.svc.auth.AddTelegram(ctx, data.UserID, data.Message); err != nil {
			err = fmt.Errorf("adding first name: %w", err)
			return
		}

		data.Message = ""
	case auth.StateCompleted:
		var u *user.User
		u, err = t.svc.auth.GetFromSession(ctx, data.UserID)
		if err != nil {
			err = fmt.Errorf("getting from session: %w", err)
			return
		}
		if u == nil {
			err = errors.New("user not found")
			return
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
			break state
		}

		if data.CallbackMap[buttons.Change] != "" {
			delete(data.CallbackMap, buttons.Change)
			err = t.svc.auth.SetSomething(ctx, data.UserID, func(old *user.User) error {
				*old = *u

				return nil
			})
			if err != nil {
				err = fmt.Errorf("changing user data: %w", err)
				return
			}
		}

		answer := data.CallbackMap[buttons.Confirm]
		switch answer {
		case "yes":
		case "no":
			opt = inline.SendMessageOpts(
				inline.Text("Изменить фамилию", models.Pair(buttons.Change, "first_name"), callbackEndpoint),
				inline.Text("Изменить имя: %s\n", models.Pair(buttons.Change, "last_name"), callbackEndpoint),
				inline.Text("Изменить отчество", models.Pair(buttons.Change, "surname"), callbackEndpoint),
				inline.Text("Изменить тип аккаунта", models.Pair(buttons.Change, "user_type"), callbackEndpoint),
				inline.Text("Изменить телефон", models.Pair(buttons.Change, "phone"), callbackEndpoint),
				inline.Text("Изменить whatsapp", models.Pair(buttons.Change, "whatsapp"), callbackEndpoint),
				inline.Text("Изменить telegram", models.Pair(buttons.Change, "telegram"), callbackEndpoint),
			)

		default:
			u, err = t.svc.auth.GetFromSession(ctx, data.UserID)
			if err != nil {
				err = fmt.Errorf("getting from session: %w", err)
				return
			}
			if u == nil {
				err = fmt.Errorf("user not found")
				return
			}

			typeUser := func(t user.Type) string {
				switch t {
				case user.TypeCustomer:
					return "Я покупатель"
				case user.TypeSupplier:
					return "Я продавец"
				default:
					return ""
				}
			}

			replyMessage = "Проверьте пожалуйста ваши данные\n"
			replyMessage += fmt.Sprintf("Фамилие: %s\n", u.LastName)
			replyMessage += fmt.Sprintf("Имя: %s\n", u.FirstName)
			replyMessage += fmt.Sprintf("Отчество: %s\n", u.FirstName)
			replyMessage += fmt.Sprintf("Тип аккаунта: %s\n", typeUser(u.Type))
			replyMessage += fmt.Sprintf("Телефон: %s\n", u.FirstName)
			replyMessage += fmt.Sprintf("WhatsApp: %s\n", u.Whatsapp)
			replyMessage += fmt.Sprintf("Telegram: %s\n", u.Telegram)

			opt = inline.SendMessageOpts(
				inline.Text(
					"Все верно",
					models.Pair(buttons.Confirm, "yes"),
					callbackEndpoint,
				),
				inline.Text(
					"Изменить",
					models.Pair(buttons.Confirm, "no"),
					callbackEndpoint,
				),
			)
		}
	}

	return
}
