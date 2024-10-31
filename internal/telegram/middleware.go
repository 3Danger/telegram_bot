package telegram

import (
	"context"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
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

			zerolog.Ctx(getContext(c)).Info().Interface("MSG_BODY", c.Message()).Send()

			if c.Sender().IsBot {
				return c.Send("Ботов пока не обслуживаем")
			}

			err := next(c)

			return err
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
