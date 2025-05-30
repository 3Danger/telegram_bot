//nolint:forcetypeassert
package inmemory

import (
	"container/list"
	"context"

	"github.com/golang/groupcache/lru"

	cs "github.com/3Danger/telegram_bot/internal/repo/chain"
)

type repo struct {
	data *lru.Cache
}

func NewRepo(maxItems int) cs.Repo {
	return &repo{
		data: lru.New(maxItems),
	}
}

func (r *repo) Push(_ context.Context, userID int, state string) error {
	var l *list.List

	curListAny, ok := r.data.Get(userID)
	if ok {
		l = curListAny.(*list.List)
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

func (r *repo) Pop(_ context.Context, userID int) (string, error) {
	l, el := r.lastState(userID)
	if el == nil || l == nil {
		return "", nil
	}

	l.Remove(el)

	if l.Len() == 0 {
		r.data.Remove(userID)
	}

	return el.Value.(string), nil
}

func (r *repo) LastState(_ context.Context, userID int) (string, error) {
	l, el := r.lastState(userID)
	if el == nil || l == nil {
		return "", nil
	}

	return el.Value.(string), nil
}

func (r *repo) lastState(userID int) (*list.List, *list.Element) {
	curListAny, ok := r.data.Get(userID)
	if !ok {
		return nil, nil
	}

	l := curListAny.(*list.List)

	el := l.Front()
	if el == nil {
		return l, nil
	}

	return l, el
}

func (r *repo) Clear(_ context.Context, userID int) error {
	r.data.Remove(userID)

	return nil
}
