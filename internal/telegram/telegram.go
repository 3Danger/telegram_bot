package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/repo/chain"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	"github.com/3Danger/telegram_bot/internal/telegram/middlewares"
	"github.com/3Danger/telegram_bot/internal/telegram/models"
	"github.com/3Danger/telegram_bot/internal/telegram/sender"
	"github.com/3Danger/telegram_bot/internal/telegram/validator"
)

type Handler interface {
	Handle(ctx context.Context, data models.Request) error
}

type Telegram struct {
	bot    *tele.Bot
	sender *sender.Sender

	router map[string]Handler

	cnf       config.Telegram
	validator *validator.MediaValidator
	repo      Repo
}

type Repo struct {
	user  userpg.Querier
	chain chain.Repo
}

func New(
	cnf config.Telegram,
	userRepo userpg.Querier,
	repoChainStates chain.Repo,
) (*Telegram, error) {
	bot, err := configureBot(cnf)
	if err != nil {
		return nil, err
	}

	sender := sender.New(bot)

	svc := &Telegram{
		bot:    bot,
		sender: sender,
		cnf:    cnf,
		router: make(map[string]Handler),
		repo: Repo{
			user:  userRepo,
			chain: repoChainStates,
		},
		validator: validator.Default(),
	}

	svc.configureRoutes()

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

	//nolint:gomnd
	opts := &tele.GetUpdatesOpts{
		Offset:         0,
		Limit:          0,
		Timeout:        60,
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

				err = t.updateProcessor(ctx, update)
				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("processing update")
				}
			}
		}
	}
}
