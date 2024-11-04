package build

import (
	"github.com/3Danger/telegram_bot/internal/repo"
	csrepo "github.com/3Danger/telegram_bot/internal/repo/chain-states"
	csmem "github.com/3Danger/telegram_bot/internal/repo/chain-states/inmemory"
	statemem "github.com/3Danger/telegram_bot/internal/repo/state/inmemory"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	userpgx "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
)

func (b *Build) RepoUserPermanent() repo.Repo[user.User] {
	repo := userpgx.NewRepo(b.db)

	return repo
}

func (b *Build) RepoUserSession() repo.Repo[user.User] {
	repo := statemem.NewRepo[user.User](b.cnf.Repo.InMemory.MaxItems)

	return repo
}

func (b *Build) RepoState() repo.Repo[string] {
	repo := statemem.NewRepo[string](b.cnf.Repo.InMemory.MaxItems)

	return repo
}

func (b *Build) RepoChainStates() csrepo.Repo {
	repo := csmem.NewRepo(b.cnf.Repo.InMemory.MaxItems)

	return repo
}
