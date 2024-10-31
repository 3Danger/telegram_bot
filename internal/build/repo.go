package build

import (
	"github.com/3Danger/telegram_bot/internal/repo/state"
	statemem "github.com/3Danger/telegram_bot/internal/repo/state/inmemory"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	userpgx "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
)

func (b *Build) RepoUser() user.Repo {
	repo := userpgx.NewRepo(b.db)

	return repo
}

func (b *Build) RepoState() state.Repo {
	repo := statemem.NewRepo(b.cnf.Repo.InMemory.MaxItems)

	return repo
}
