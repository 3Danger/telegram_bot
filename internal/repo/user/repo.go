package user

import (
	"context"
	"errors"
	"time"
)

type Repo interface {
	User(ctx context.Context, userID int64) (*User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUserContactTelegram(ctx context.Context, userID int64, telegram string) error
	UpdateUserContactWhatsapp(ctx context.Context, userID int64, whatsapp string) error
	UpdateUserContactPhone(ctx context.Context, userID int64, phone string) error
}

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID         int64
	IsSupplier bool
	FirstName  string
	LastName   string
	Surname    string
	Phone      string
	Whatsapp   string
	Telegram   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
