package telegram

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/repo/state"
	"github.com/3Danger/telegram_bot/internal/repo/user"
)

type Telegram struct {
	bot *tele.Bot
	cnf config.Telegram

	repo repo
}

type repo struct {
	user  user.Repo
	state state.Repo
}

func New(
	ctx context.Context,
	cnf config.Telegram,
	userRepo user.Repo,
	stateRepo state.Repo,
) (*Telegram, error) {
	bot, err := configureBot(cnf)
	if err != nil {
		return nil, err
	}

	svc := &Telegram{
		bot: bot,
		cnf: cnf,
		repo: repo{
			user:  userRepo,
			state: stateRepo,
		},
	}

	if err = svc.configureRoute(ctx); err != nil {
		return nil, fmt.Errorf("configure telegram routes: %w", err)
	}

	return svc, nil
}

func configureBot(cnf config.Telegram) (*tele.Bot, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  cnf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		//Verbose: cnf.Debug,
		Verbose: false,
		OnError: onError,
	})
	if err != nil {
		return nil, fmt.Errorf("create new telegram bot: %w", err)
	}

	return bot, nil
}

func onError(err error, c tele.Context) {
	l, ok := c.Get("zerolog").(zerolog.Logger)
	if !ok {
		l = zerolog.New(os.Stdout)
		l.Error().Msg("DEFAULT LOGGER")
	}

	l.Err(err).Interface("message", c.Message()).Send()
}

// Обновляем configureRoute для добавления нового обработчика
func (t *Telegram) configureRoute(ctx context.Context) error {

	t.bot.Use(
		middlewareContext(ctx),
		middlewareFilterBot(),
		t.middlewareSaveLastState(),
	)

	start := t.bot.Group()
	{
		start.Handle("/start", t.handlerHome)
		start.Handle(home, t.handlerHome)
		start.Handle(auth, t.handlerAuth)
		start.Handle(tele.OnText, t.handlerAuth)
		start.Handle(tele.OnContact, t.handlerAuth)
	}

	return nil
}

func (t *Telegram) Start(ctx context.Context) error {
	t.bot.Start()

	<-ctx.Done()

	return nil
}
