package state

import "context"

type Repo interface {
	State(ctx context.Context, userID int) (string, error)
	SetState(ctx context.Context, userID int, state string) error
}
