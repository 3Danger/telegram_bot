package inline

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/models"
)

func Text(text string, callbackMap models.Pair) tele.InlineKeyboardButton {
	return tele.InlineKeyboardButton{
		Text:         text,
		CallbackData: models.NewCallback(callbackMap),
	}
}
