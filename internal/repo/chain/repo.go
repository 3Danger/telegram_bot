package chain

import "context"

type Repo interface {
	Push(ctx context.Context, userID int, state string) error
	Pop(ctx context.Context, userID int) (string, error)
	LastState(ctx context.Context, userID int) (string, error)
	Clear(ctx context.Context, userID int) error
}
