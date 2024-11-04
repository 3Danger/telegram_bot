package inline

import tele "github.com/PaulSonOfLars/gotgbot/v2"

func SendMessageOpts(buttons ...[]tele.InlineKeyboardButton) *tele.SendMessageOpts {
	return &tele.SendMessageOpts{
		ReplyMarkup: tele.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	}
}

func Row(btns ...tele.InlineKeyboardButton) []tele.InlineKeyboardButton {
	return btns
}
