package auth

import "errors"

var (
	ErrWrongName     = newValidationErr("Не неверно задано имя")
	ErrWrongUserType = newValidationErr("Не неверно задан тип пользователя")
	ErrWrongPhone    = newValidationErr("Не неверно задан телефон")
	ErrWrongNickName = newValidationErr("Не неверно задан ник")

	ErrUserNotFound = newValidationErr("Пользователь не найден")
)

type ValidationErr struct {
	msg string
}

func (v *ValidationErr) Error() string {
	return v.msg
}

func newValidationErr(msg string) error {
	return &ValidationErr{msg: msg}
}

func IsValidationErr(err error) bool {
	var v *ValidationErr

	return errors.As(err, &v)
}
