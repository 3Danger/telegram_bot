package build

import (
	"time"

	csrepo "github.com/3Danger/telegram_bot/internal/repo/chain"
	csmem "github.com/3Danger/telegram_bot/internal/repo/chain/inmemory"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	userwrap "github.com/3Danger/telegram_bot/internal/repo/user/wrappers"
)

const timeout = time.Second * 30

func (b *Build) RepoUser() userpg.Querier {
	var repo userpg.Querier

	repo = userpg.New(b.db)

	repo = userwrap.NewQuerierWithTimeout(repo,
		userwrap.QuerierWithTimeoutConfig{
			DeleteTimeout: timeout,
			GetTimeout:    timeout,
			UpsertTimeout: timeout,
		})

	repo = userwrap.WithSkipNoRows(repo)

	return repo
}

func (b *Build) RepoChainStates() csrepo.Repo {
	repo := csmem.NewRepo(b.cnf.Repo.InMemory.MaxItems)

	return repo
}
