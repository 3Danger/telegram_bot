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

func (r *repo) User(ctx context.Context, userID int64) (*user.User, error) {
	row, err := r.q.User(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("making query: %w", err)
	}

	return mapper.UserToRepo(row), nil
}

func (r *repo) CreateUser(ctx context.Context, user user.User) error {
	if err := r.q.CreateUser(ctx, mapper.UserToQuery(user)); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) UpdateUserContactTelegram(ctx context.Context, userID int64, telegram string) error {
	affected, err := r.q.UpdateUserContactTelegram(ctx, query.UpdateUserContactTelegramParams{
		ID:       userID,
		Telegram: telegram,
	})
	if err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	if affected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *repo) UpdateUserContactWhatsapp(ctx context.Context, userID int64, whatsapp string) error {
	affected, err := r.q.UpdateUserContactWhatsapp(ctx, query.UpdateUserContactWhatsappParams{
		ID:       userID,
		Whatsapp: whatsapp,
	})
	if err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	if affected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *repo) UpdateUserContactPhone(ctx context.Context, userID int64, phone string) error {
	affected, err := r.q.UpdateUserContactPhone(ctx, query.UpdateUserContactPhoneParams{
		ID:    userID,
		Phone: phone,
	})
	if err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	if affected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}
