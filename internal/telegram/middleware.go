package telegram

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/telegram/constants"
)

func getContext(c tele.Context) context.Context {
	return c.Get("ctx").(context.Context)
}

func setContext(c tele.Context, ctx context.Context) {
	c.Set("ctx", ctx)
}

func middlewareContext(ctx context.Context) tele.MiddlewareFunc {
	l := *zerolog.Ctx(ctx)

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			l := l
			sender := c.Sender()
			message := c.Message()

			if message != nil && sender != nil {
				l = l.With().
					Int64("userID", sender.ID).
					Int64("chatID", message.Chat.ID).
					Logger()
			}

			setContext(c, l.WithContext(ctx))

			err := next(c)

			return err
		}
	}
}
func middlewareFilterBot() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().IsBot {
				return c.Send("Ботов пока не обслуживаем")
			}

			err := next(c)

			return err
		}
	}
}

func (t *Telegram) middlewareSaveLastState() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if msg := c.Message(); msg != nil && constants.IsValid(msg.Text) {
				ctx := getContext(c)
				switch msg.Text {
				case constants.Back:
					prev, err := t.repo.chainStates.Pop(ctx, c.Sender().ID)
					if err != nil {
						return fmt.Errorf("cleanning chain states: %w", err)
					}
					return c.Send("", createMenu(prev))
				case constants.Home:
					if err := t.repo.chainStates.Clear(ctx, c.Sender().ID); err != nil {
						return fmt.Errorf("cleanning chain states: %w", err)
					}
				default:
					if err := next(c); err != nil {
						return err
					}

					if err := t.repo.chainStates.Push(ctx, c.Sender().ID, msg.Text); err != nil {
						return fmt.Errorf("pushing chain state: %w", err)
					}

					return nil
				}
			}

			if err := next(c); err != nil {
				return err
			}

			return nil
		}
	}
}

//func (t *Telegram) middlewareAuthCheck(p tele.HandlerFunc) tele.HandlerFunc {
//	return func(c tele.Context) error {
//		user, err := t.repo.user.User(ctx, c.Message().Contact.UserID)
//		if err != nil {
//			return
//		}
//
//		if user == nil || !user.HasRegistered {
//			return
//		}
//
//		ctx = setUser(ctx, user)
//
//		return p(ctx, msg)
//	}
//}
