package build

import (
	"time"

	csrepo "github.com/3Danger/telegram_bot/internal/repo/chain-states"
	csmem "github.com/3Danger/telegram_bot/internal/repo/chain-states/inmemory"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	userwrap "github.com/3Danger/telegram_bot/internal/repo/user/wrappers"
)

func (b *Build) RepoUserPermanent() userpg.Querier {
	var repo userpg.Querier

	repo = userpg.New(b.db)

	repo = userwrap.NewQuerierWithTimeout(repo,
		userwrap.QuerierWithTimeoutConfig{
			DeleteTimeout: time.Second * 30,
			GetTimeout:    time.Second * 30,
			UpsertTimeout: time.Second * 30,
		})

	repo = userwrap.WithSkipNoRows(repo)

	return repo
}

func (b *Build) RepoChainStates() csrepo.Repo {
	repo := csmem.NewRepo(b.cnf.Repo.InMemory.MaxItems)

	return repo
}
