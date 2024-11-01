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

func (r *repo) Get(_ context.Context, userID int64) (string, error) {
	s, ok := r.data.Get(userID)
	if !ok {
		return "", nil
	}

	return s.(string), nil
}

func (r *repo) Set(_ context.Context, userID int64, state string) error {
	r.data.Remove(userID)
	r.data.Add(userID, state)

	return nil
}

func (r *repo) Delete(_ context.Context, userID int64) error {
	r.data.Remove(userID)

	return nil
}
