package telegram

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/3Danger/telegram_bot/internal/repo/user"
	c "github.com/3Danger/telegram_bot/pkg/command"
)

/*
Выберите действие
	1. Выложить товар
	2. Показать мои товары
	3. Мои данные
*/

type Command interface {
	GetMessage() string
	NextCommand(path string) Command
}

const (
	commandHome                c.Name = "В начало"
	commandCommit              c.Name = "сохранить"
	commandSupplierAuth        c.Name = "Авторизоваться как поставщик"
	commandSupplierListMyGoods c.Name = "Мои товары"
	commandSupplierAddGoods    c.Name = "Добавить товары"
	commandCustomerAuth        c.Name = "Авторизоваться как клиент"
	commandCustomerListGoods   c.Name = "Посмотреть товары"
)

func (t *Telegram) saveCommand(commandName c.Name) c.MiddlewareFn {
	return func(p c.ProcessorFn) c.ProcessorFn {
		return func(ctx context.Context, msg *c.Message) (*c.Message, []c.Name) {
			if err := t.repo.state.SetState(ctx, msg.UserID, commandName.String()); err != nil {
				return handleErr(ctx, msg, err)
			}

			return p(ctx, msg)
		}
	}
}

func (t *Telegram) middlewareAuthCheck(p c.ProcessorFn) c.ProcessorFn {
	return func(ctx context.Context, msg *c.Message) (*c.Message, []c.Name) {
		user, err := t.repo.user.User(ctx, msg.UserID)
		if err != nil {
			return handleErr(ctx, msg, err)
		}

		if user == nil || !user.HasRegistered {
			return handleUnauthorised(ctx, msg)
		}

		ctx = setUser(ctx, user)

		return p(ctx, msg)
	}
}

func setUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func getUser(ctx context.Context) *user.User {
	return ctx.Value("user").(*user.User)
}

func (t *Telegram) home(ctx context.Context, msg *c.Message) (*c.Message, []c.Name) {
	respMsg := c.NewMessage(msg.UserID, msg.ChatID)

	if getUser(ctx).IsSupplier {
		return respMsg, []c.Name{commandSupplierAddGoods, commandSupplierListMyGoods}
	}

	return respMsg, []c.Name{commandCustomerListGoods}
}

func (t *Telegram) authSupplier(ctx context.Context, msg *c.Message) (*c.Message, []c.Name) {
	respMsg := c.NewMessage(msg.UserID, msg.ChatID).
		WithText("Введите ФИО\nПример:\nИванов Иван Иванович")

	return respMsg, []c.Name{}
}

func handleUnauthorised(ctx context.Context, msg *c.Message) (*c.Message, []c.Name) {
	zerolog.Ctx(ctx).Warn().Msg("unauthorised")

	msg.WithText("Вы не авторизованны")

	return msg, []c.Name{commandSupplierAuth, commandCustomerAuth, commandHome}
}

func handleErr(ctx context.Context, msg *c.Message, err error) (*c.Message, []c.Name) {
	zerolog.Ctx(ctx).Err(err).Send()

	msg = msg.WithText("Ошибка, попробуйте позже")

	return msg, []c.Name{commandHome}
}
