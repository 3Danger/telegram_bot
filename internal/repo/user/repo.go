package user

import "context"

type Repo interface {
	User(ctx context.Context, userID int) (*User, error)
	State(ctx context.Context, chatID int) (State, error)
}

type User struct {
	IsSupplier bool
	FIO        struct {
		FirstName string
		LastName  string
		Surname   string
	}
	Contact struct {
		Telephone string
		Whatsapp  string
		Telegram  string
	}
}

type name struct{}

type State string

const (
	StateStart State = "start"
)
