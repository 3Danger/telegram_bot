package sender

import (
	tele "github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/3Danger/telegram_bot/internal/telegram/keyboard"
)

type Sender struct {
	bot *tele.Bot
}

func New(bot *tele.Bot) *Sender {
	return &Sender{
		bot: bot,
	}
}

func (s *Sender) Send(chatID int64, msg string, menu keyboard.Menu) error {
	_, err := s.bot.SendMessage(chatID, msg, menu.Menu())

	return err
}
