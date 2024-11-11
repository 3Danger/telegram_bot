package sender

import (
	"fmt"

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
	if _, err := s.bot.SendMessage(chatID, msg, menu.Menu()); err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}
