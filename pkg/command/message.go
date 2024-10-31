package command

import "io"

type Message struct {
	ChatID  int
	UserID  int
	Photo   func() ([]io.ReadCloser, error)
	Video   func() ([]io.ReadCloser, error)
	Text    string
	Command Name
}

func NewMessage(userID, chatID int) *Message {
	return &Message{
		ChatID: chatID,
		UserID: userID,
		Photo:  nil,
		Video:  nil,
		Text:   "",
	}
}

func (m *Message) WithText(text string) *Message {
	m.Text = text

	return m
}
