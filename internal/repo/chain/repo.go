package chain

import "context"

type Repo interface {
	Push(ctx context.Context, userID int64, state string) error
	Pop(ctx context.Context, userID int64) (string, error)
	LastState(ctx context.Context, userID int64) (string, error)
	Clear(ctx context.Context, userID int64) error
}
