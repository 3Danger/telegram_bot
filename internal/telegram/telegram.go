package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"
	"github.com/samber/lo"

	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/telegram/middlewares"
)

type ServiceMessageProcessor interface {
	MessageProcessor(
		ctx context.Context, msg models.Request,
	) ([]models.Response, error)
}

type Telegram struct {
	bot *tele.Bot
	smp ServiceMessageProcessor

	cnf config.Telegram
}

func New(
	cnf config.Telegram,
	smp ServiceMessageProcessor,
) (*Telegram, error) {
	bot, err := configureBot(cnf)
	if err != nil {
		return nil, err
	}

	svc := &Telegram{
		bot: bot,
		cnf: cnf,
		smp: smp,
	}

	return svc, nil
}

const requestTimeout = time.Second * 30

func configureBot(cnf config.Telegram) (*tele.Bot, error) {
	opts := &tele.BotOpts{
		BotClient:         nil,
		DisableTokenCheck: false,
		RequestOpts: &tele.RequestOpts{
			Timeout: requestTimeout,
			APIURL:  "",
		},
	}

	bot, err := tele.NewBot(cnf.Token, opts)
	if err != nil {
		return nil, fmt.Errorf("create new telegram bot: %w", err)
	}

	return bot, nil
}

func (t *Telegram) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	t.bot.BotClient = middlewares.New(t.bot.BotClient)

	opts := &tele.GetUpdatesOpts{
		Offset:         0,
		Limit:          0,
		Timeout:        60, //nolint:mnd
		AllowedUpdates: nil,
		RequestOpts:    nil,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err() //nolint:wrapcheck
		case <-ticker.C:
			updates, err := t.bot.GetUpdatesWithContext(ctx, opts)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}

				zerolog.Ctx(ctx).Err(err).Msg("getting updates")

				continue
			}

			for _, update := range updates {
				opts.Offset = max(update.UpdateId+1, opts.Offset)

				fmt.Println(string(lo.Must(json.MarshalIndent(update, "", "\t")))) //nolint:forbidigo

				msg := models.NewRequest(update)

				resp, err := t.smp.MessageProcessor(ctx, msg)
				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("processing update")
				}

				for _, r := range resp {
					if _, err := t.bot.SendMessage(r.ChatID, r.Text, r.Menu.Menu()); err != nil {
						return fmt.Errorf("sending message: %w", err)
					}
				}
			}
		}
	}
}
