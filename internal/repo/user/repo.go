package user

import (
	"context"
	"time"
)

type Repo interface {
	User(ctx context.Context, userID int) (*User, error)
}

type User struct {
	ID            int
	HasRegistered bool
	IsSupplier    bool
	FirstName     string
	LastName      string
	Surname       string
	Telephone     string
	Whatsapp      string
	Telegram      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
