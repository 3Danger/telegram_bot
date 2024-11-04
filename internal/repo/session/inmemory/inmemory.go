package inmemory

import (
	"context"

	"github.com/golang/groupcache/lru"
	"github.com/samber/lo"

	r "github.com/3Danger/telegram_bot/internal/repo"
)

type repo[T any] struct {
	data *lru.Cache
}

func NewRepo[T any](maxItems int) r.Repo[T] {
	return &repo[T]{
		data: lru.New(maxItems),
	}
}

func (r *repo[T]) Get(_ context.Context, userID int64) (*T, error) {
	s, ok := r.data.Get(userID)
	if !ok {
		return nil, nil
	}

	return lo.ToPtr(s.(T)), nil
}

func (r *repo[T]) Set(_ context.Context, userID int64, state T) error {
	r.data.Remove(userID)
	r.data.Add(userID, state)

	return nil
}

func (r *repo[T]) Delete(_ context.Context, userID int64) error {
	r.data.Remove(userID)

	return nil
}
