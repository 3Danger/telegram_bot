package telegram

import tele "gopkg.in/telebot.v4"

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
