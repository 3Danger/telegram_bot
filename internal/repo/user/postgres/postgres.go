package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	r "github.com/3Danger/telegram_bot/internal/repo"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/mapper"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

type repo struct {
	q *query.Queries
}

func NewRepo(db query.DBTX) r.Repo[user.User] {
	q := query.New(db)

	return &repo{
		q: q,
	}
}

func (r *repo) Get(ctx context.Context, userID int64) (*user.User, error) {
	row, err := r.q.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("making query: %w", err)
	}

	return mapper.UserToRepo(row), nil
}

func (r *repo) Set(ctx context.Context, userID int64, user user.User) error {
	if err := r.q.Set(ctx, mapper.UserToQuery(userID, user)); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, userID int64) error {
	if err := r.q.Delete(ctx, userID); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}
