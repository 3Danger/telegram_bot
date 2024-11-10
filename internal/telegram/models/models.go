package models

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons"
	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons/callback"
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
	UserId      int64
}

func NewRequest(update tele.Update) Request {
	msg := Request{}

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
				UserId:      contact.UserId,
			}
		}

		msg.chatID = um.Chat.Id
		if um.From != nil {
			msg.userID = um.From.Id
		}
		if um.Text == buttons.Home.Button().Url {
			msg.callback.SetEndpoint(buttons.Home.Button().Url)
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

func (r *Request) UserID() int64 {
	return r.userID
}

func (r *Request) ChatID() int64 {
	return r.chatID
}

func (r *Request) Message() string {
	return r.message
}
