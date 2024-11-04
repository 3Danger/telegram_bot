package middlewares

import (
	"context"
	"encoding/json"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type middleware struct {
	botClient tele.BotClient
}

func New(botClient tele.BotClient) tele.BotClient {
	return &middleware{
		botClient: botClient,
	}
}

func (m *middleware) RequestWithContext(
	ctx context.Context,
	token string,
	method string,
	params map[string]string,
	data map[string]tele.FileReader,
	opts *tele.RequestOpts,
) (json.RawMessage, error) {
	return m.botClient.RequestWithContext(ctx, token, method, params, data, opts)
}

func (m *middleware) GetAPIURL(opts *tele.RequestOpts) string {
	return m.botClient.GetAPIURL(opts)
}

func (m *middleware) FileURL(token string, tgFilePath string, opts *tele.RequestOpts) string {
	return m.botClient.FileURL(token, tgFilePath, opts)
}
