package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/mapper"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

type repo struct {
	q *query.Queries
}

func NewRepo(db query.DBTX) user.Repo {
	q := query.New(db)

	return &repo{
		q: q,
	}
}

func (r *repo) User(ctx context.Context, userID int) (*user.User, error) {
	row, err := r.q.User(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("making query: %w", err)
	}

	return mapper.UserToRepo(row), nil
}

func (r *repo) UpsertUser(ctx context.Context, user user.User) error {
	if err := r.q.UpsertUser(ctx, mapper.UserToQuery(user)); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}
