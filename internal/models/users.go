package models

import (
	userpg "github.com/3Danger/telegram_bot/internal/repo/user/postgres/query"
)

type User struct {
	ID         int      `json:"id"`
	UserType   UserType `json:"user_type"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Phone      string   `json:"phone"`
	Additional string   `json:"additional"`
}

type UserType string

const (
	UserTypeSupplier = UserType(userpg.UserTypeSupplier)
	UserTypeCustomer = UserType(userpg.UserTypeCustomer)
	UserTypeUnknown  = UserType(userpg.UserTypeUnknown)
)
