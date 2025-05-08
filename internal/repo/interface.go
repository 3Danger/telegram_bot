package repo

import "context"

type Repo[T any] interface {
	Get(ctx context.Context, userID int) (T, error)
	Set(ctx context.Context, userID int, data T) error
	Delete(ctx context.Context, userID int) error
}
