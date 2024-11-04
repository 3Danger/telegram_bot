package models

import (
	"strings"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type Data struct {
	UserID        int64
	ChatID        int64
	Message       string
	CallbackOrder []string
	CallbackMap   map[string]string
}

const (
	expSeparator  = "|"
	pairSeparator = "="
)

type PairKeyValues struct {
	Key, Value string
}

func Pair(key, value string) PairKeyValues {
	return PairKeyValues{
		Key:   key,
		Value: value,
	}
}

func NewCallback(pairs ...PairKeyValues) string {
	rows := make([]string, 0, len(pairs))
	for _, item := range pairs {
		if item.Value != "" {
			rows = append(rows, item.Key+pairSeparator+item.Value)
			continue
		}
		rows = append(rows, item.Key)
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
	msg.CallbackOrder, msg.CallbackMap = extractCallback(cbData)

	return msg
}

func extractCallback(data string) ([]string, map[string]string) {
	rows := strings.Split(data, expSeparator)

	var (
		callbackOrder = make([]string, 0, len(rows))
		callbackMap   = make(map[string]string, len(rows))
	)

	for _, row := range rows {
		item := strings.Split(row, pairSeparator)
		if len(item) == 0 || item[0] == "" {
			continue
		}

		callbackOrder = append(callbackOrder, item[0])
		if len(item) == 1 {
			callbackMap[item[0]] = ""
			continue
		}

		callbackMap[item[0]] = item[1]
	}

	return callbackOrder, callbackMap
}
