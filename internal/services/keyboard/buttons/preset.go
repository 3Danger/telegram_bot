//nolint:gochecknoglobals
package buttons

import (
	"github.com/3Danger/telegram_bot/internal/services/keyboard"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/reply"
)

const (
	ConstHome         = "Домой"
	ConstBack         = "Назад"
	ConstRegistration = "Регистрация"
	ConstLocation     = "Поделиться локацией"
	ConstContact      = "Поделиться контактами"

	ConstAuthChoiceImSupplier = "Я продавец"
	ConstAuthChoiceImCustomer = "Я покупатель"
	ConstAuthSave             = "Сохранить"
	ConstAuthEdit             = "Изменить"
)

var (
	Home         keyboard.InlineButton = inline.NewWithEndpoint(ConstHome, "/start")
	Back         keyboard.InlineButton = inline.NewWithEndpoint(ConstBack, "/back")
	Registration keyboard.InlineButton = inline.NewWithEndpoint(ConstRegistration, "/registration")
	Location     keyboard.ReplyButton  = reply.New(ConstLocation).WithLocation()
	Contact      keyboard.ReplyButton  = reply.New(ConstContact).WithLocation()
)
