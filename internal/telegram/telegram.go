package telegram

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/config"
	cs "github.com/3Danger/telegram_bot/internal/repo/chain-states"
	"github.com/3Danger/telegram_bot/internal/repo/state"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/telegram/constants"
	"github.com/3Danger/telegram_bot/internal/telegram/media"
)

type Telegram struct {
	bot  *tele.Bot
	cnf  config.Telegram
	v    *media.Validator
	repo repo
}

type repo struct {
	user        user.Repo
	state       state.Repo
	chainStates cs.Repo
}

func New(
	ctx context.Context,
	cnf config.Telegram,
	userRepo user.Repo,
	stateRepo state.Repo,
	csRepo cs.Repo,
) (*Telegram, error) {
	bot, err := configureBot(cnf)
	if err != nil {
		return nil, err
	}

	const megaByte = 1024 * 1024

	mediaValidator := media.NewValidator(
		media.NewBound[int](1024, 8192),
		media.NewBound[int64](megaByte/10, megaByte*10),
		media.NewBound[int](480, 2592),
		media.NewBound[int64](0, megaByte*15),
		media.NewBound[time.Duration](time.Second*15, time.Second*60),
	)

	svc := &Telegram{
		bot: bot,
		cnf: cnf,
		repo: repo{
			user:        userRepo,
			state:       stateRepo,
			chainStates: csRepo,
		},
		v: mediaValidator,
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
		start.Handle(constants.Home, t.handlerHome)
	}

	supplier := t.bot.Group()
	{
		supplier.Handle(constants.SupplierShowItems, t.handlerSupplierShowItems)
		supplier.Handle(constants.SupplierPostItems, t.handlerSupplierPostItems)

		//TODO Мне не нравится что эндпоинты (tele.OnPhoto, tele.OnMedia) "всеядные",
		// нужен какой то единый конструктор который будет разделять по под-эндпоинтам
		supplier.Handle(tele.OnPhoto, t.handlerSupplierPostItems)
		supplier.Handle(tele.OnMedia, t.handlerSupplierPostItems)
	}
	customer := t.bot.Group()
	{
		customer.Handle(constants.CustomerShowItems, t.handlerCustomerHome)
	}

	auth := t.bot.Group()
	{
		// Обработка команд редактирования
		auth.Handle(constants.Auth, t.handlerAuth)
		auth.Handle(tele.OnText, t.handlerAuth)
		auth.Handle(tele.OnContact, t.handlerAuth)
		auth.Handle(constants.Back, t.handlerBackNavigation)
		auth.Handle(constants.AuthEditName, t.handlerEditName)
		auth.Handle(constants.AuthEditPhone, t.handlerEditPhone)
		auth.Handle(constants.AuthConfirm, t.handlerConfirmRegistration)
	}

	return nil
}

func (t *Telegram) Start(ctx context.Context) error {
	t.bot.Start()

	<-ctx.Done()

	return nil
}
