package telegrambot

import (
	"context"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/repo/chain"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	"github.com/3Danger/telegram_bot/internal/services/telegram_bot/validator"
)

type Handler interface {
	Handle(ctx context.Context, data models.Request) (models.Responses, error)
}

type Telegram struct {
	router map[string]Handler

	validator *validator.MediaValidator
	repo      Repo
}

type Repo struct {
	user  userpg.Querier
	chain chain.Repo
}

func New(
	userRepo userpg.Querier,
	repoChainStates chain.Repo,
) *Telegram {
	svc := &Telegram{
		router: make(map[string]Handler),
		repo: Repo{
			user:  userRepo,
			chain: repoChainStates,
		},
		validator: validator.Default(),
	}

	svc.configureRoutes()

	return svc
}
