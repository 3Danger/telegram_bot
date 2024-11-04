package user

import (
	"errors"
	"time"

	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

var ErrUserNotFound = errors.New("user not found")

type Type string

func (t Type) String() string { return string(t) }

func (t Type) Valid() bool {
	switch t {
	case TypeCustomer:
		return true
	case TypeSupplier:
		return true
	}

	return false
}

const (
	TypeUndefined = Type(query.UserTypeUndefined)
	TypeCustomer  = Type(query.UserTypeCustomer)
	TypeSupplier  = Type(query.UserTypeSupplier)
)

type User struct {
	ID        int64
	Type      Type
	FirstName string
	LastName  string
	Surname   string
	Phone     string
	Whatsapp  string
	Telegram  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
