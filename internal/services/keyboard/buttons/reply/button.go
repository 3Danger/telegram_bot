package reply

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type Button struct {
	text     string
	contact  bool
	location bool
}

func New[T ~string](text T) *Button {
	return &Button{
		text:     string(text),
		contact:  false,
		location: false,
	}
}

func (b *Button) WithLocation() *Button {
	b.location = true

	return b
}

func (b *Button) WithContact() *Button {
	b.contact = true

	return b
}

func (b *Button) Button() tele.KeyboardButton {
	return tele.KeyboardButton{
		Text:            b.text,
		RequestUsers:    nil,
		RequestChat:     nil,
		RequestContact:  b.contact,
		RequestLocation: b.location,
		RequestPoll:     nil,
		WebApp:          nil,
	}
}
