// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package query

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type UserType string

const (
	UserTypeUndefined UserType = "undefined"
	UserTypeSupplier  UserType = "supplier"
	UserTypeCustomer  UserType = "customer"
)

func (e *UserType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserType(s)
	case string:
		*e = UserType(s)
	default:
		return fmt.Errorf("unsupported scan type for UserType: %T", src)
	}
	return nil
}

type NullUserType struct {
	UserType UserType
	Valid    bool // Valid is true if UserType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserType) Scan(value interface{}) error {
	if value == nil {
		ns.UserType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserType), nil
}

// Все пользователи
type User struct {
	ID int64
	// Тип пользователя
	UserType UserType
	// Имя
	FirstName string
	// Фамилия
	LastName string
	// Отчество
	Surname string
	// Телефон
	Phone string
	// Номер whatsapp
	Whatsapp string
	// Ник в телеграмме
	Telegram  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
