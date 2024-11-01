package inmemory

import (
	"container/list"
	"context"

	"github.com/golang/groupcache/lru"

	cs "github.com/3Danger/telegram_bot/internal/repo/chain-states"
)

type repo struct {
	data *lru.Cache
}

func NewRepo(maxItems int) cs.Repo {
	return &repo{
		data: lru.New(maxItems),
	}
}

func (r *repo) Push(ctx context.Context, userID int64, state string) error {
	var l *list.List

	curListAny, ok := r.data.Get(userID)
	if ok {
		l = curListAny.(*list.List)
		r.data.Remove(userID)
	} else {
		l = list.New().Init()
	}

	if el := l.Front(); el != nil && el.Value.(string) == state {
		return nil
	}

	l.PushFront(state)

	r.data.Add(userID, l)

	return nil
}

func (r *repo) Pop(ctx context.Context, userID int64) (string, error) {
	curListAny, ok := r.data.Get(userID)
	if !ok {
		return "", nil
	}

	l := curListAny.(*list.List)

	el := l.Front()
	if el == nil {
		return "", nil
	}
	l.Remove(el)

	if l.Len() == 0 {
		r.data.Remove(userID)
	}

	return el.Value.(string), nil
}

func (r *repo) Clear(ctx context.Context, userID int64) error {
	r.data.Remove(userID)

	return nil
}
