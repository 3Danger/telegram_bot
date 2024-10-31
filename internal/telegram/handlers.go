package telegram

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

const (
	home       = "🏠На главную"
	auth       = "⚙️Регистрация"
	imSupplier = "📦Я продавец"
	imCustomer = "💝Я покупатель"
)

func (t *Telegram) handlerHome(c tele.Context) error {
	u, err := t.repo.user.User(getContext(c), c.Sender().ID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}
	if u == nil {
		return c.Send(
			"Добро пожаловать!\nДля пользования необходимо регистрация",
			createMenu(auth),
		)
	}

	return c.Reply("", createMenu(home))
}

func createRow(buttonsStr []string) tele.Row {
	btns := make([]tele.Btn, len(buttonsStr))
	for _, item := range buttonsStr {
		btns = append(btns, tele.Btn{Text: item})
	}

	return btns
}

func createMenu(buttonsStr ...string) *tele.ReplyMarkup {
	menu := tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(createRow(buttonsStr))

	return &menu
}

func createBigMenu(buttonsStr ...[]string) *tele.ReplyMarkup {
	menu := tele.ReplyMarkup{ResizeKeyboard: true}

	rows := make([]tele.Row, 0, len(buttonsStr))
	for _, row := range buttonsStr {
		rows = append(rows, createRow(row))
	}

	menu.Reply(rows...)

	return &menu
}
