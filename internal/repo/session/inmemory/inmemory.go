//nolint:forcetypeassert
package inmemory

import (
	"context"

	"github.com/golang/groupcache/lru"

	r "github.com/3Danger/telegram_bot/internal/repo"
)

var _ r.Repo[struct{}] = &repo[struct{}]{} //nolint:exhaustruct

type repo[T any] struct {
	data *lru.Cache
}

func NewRepo[T any](maxItems int) r.Repo[T] {
	return &repo[T]{
		data: lru.New(maxItems),
	}
}

func (r *repo[T]) Get(_ context.Context, userID int) (T, error) {
	s, ok := r.data.Get(userID)
	if !ok {
		var t T

		return t, nil
	}

	return s.(T), nil
}

func (r *repo[T]) Set(_ context.Context, userID int, state T) error {
	r.data.Remove(userID)
	r.data.Add(userID, state)

	return nil
}

func (r *repo[T]) Delete(_ context.Context, userID int) error {
	r.data.Remove(userID)

	return nil
}
