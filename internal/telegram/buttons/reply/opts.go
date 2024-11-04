package reply

import tele "github.com/PaulSonOfLars/gotgbot/v2"

func SendMessageOpts(buttons ...tele.KeyboardButton) *tele.SendMessageOpts {
	return &tele.SendMessageOpts{
		ReplyMarkup: tele.ReplyKeyboardMarkup{
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
			Keyboard: [][]tele.KeyboardButton{
				buttons,
			},
		},
	}
}
