package build

import (
	"time"

	csrepo "github.com/3Danger/telegram_bot/internal/repo/chain"
	csmem "github.com/3Danger/telegram_bot/internal/repo/chain/inmemory"
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres"
	userwrap "github.com/3Danger/telegram_bot/internal/repo/user/wrappers"
)

const timeout = time.Second * 30

func (b *Build) RepoUser() user.Repo {
	var repo user.Repo

	repo = postgres.NewRepo(b.db)

	repo = userwrap.NewRepoWithTimeout(repo,
		userwrap.RepoWithTimeoutConfig{
			ApproveChangesTimeout: timeout,
			DeleteDraftTimeout:    timeout,
			GetCompletedTimeout:   timeout,
			GetDraftTimeout:       timeout,
			SetAdditionalTimeout:  timeout,
			SetFirstNameTimeout:   timeout,
			SetLastNameTimeout:    timeout,
			SetPhoneTimeout:       timeout,
			SetUserTypeTimeout:    timeout,
			UpsertDraftTimeout:    timeout,
		})

	repo = userwrap.WithSkipNoRows(repo)

	return repo
}

func (b *Build) RepoChainStates() csrepo.Repo {
	repo := csmem.NewRepo(b.cnf.Repo.InMemory.MaxItems)

	return repo
}
