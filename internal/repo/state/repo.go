package state

import "context"

type Repo interface {
	Get(ctx context.Context, userID int64) (string, error)
	Set(ctx context.Context, userID int64, state string) error
}
