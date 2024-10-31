package mapper

import (
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

func UserToRepo(row query.User) *user.User {
	return &user.User{
		ID:         row.ID,
		IsSupplier: row.IsSupplier,
		FirstName:  row.FirstName,
		LastName:   row.LastName,
		Surname:    row.Surname,
		Phone:      row.Phone,
		Whatsapp:   row.Whatsapp,
		Telegram:   row.Telegram,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func UserToQuery(row user.User) query.CreateUserParams {
	return query.CreateUserParams{
		ID:         row.ID,
		IsSupplier: row.IsSupplier,
		FirstName:  row.FirstName,
		LastName:   row.LastName,
		Surname:    row.Surname,
	}
}
