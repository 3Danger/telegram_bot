package telegram

/*
Выберите действие
	1. Выложить товар
	2. Показать мои товары
	3. Мои данные
*/

type Command interface {
	GetMessage() string
	NextCommand(path string) Command
}
