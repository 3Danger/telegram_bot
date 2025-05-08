package models

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/services/keyboard"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/callback"
)

type Request struct {
	userID   int64
	chatID   int64
	message  string
	callback *callback.Callback
	contact  *Contact
}

type Contact struct {
	PhoneNumber string
	FirstName   string
	LastName    string
	UserID      int
}

type Response struct {
	ChatID int
	Text   string
	Menu   keyboard.Menu
}

type Responses []Response

func NewResponses(chatID int, text string, menu keyboard.Menu) Responses {
	return Responses{
		Response{
			ChatID: chatID,
			Text:   text,
			Menu:   menu,
		},
	}
}

func (r *Responses) Add(chatID int, text string, menu keyboard.Menu) {
	*r = append(*r, Response{
		ChatID: chatID,
		Text:   text,
		Menu:   menu,
	})
}

func NewRequest(update tele.Update) Request {
	msg := Request{} //nolint:exhaustruct

	cbData := ""

	if cb := update.CallbackQuery; cb != nil {
		msg.chatID = cb.Message.GetChat().Id
		msg.userID = cb.From.Id
		cbData = cb.Data
	}

	msg.callback = callback.FromString(cbData)

	if um := update.Message; um != nil {
		if contact := um.Contact; contact != nil {
			msg.contact = &Contact{
				PhoneNumber: contact.PhoneNumber,
				FirstName:   contact.FirstName,
				LastName:    contact.LastName,
				UserID:      int(contact.UserId),
			}
		}

		msg.chatID = um.Chat.Id
		if um.From != nil {
			msg.userID = um.From.Id
		}

		if isHomePage := um.Text == buttons.Home.Button().Text; isHomePage {
			msg.callback.SetEndpoint(buttons.Home.Button().Text)

			um.Text = ""
		}

		msg.message = um.Text
	}

	return msg
}

func (r *Request) Contact() *Contact {
	return r.contact
}

func (r *Request) Endpoint() string {
	return r.callback.Endpoint()
}

func (r *Request) Value(key string) string {
	return r.callback.Value(key)
}

func (r *Request) UserID() int {
	return int(r.userID)
}

func (r *Request) ChatID() int {
	return int(r.chatID)
}

func (r *Request) Message() string {
	return r.message
}
