package keyboard

import "github.com/PaulSonOfLars/gotgbot/v2"

type Menu interface {
	Menu() *gotgbot.SendMessageOpts
}

type InlineButton interface {
	Button() gotgbot.InlineKeyboardButton
	Endpoint() string
}

type ReplyButton interface {
	Button() gotgbot.KeyboardButton
}
