package inmemory

import (
	"context"

	"github.com/golang/groupcache/lru"

	"github.com/3Danger/telegram_bot/internal/repo/state"
)

type repo struct {
	data *lru.Cache
}

func NewRepo(maxItems int) state.Repo {
	return &repo{
		data: lru.New(maxItems),
	}
}

func (r *repo) State(_ context.Context, userID int) (string, error) {
	s, ok := r.data.Get(userID)
	if !ok {
		return "", nil
	}

	return s.(string), nil
}

func (r *repo) SetState(_ context.Context, userID int, command string) error {
	r.data.Remove(userID)
	r.data.Add(userID, command)

	return nil
}
