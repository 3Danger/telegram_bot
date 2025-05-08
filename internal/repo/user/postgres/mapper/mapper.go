package mapper

import (
	"github.com/3Danger/telegram_bot/internal/models"
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

func UserRepoToModel(user *userpg.User) *models.User {
	return &models.User{
		ID:         int(user.ID),
		UserType:   models.UserType(user.UserType),
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		Additional: user.Additional,
	}
}
