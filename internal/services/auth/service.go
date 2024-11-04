package auth

import (
	"context"
	"fmt"
	"regexp"

	r "github.com/3Danger/telegram_bot/internal/repo"
	"github.com/3Danger/telegram_bot/internal/repo/user"
)

type repo struct {
	session   r.Repo[user.User]
	permanent r.Repo[user.User]
}

type Service struct {
	repo repo
	word *regexp.Regexp
	//phone      *regexp.Regexp
	//tgUsername *regexp.Regexp
}

func NewService(permanentRepo, sessionRepo r.Repo[user.User]) *Service {
	word, err := regexp.Compile(`^[a-zA-Zа-яА-Я]+$`)
	if err != nil {
		panic(err)
	}
	// TODO пофиксить
	phone, err := regexp.Compile(`^(\+?\d{1,3})\D?(\d{1,3}\D?)\D?(\d{2,4}.?){1,2}(\d{2,4})$`)
	if err != nil {
		panic(err)
	}
	_ = phone
	//tgUsername, err := regexp.Compile(`(?:@|(?:(?:(?:https?://)?t(?:elegram)?)\.me\/))(\w{4,})$`)
	//if err != nil {
	//	panic(err)
	//}

	return &Service{
		word: word,
		//phone:      phone,
		//tgUsername: tgUsername,
		repo: repo{
			permanent: permanentRepo,
			session:   sessionRepo,
		},
	}
}

func (s *Service) SaveToPermanent(ctx context.Context, userID int64) error {
	u, err := s.repo.session.Get(ctx, userID)
	if err != nil {
		return fmt.Errorf("getting user from sesson: %w", err)
	}

	if u == nil {
		return ErrUserNotFound
	}

	if err = s.repo.permanent.Set(ctx, userID, *u); err != nil {
		return fmt.Errorf("saving user to permanent repo: %w", err)
	}

	return nil
}

func (s *Service) AddUserID(ctx context.Context, userId int64) error {
	return s.SetSomething(ctx, userId, func(u *user.User) error {
		return nil
	})
}

func (s *Service) AddFirstName(ctx context.Context, userId int64, data string) error {
	if len(s.word.FindStringSubmatch(data)) != 1 {
		return ErrWrongName
	}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.FirstName = data

		return nil
	})
}

func (s *Service) AddLastName(ctx context.Context, userId int64, data string) error {
	if len(s.word.FindStringSubmatch(data)) != 1 {
		return ErrWrongName
	}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.LastName = data

		return nil
	})
}

func (s *Service) AddSurname(ctx context.Context, userId int64, data string) error {
	if len(s.word.FindStringSubmatch(data)) != 1 {
		return ErrWrongName
	}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.Surname = data

		return nil
	})
}

func (s *Service) AddPhone(ctx context.Context, userId int64, data string) error {
	//if len(s.phone.FindStringSubmatch(data)) != 1 {
	//	return ErrWrongPhone
	//}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.Phone = data

		return nil
	})
}

func (s *Service) AddTelegram(ctx context.Context, userId int64, data string) error {
	//if len(s.tgUsername.FindStringSubmatch(data)) != 1 {
	//	return ErrWrongNickName
	//}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.Telegram = data

		return nil
	})
}

func (s *Service) AddWhatsapp(ctx context.Context, userId int64, data string) error {
	//if len(s.phone.FindStringSubmatch(data)) != 1 {
	//	return ErrWrongPhone
	//}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.Whatsapp = data

		return nil
	})
}

func (s *Service) AddUserType(ctx context.Context, userId int64, data string) error {
	if !user.Type(data).Valid() {
		return ErrWrongUserType
	}

	return s.SetSomething(ctx, userId, func(u *user.User) error {
		u.Type = user.Type(data)

		return nil
	})
}

func (s *Service) SetSomething(ctx context.Context, userId int64, setter func(u *user.User) error) error {
	u, err := s.repo.session.Get(ctx, userId)
	if err != nil {
		return fmt.Errorf("getting user from session repo: %w", err)
	}
	if u == nil {
		u = &user.User{
			ID: userId,
		}
	}
	if err = setter(u); err != nil {
		return fmt.Errorf("setting values: %w", err)
	}
	if err = s.repo.session.Set(ctx, userId, *u); err != nil {
		return fmt.Errorf("setting user to session repo: %w", err)
	}

	return nil
}

func (s *Service) RegisteredSession(ctx context.Context, userID int64) (bool, error) {
	u, err := s.repo.session.Get(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("getting user from session repo: %w", err)
	}
	if u == nil {
		return false, nil
	}

	return hasRegistered(u), nil
}

func (s *Service) GetFromSession(ctx context.Context, userID int64) (*user.User, error) {
	u, err := s.repo.session.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting user from session repo: %w", err)
	}

	return u, nil
}

func (s *Service) RegisteredPermanent(ctx context.Context, userID int64) (bool, error) {
	u, err := s.repo.permanent.Get(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("getting user from session repo: %w", err)
	}

	return u != nil, nil
}

type State string

const (
	StateNew       = State("StateNew")
	StateType      = State("StateType")
	StateFirstName = State("StateFirstName")
	StateLastName  = State("StateLastName")
	StateSurname   = State("StateSurname")
	StatePhone     = State("StatePhone")
	StateWhatsapp  = State("StateWhatsapp")
	StateTelegram  = State("StateTelegram")
	StateCompleted = State("StateCompleted")
)

func (s *Service) GetSessionState(ctx context.Context, userId int64) (State, error) {
	u, err := s.repo.session.Get(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("getting user from session repo: %w", err)
	}

	return getState(u), nil
}

func getState(u *user.User) State {
	switch {
	case u == nil:
		return StateNew
	case u.ID == 0:
		return StateNew
	case u.Type == "" || u.Type == user.TypeUndefined:
		return StateType
	case u.FirstName == "":
		return StateFirstName
	case u.LastName == "":
		return StateLastName
	case u.Surname == "":
		return StateSurname
	case u.Phone == "":
		return StatePhone
	case u.Whatsapp == "":
		return StateWhatsapp
	case u.Telegram == "":
		return StateTelegram
	}

	return StateCompleted
}

func hasRegistered(u *user.User) bool {
	return getState(u) == StateCompleted
}
