package reply

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/buttons"
)

var (
	ButtonAuth              = Text(buttons.EndpointRegistration)
	ButtonHome              = Text(buttons.EndpointHome)
	ButtonCustomerShowItems = Text(buttons.ButtonCustomerShowItems)
	ButtonSupplierShowItems = Text(buttons.ButtonSupplierShowItems)
	ButtonSupplierPostItems = Text(buttons.ButtonSupplierPostItems)
)

func Text(text string) tele.KeyboardButton {
	return tele.KeyboardButton{
		Text: text,
	}
}
