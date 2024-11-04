package telegram

import (
	"context"
	"errors"
	"fmt"
	"time"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/repo"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/services/auth"
	"github.com/3Danger/telegram_bot/internal/telegram/middlewares"
	"github.com/3Danger/telegram_bot/internal/telegram/validator"
)

type Telegram struct {
	bot *tele.Bot

	cnf       config.Telegram
	validator *validator.MediaValidator
	svc       Services
	repo      Repo
}

type Services struct {
	auth *auth.Service
}

type Repo struct {
	user    repo.Repo[user.User]
	state   repo.Repo[string]
	command repo.Repo[string]
}

func New(
	cnf config.Telegram,
	userRepo repo.Repo[user.User],
	stateRepo repo.Repo[string],
	commandRepo repo.Repo[string],
	auth *auth.Service,
) (*Telegram, error) {
	bot, err := configureBot(cnf)
	if err != nil {
		return nil, err
	}

	svc := &Telegram{
		bot: bot,
		cnf: cnf,
		repo: Repo{
			user:    userRepo,
			state:   stateRepo,
			command: commandRepo,
		},
		validator: configureMediaValidator(),
		svc:       Services{auth: auth},
	}

	return svc, nil
}

func configureMediaValidator() *validator.MediaValidator {
	const megaByte = 1024 * 1024

	return validator.New(
		validator.NewBound[int64](1024, 8192),
		validator.NewBound[int64](megaByte/10, megaByte*10),
		validator.NewBound[int64](480, 2592),
		validator.NewBound[int64](0, megaByte*15),
		validator.NewBound[time.Duration](time.Second*15, time.Second*60),
	)
}

func configureBot(cnf config.Telegram) (*tele.Bot, error) {
	opts := &tele.BotOpts{
		BotClient:         nil,
		DisableTokenCheck: false,
		RequestOpts: &tele.RequestOpts{
			Timeout: time.Second * 30,
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
		Timeout:        60,
		AllowedUpdates: nil,
		RequestOpts:    nil,
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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

				if auth.IsValidationErr(err) {
					if _, err := t.bot.SendMessage(getChatID(update), err.Error(), nil); err != nil {
						zerolog.Ctx(ctx).Err(err).Msg("sending message")
					}

					continue
				}

				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("processing update")
				}
			}
		}
	}

}

func getChatID(update tele.Update) int64 {
	if m := update.Message; m != nil {
		return m.Chat.Id
	}
	if c := update.CallbackQuery; c != nil {
		return c.Message.GetMessageId()
	}

	return 0
}
