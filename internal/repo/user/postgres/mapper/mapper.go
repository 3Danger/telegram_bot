package mapper

import (
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

func UserToRepo(row query.User) *user.User {
	return &user.User{
		ID:        row.ID,
		Type:      user.Type(row.UserType),
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Surname:   row.Surname,
		Phone:     row.Phone,
		Whatsapp:  row.Whatsapp,
		Telegram:  row.Telegram,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func UserToQuery(userID int64, row user.User) query.SetParams {
	return query.SetParams{
		ID:        userID,
		UserType:  query.UserType(row.Type),
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Surname:   row.Surname,
		Phone:     row.Phone,
		Whatsapp:  row.Whatsapp,
		Telegram:  row.Telegram,
	}
}
