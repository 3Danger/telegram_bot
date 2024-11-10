package reply

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type Button struct {
	text     string
	contact  bool
	location bool
}

func New(text string) *Button {
	return &Button{
		text: text,
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
		RequestContact:  b.contact,
		RequestLocation: b.location,
	}
}
