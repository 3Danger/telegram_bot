package handlers

import (
	"context"

	"github.com/3Danger/telegram_bot/internal/models"
	repocache "github.com/3Danger/telegram_bot/internal/repo"
)

type Repo repocache.Repo[*models.User]

type Handler interface {
	Process(ctx context.Context, cache Repo, data models.Request) (models.Responses, error)
	Name() string
}
