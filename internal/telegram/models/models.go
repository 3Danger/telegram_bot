package models

import (
	"strings"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type Data struct {
	UserID      int64
	ChatID      int64
	Message     string
	CallbackMap map[string]string
}

const (
	expSeparator  = "|"
	pairSeparator = "="
)

type Pair map[string]string

func (p Pair) With(k, v string) Pair {
	cp := p.Clone()
	if cp == nil {
		cp = make(Pair)
	}

	cp[k] = v
	return cp
}

func (p Pair) Clone() Pair {
	cp := make(Pair, len(p))
	for k, v := range p {
		cp[k] = v
	}

	return cp
}

func NewCallback(data Pair) string {
	if len(data) == 0 {
		return ""
	}

	rows := make([]string, 0, len(data))
	for k, v := range data {
		if v != "" {
			rows = append(rows, k+pairSeparator+v)
			continue
		}
		rows = append(rows, k)
	}

	return strings.Join(rows, expSeparator)
}

func NewMessage(update tele.Update) Data {
	msg := Data{}

	if um := update.Message; um != nil {
		msg.ChatID = um.Chat.Id
		if um.From != nil {
			msg.UserID = um.From.Id
		}
		msg.Message = um.Text
	}

	cbData := ""
	if cb := update.CallbackQuery; cb != nil {
		msg.ChatID = cb.Message.GetChat().Id
		msg.UserID = cb.From.Id
		cbData = cb.Data
	}

	msg.CallbackMap = extractCallback(cbData)

	return msg
}

func extractCallback(data string) Pair {
	rows := strings.Split(data, expSeparator)

	var (
		callbackMap = make(map[string]string, len(rows))
	)

	for _, row := range rows {
		item := strings.Split(row, pairSeparator)
		if len(item) == 0 || item[0] == "" {
			continue
		}

		if len(item) == 1 {
			callbackMap[item[0]] = ""
			continue
		}

		callbackMap[item[0]] = item[1]
	}

	return callbackMap
}
