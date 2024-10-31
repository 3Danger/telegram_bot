package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/3Danger/telegram_bot/internal/config"
	"github.com/3Danger/telegram_bot/internal/repo/state"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	c "github.com/3Danger/telegram_bot/pkg/command"
)

type Telegram struct {
	api *tgbotapi.BotAPI
	cnf config.Telegram

	commands map[c.Name]*c.Command
	repo     repo
}

type repo struct {
	user  user.Repo
	state state.Repo
}

func New(cnf config.Telegram, userRepo user.Repo, stateRepo state.Repo) (*Telegram, error) {
	api, err := tgbotapi.NewBotAPI(cnf.Token)
	if err != nil {
		return nil, fmt.Errorf("creating new telegram api: %w", err)
	}

	api.Debug = cnf.Debug

	return &Telegram{
		api:      api,
		cnf:      cnf,
		commands: make(map[c.Name]*c.Command),
		repo: repo{
			user:  userRepo,
			state: stateRepo,
		},
	}, nil
}

func (t *Telegram) addCommand(
	name c.Name, processor c.ProcessorFn, middleware ...c.MiddlewareFn,
) {
	m := make([]c.MiddlewareFn, len(middleware)+1)
	m = append(m, t.saveCommand(name))
	m = append(m, middleware...)

	t.commands[name] = c.New(processor, m...)
}

func (t *Telegram) InitCommands() {
	t.addCommand(commandHome, t.home, t.middlewareAuthCheck)
	t.addCommand(commandHome, t.home, t.middlewareAuthCheck)
}
