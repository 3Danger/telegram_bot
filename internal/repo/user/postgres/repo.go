package postgres

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/models"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/mapper"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

//go:generate ifacemaker -f repo.go -s repo -i Repo -p user -o ../repo.go -D -y "Репозиторий пользователей"
//go:generate gowrap gen -p .. -i Repo -o ../wrappers/timeout.go -t timeout
//go:generate gowrap gen -p .. -i Repo -o ../wrappers/skip_no_rows.go -t ../../../../gowrap/skip_no_rows.tmpl
type repo struct {
	q query.Querier
}

func NewRepo(dbtx query.DBTX) user.Repo {
	return &repo{
		q: query.New(dbtx),
	}
}

func (r *repo) ApproveChanges(ctx context.Context, id int) error {
	if err := r.q.ApproveChanges(ctx, int64(id)); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) DeleteDraft(ctx context.Context, id int) error {
	if err := r.q.Delete(ctx, &query.DeleteParams{
		ID:            int64(id),
		RecordingMode: query.RecordingModeDraft,
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) GetCompleted(ctx context.Context, id int) (*models.User, error) {
	row, err := r.q.Get(ctx, &query.GetParams{
		ID:            int64(id),
		RecordingMode: query.RecordingModeCompleted,
	})
	if err != nil {
		return nil, fmt.Errorf("making query: %w", err)
	}

	return mapper.UserRepoToModel(row), nil
}

func (r *repo) GetDraft(ctx context.Context, id int) (*models.User, error) {
	row, err := r.q.Get(ctx, &query.GetParams{
		ID:            int64(id),
		RecordingMode: query.RecordingModeDraft,
	})
	if err != nil {
		return nil, fmt.Errorf("making query: %w", err)
	}

	return mapper.UserRepoToModel(row), nil
}

func (r *repo) SetAdditional(ctx context.Context, id int, additionalParams string) error {
	if err := r.q.SetAdditional(ctx, &query.SetAdditionalParams{
		Additional: additionalParams,
		ID:         int64(id),
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) SetFirstName(ctx context.Context, id int, firstNameParams string) error {
	if err := r.q.SetFirstName(ctx, &query.SetFirstNameParams{
		FirstName: firstNameParams,
		ID:        int64(id),
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) SetLastName(ctx context.Context, id int, lastNameParams string) error {
	if err := r.q.SetLastName(ctx, &query.SetLastNameParams{
		LastName: lastNameParams,
		ID:       int64(id),
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) SetPhone(ctx context.Context, id int, phoneParams string) error {
	if err := r.q.SetPhone(ctx, &query.SetPhoneParams{
		Phone: phoneParams,
		ID:    int64(id),
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) SetUserType(ctx context.Context, id int, userTypeParams string) error {
	if err := r.q.SetUserType(ctx, &query.SetUserTypeParams{
		UserType: query.UserType(userTypeParams),
		ID:       int64(id),
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}

func (r *repo) UpsertDraft(ctx context.Context, user models.User) error {
	if err := r.q.Upsert(ctx, &query.UpsertParams{
		ID:            int64(user.ID),
		RecordingMode: query.RecordingModeDraft,
		UserType:      query.UserType(user.UserType),
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Additional:    user.Additional,
	}); err != nil {
		return fmt.Errorf("making query: %w", err)
	}

	return nil
}
