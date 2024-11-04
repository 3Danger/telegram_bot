package state_consider

import (
	"context"
	"fmt"
	"strings"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type repo interface {
	Get(ctx context.Context, userID int64) (string, error)
	Set(ctx context.Context, userID int64, state string) error
}

// StateConsider - Отвечает за:
// 1. Переадресацию переадресацию запросов исходя из последней команды
// Переадресация будет происходить ИЗ общих url, ожидаемые url - /photo /media /contact etc...
// Переадресация будет происходить НА url которые мы добавим через Handle()
type StateConsider struct {
	repo            repo
	group           *tele.Group
	stateHandler    map[string]tele.HandlerFunc
	commonEndpoints map[string]struct{} // Общий команды из которых следует переадресовать /photo /media /contact etc...
}

func New(
	group *tele.Group,
	repo repo,
	commonEndpoints ...string,
) *StateConsider {
	s := &StateConsider{
		group:           group,
		repo:            repo,
		stateHandler:    map[string]tele.HandlerFunc{},
		commonEndpoints: map[string]struct{}{},
	}

	for _, e := range commonEndpoints {
		s.group.Handle(e, s.handleByState)
		s.commonEndpoints[e] = struct{}{}
	}

	return s
}

func (s *StateConsider) handleByState(c tele.Context) error {
	state, err := s.repo.Get(utils.GetContext(c), c.Sender().ID)
	if err != nil {
		return fmt.Errorf("getting last state: %w", err)
	}

	hander, ok := s.stateHandler[state]
	if !ok {
		return c.Send("вы просите странного", buttons.Reply(buttons.Home))
	}

	return hander(c)
}

func (s *StateConsider) Handle(endpoint interface{}, h tele.HandlerFunc) {
	s.stateHandler[extractEndpoint(endpoint)] = h
	s.group.Handle(endpoint, h)
}

func extractEndpoint(endpoint any) string {
	switch endpoint := endpoint.(type) {
	case tele.CallbackEndpoint:
		return endpoint.CallbackUnique()
	case string:
		return endpoint
	}

	panic(fmt.Errorf("telebot: unsupported endpoint"))
}

func (s *StateConsider) SaveLastStateMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			data := strings.Split(c.Data(), "|")
			if len(data) > 1 {
				fmt.Println()
			}

			if err := next(c); err != nil {
				return err
			}

			msg := c.Message()
			if msg == nil {
				return nil
			}

			if _, ok := s.commonEndpoints[msg.Text]; ok {
				return nil
			}

			if err := s.repo.Set(utils.GetContext(c), c.Sender().ID, msg.Text); err != nil {
				return fmt.Errorf("pushing chain state: %w", err)
			}

			return nil
		}
	}
}
