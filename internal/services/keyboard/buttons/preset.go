//nolint:gochecknoglobals
package buttons

import (
	"github.com/3Danger/telegram_bot/internal/services/keyboard"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/inline"
	"github.com/3Danger/telegram_bot/internal/services/keyboard/buttons/reply"
)

type Button string

func (b Button) String() string {
	return string(b)
}

func (b Button) InlineEndpoint(url string) *inline.Button {
	return inline.NewWithEndpoint(string(b), url)
}

func (b Button) Inline() *inline.Button {
	return inline.New(b.String())
}

func (b Button) Reply() *reply.Button {
	return reply.New(b.String())
}

const (
	ConstHome         Button = "Домой"
	ConstBack         Button = "Назад"
	ConstRegistration Button = "Регистрация"
	ConstLocation     Button = "Поделиться локацией"
	ConstContact      Button = "Поделиться контактами"

	ConstChangeName       Button = "Изменить имя"
	ConstChangeFamily     Button = "Изменить фамилию"
	ConstChangePhone      Button = "Изменить телефон"
	ConstChangeUserType   Button = "Изменить тип аккаунта"
	ConstChangeAdditional Button = "Править дополнительную информацию"

	ConstAuthChoiceImSupplier Button = "Я продавец"
	ConstAuthChoiceImCustomer Button = "Я покупатель"
	ConstAuthSave             Button = "Сохранить"
	ConstAuthEdit             Button = "Изменить"
)

const (
	URLHome         = "/start"
	URLBack         = "/back"
	URLRegistration = "/registration"
)

var (
	Home         keyboard.InlineButton = ConstHome.InlineEndpoint(URLHome)
	Back         keyboard.InlineButton = ConstBack.InlineEndpoint(URLBack)
	Registration keyboard.InlineButton = ConstRegistration.InlineEndpoint(URLRegistration)
	Location     keyboard.ReplyButton  = ConstLocation.Reply().WithLocation()
	Contact      keyboard.ReplyButton  = ConstContact.Reply().WithContact()
)
