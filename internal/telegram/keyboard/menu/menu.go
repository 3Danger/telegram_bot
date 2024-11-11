package menu

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard"
)

type Inline struct {
	buttons [][]keyboard.InlineButton
}

func NewInline(row ...keyboard.InlineButton) *Inline {
	return &Inline{buttons: [][]keyboard.InlineButton{row}}
}

func (m *Inline) Add(row ...keyboard.InlineButton) *Inline {
	m.buttons = append(m.buttons, row)

	return m
}

func (m *Inline) Menu() *tele.SendMessageOpts {
	inlineKeyboard := make([][]tele.InlineKeyboardButton, len(m.buttons))
	for i, row := range m.buttons {
		inlineKeyboard[i] = make([]tele.InlineKeyboardButton, len(row))
		for j, button := range row {
			inlineKeyboard[i][j] = button.Button()
		}
	}

	return &tele.SendMessageOpts{
		BusinessConnectionId: "",
		MessageThreadId:      0,
		ParseMode:            "",
		Entities:             nil,
		LinkPreviewOptions:   nil,
		DisableNotification:  false,
		ProtectContent:       false,
		AllowPaidBroadcast:   false,
		MessageEffectId:      "",
		ReplyParameters:      nil,
		ReplyMarkup:          tele.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard},
		RequestOpts:          nil,
	}
}

type Reply struct {
	buttons         [][]keyboard.ReplyButton
	oneTimeKeyboard bool
	persistent      bool
}

func NewReply(row ...keyboard.ReplyButton) *Reply {
	return &Reply{
		buttons:         [][]keyboard.ReplyButton{row},
		oneTimeKeyboard: false,
		persistent:      false,
	}
}

func (m *Reply) Add(row ...keyboard.ReplyButton) *Reply {
	m.buttons = append(m.buttons, row)

	return m
}

func (m *Reply) OneTime(is bool) *Reply {
	m.oneTimeKeyboard = is

	return m
}

func (m *Reply) Persistent(is bool) *Reply {
	m.persistent = is

	return m
}

func (m *Reply) Menu() *tele.SendMessageOpts {
	replyKeyboard := make([][]tele.KeyboardButton, len(m.buttons))
	for i, row := range m.buttons {
		replyKeyboard[i] = make([]tele.KeyboardButton, len(row))
		for j, button := range row {
			replyKeyboard[i][j] = button.Button()
		}
	}

	return &tele.SendMessageOpts{
		BusinessConnectionId: "",
		MessageThreadId:      0,
		ParseMode:            "",
		Entities:             nil,
		LinkPreviewOptions:   nil,
		DisableNotification:  false,
		ProtectContent:       false,
		AllowPaidBroadcast:   false,
		MessageEffectId:      "",
		ReplyParameters:      nil,
		ReplyMarkup: tele.ReplyKeyboardMarkup{
			Keyboard:              replyKeyboard,
			IsPersistent:          m.persistent,
			ResizeKeyboard:        true,
			OneTimeKeyboard:       m.oneTimeKeyboard,
			InputFieldPlaceholder: "testtesttesttest",
			Selective:             false,
		},
		RequestOpts: nil,
	}
}
