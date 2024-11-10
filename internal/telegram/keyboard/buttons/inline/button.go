package inline

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard/buttons/callback"
)

type Button struct {
	text     string
	callback *callback.Callback
}

func New(text string) *Button {
	return &Button{
		text: text,
	}
}

func NewWithEndpoint(text, url string) *Button {
	return &Button{
		text:     text,
		callback: callback.New().With("endpoint", url),
	}
}

func (b *Button) WithValue(k, v string) *Button {
	b.callback = b.callback.With(k, v)

	return b
}

func (b *Button) Button() tele.InlineKeyboardButton {
	return tele.InlineKeyboardButton{
		Text:         b.text,
		CallbackData: b.callback.Data(),
	}
}
