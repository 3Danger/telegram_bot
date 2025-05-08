//nolint:exhaustruct
package auth

import (
	"context"
	"fmt"

	"github.com/3Danger/telegram_bot/internal/models"
	r "github.com/3Danger/telegram_bot/internal/repo"
	session "github.com/3Danger/telegram_bot/internal/repo/session/inmemory"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/services/telegram_bot/auth/handlers"
)

type repo struct {
	state r.Repo[string]
	user  user.Repo
}

type Auth struct {
	repo        repo
	handlersMap map[string]handlers.Handler
}

func NewAuth(users user.Repo) *Auth {
	handlersMap := make(map[string]handlers.Handler)
	{
		hChange := handlers.ChangeData{
			Next: nil,
		}
		hContact := handlers.ContactFromTg{
			Next: hChange,
		}

		handlersMap[""] = hContact
		handlersMap[hContact.Name()] = hContact
		handlersMap[hChange.Name()] = hChange
	}

	const cacheSize = 10000

	return &Auth{
		repo: repo{
			state: session.NewRepo[string](cacheSize),
			user:  users,
		},
		handlersMap: handlersMap,
	}
}

func (a *Auth) Handle(ctx context.Context, data models.Request) (models.Responses, error) {
	userData, err := a.createDraftIfNotExists(ctx, data.UserID())
	if err != nil {
		return nil, fmt.Errorf("creating draft user: %w", err)
	}

	resp, err := a.handle(ctx, userData, data)
	if err != nil {
		return nil, fmt.Errorf("handling draft: %w", err)
	}

	if err = a.repo.user.UpsertDraft(ctx, *userData); err != nil {
		return nil, fmt.Errorf("upserting draft user: %w", err)
	}

	return resp, nil
}

func (a *Auth) handle(ctx context.Context, user *models.User, data models.Request) (models.Responses, error) {
	state, err := a.repo.state.Get(ctx, data.UserID())
	if err != nil {
		return nil, fmt.Errorf("getting state from temporary repo: %w", err)
	}
	_ = state
	_ = user
	// TODO ....

	return nil, nil
}

func (a *Auth) createDraftIfNotExists(ctx context.Context, userID int) (*models.User, error) {
	userCompleted, err := a.repo.user.GetCompleted(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user from temporary repo: %w", err)
	}

	userDraft, err := a.repo.user.GetDraft(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user from temporary repo: %w", err)
	}

	if userCompleted != nil && userDraft == nil {
		userDraft = userCompleted
	}

	if userDraft == nil {
		userDraft = &models.User{
			ID:         userID,
			UserType:   models.UserTypeUnknown,
			FirstName:  "",
			LastName:   "",
			Phone:      "",
			Additional: "",
		}
	}

	if err = a.repo.user.UpsertDraft(ctx, *userDraft); err != nil {
		return nil, fmt.Errorf("upserting draft user: %w", err)
	}

	return userDraft, nil
}
