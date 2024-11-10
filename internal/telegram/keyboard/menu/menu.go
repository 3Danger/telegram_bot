package menu

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard"
)

type MenuInline struct {
	buttons [][]keyboard.InlineButton
}

func NewInline(row ...keyboard.InlineButton) *MenuInline {
	return &MenuInline{buttons: [][]keyboard.InlineButton{row}}
}

func (m *MenuInline) Add(row ...keyboard.InlineButton) *MenuInline {
	m.buttons = append(m.buttons, row)

	return m
}

func (m *MenuInline) Menu() *tele.SendMessageOpts {
	inlineKeyboard := make([][]tele.InlineKeyboardButton, len(m.buttons))
	for i, row := range m.buttons {
		inlineKeyboard[i] = make([]tele.InlineKeyboardButton, len(row))
		for j, button := range row {
			inlineKeyboard[i][j] = button.Button()
		}
	}
	return &tele.SendMessageOpts{
		ReplyMarkup: tele.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard},
	}
}

type MenuReply struct {
	buttons         [][]keyboard.ReplyButton
	oneTimeKeyboard bool
	persistent      bool
}

func NewReply(row ...keyboard.ReplyButton) *MenuReply {
	return &MenuReply{buttons: [][]keyboard.ReplyButton{row}}
}

func (m *MenuReply) Add(row ...keyboard.ReplyButton) *MenuReply {
	m.buttons = append(m.buttons, row)

	return m
}

func (m *MenuReply) OneTime(is bool) *MenuReply {
	m.oneTimeKeyboard = is
	return m
}

func (m *MenuReply) Persistent(is bool) *MenuReply {
	m.persistent = is
	return m
}

func (m *MenuReply) Menu() *tele.SendMessageOpts {
	replyKeyboard := make([][]tele.KeyboardButton, len(m.buttons))
	for i, row := range m.buttons {
		replyKeyboard[i] = make([]tele.KeyboardButton, len(row))
		for j, button := range row {
			replyKeyboard[i][j] = button.Button()
		}
	}
	return &tele.SendMessageOpts{
		ReplyMarkup: tele.ReplyKeyboardMarkup{
			Keyboard:              replyKeyboard,
			IsPersistent:          m.persistent,
			ResizeKeyboard:        true,
			OneTimeKeyboard:       m.oneTimeKeyboard,
			InputFieldPlaceholder: "testtesttesttest",
			Selective:             false,
		},
	}
}
