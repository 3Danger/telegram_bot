package telegram

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	tele "gopkg.in/telebot.v4"

	"github.com/3Danger/telegram_bot/internal/telegram/constants"
)

func (t *Telegram) handlerHome(c tele.Context) error {
	u, err := t.repo.user.User(getContext(c), c.Sender().ID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}
	if u == nil {
		return c.Send(
			"Добро пожаловать!\nДля пользования необходимо регистрация",
			createMenu(constants.Auth),
		)
	}

	if u.IsSupplier {
		return t.handlerSupplierHome(c)
	}

	return t.handlerCustomerHome(c)
}

// Добавить товары TODO на доработке
func (t *Telegram) handlerSupplierPostItems(c tele.Context) error {
	msg := c.Message()
	if msg != nil || msg.Media() != nil {
		return c.Send("Пришлите фото/видео товара")
	}

	media := msg.Media()

	file, err := c.Bot().File(media.MediaFile())
	if err != nil {
		return fmt.Errorf("getting photo file: %w", err)
	}
	defer file.Close()

	switch media := media.(type) {
	case *tele.Photo:
		if err := t.v.ValidatePhoto(media); err != nil {
			return err // TODO пояснить ошибку для пользователя
		}

	case *tele.Video:
		if err := t.v.ValidateVideo(media); err != nil {
			return err // TODO пояснить ошибку для пользователя
		}

	case *tele.VideoNote:
		if err := t.v.ValidateVideoNote(media); err != nil {
			return err // TODO пояснить ошибку для пользователя
		}
	}

	// Проверка что файл не битый для дебаша
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting currenc working dir: %w", err)
	}

	f, err := os.Create(filepath.Join(path, "file.jpg"))
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	b, err := io.Copy(f, file)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	fmt.Println("PHOTO", b)
	//msg.Photo.
	//msg.Photo.File.FileReader
	//os.Create("tmp")

	return c.Send("Пришлите пожалуйста фото")
}

// Показать мои товары
func (t *Telegram) handlerSupplierShowItems(c tele.Context) error {
	return c.Send("", createMenu(constants.Back))
}

func (t *Telegram) handlerSupplierHome(c tele.Context) error {
	return c.Send("", createMenu(
		//constants.SupplierShowItems,
		constants.SupplierPostItems,
	))
}

func (t *Telegram) handlerCustomerHome(c tele.Context) error {
	return c.Reply("", createMenu(constants.CustomerShowItems))
}
