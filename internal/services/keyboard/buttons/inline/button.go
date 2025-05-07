package inline

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/callback"
)

type Button struct {
	text     string
	callback *callback.Callback
}

func New(text string) *Button {
	return &Button{
		text:     text,
		callback: &callback.Callback{},
	}
}

func NewWithEndpoint[T ~string](text, url T) *Button {
	return &Button{
		text:     string(text),
		callback: callback.New().With("endpoint", string(url)),
	}
}

func (b *Button) WithValue(k, v string) *Button {
	b.callback = b.callback.With(k, v)

	return b
}

func (b *Button) Endpoint() string {
	return b.callback.Endpoint()
}

func (b *Button) Button() tele.InlineKeyboardButton {
	return tele.InlineKeyboardButton{
		Text:                         b.text,
		Url:                          "",
		CallbackData:                 b.callback.Data(),
		WebApp:                       nil,
		LoginUrl:                     nil,
		SwitchInlineQuery:            nil,
		SwitchInlineQueryCurrentChat: nil,
		SwitchInlineQueryChosenChat:  nil,
		CopyText:                     nil,
		CallbackGame:                 nil,
		Pay:                          false,
	}
}
