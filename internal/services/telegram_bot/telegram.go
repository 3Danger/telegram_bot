package telegrambot

import (
	"context"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/repo/chain"
	"github.com/3Danger/telegram_bot/internal/repo/user"
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
	user  user.Repo
	chain chain.Repo
}

func New(
	userRepo user.Repo,
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
