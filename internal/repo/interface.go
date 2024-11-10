package repo

import "context"

type Repo[T any] interface {
	Get(ctx context.Context, userID int64) (T, error)
	Set(ctx context.Context, userID int64, data T) error
	Delete(ctx context.Context, userID int64) error
}
