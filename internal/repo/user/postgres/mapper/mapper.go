package mapper

import (
	"github.com/3Danger/telegram_bot/internal/repo/user"
	"github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

func UserToRepo(row query.User) *user.User {
	return &user.User{
		ID:            int(row.ID),
		HasRegistered: row.HasRegistered,
		IsSupplier:    row.IsSupplier,
		FirstName:     row.FirstName,
		LastName:      row.LastName,
		Surname:       row.Surname,
		Telephone:     row.Telephone,
		Whatsapp:      row.Whatsapp,
		Telegram:      row.Telegram,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func UserToQuery(row user.User) query.UpsertUserParams {
	return query.UpsertUserParams{
		ID:            int64(row.ID),
		HasRegistered: row.HasRegistered,
		IsSupplier:    row.IsSupplier,
		FirstName:     row.FirstName,
		LastName:      row.LastName,
		Surname:       row.Surname,
		Telephone:     row.Telephone,
		Whatsapp:      row.Whatsapp,
		Telegram:      row.Telegram,
	}
}
