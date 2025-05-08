package handlers

//
//func subHandlerUserType(ctx context.Context, r Repo, data models.Request) (models.Responses, error) {
//	const (
//		supplier = "supplier"
//		customer = "customer"
//		userType = "user_type"
//	)
//
//	u, err := r.Get(ctx, data.UserID())
//	if err != nil {
//		return nil, fmt.Errorf("getting user from cache: %w", err)
//	}
//
//	if u == nil {
//		u = &userpg.User{ID: data.UserID()}
//	}
//
//	switch data.Value(userType) {
//	case supplier:
//		u.UserType = userpg.UserTypeSupplier
//	case customer:
//		u.UserType = userpg.UserTypeCustomer
//	default:
//		return models.NewResponses(
//			data.ChatID(),
//			"Выберите тип аккаунта",
//			menu.NewInline(
//				buttons.ConstAuthChoiceImCustomer.Inline().WithValue(userType, customer),
//				buttons.ConstAuthChoiceImSupplier.Inline().WithValue(userType, supplier),
//			),
//		), nil
//	}
//
//	if err = r.Set(ctx, data.UserID(), u); err != nil {
//		return nil, fmt.Errorf("setting user to cache: %w", err)
//	}
//
//	return nil, nil
//}
//
//func subHandlerComplete(ctx context.Context, r Repo, data models.Request) (models.Responses, error) {
//	const (
//		ask     = "ask"
//		confirm = "confirm"
//		edit    = "edit"
//	)
//
//	switch data.Value(ask) {
//	case edit:
//		return subHandlerContact(ctx, a, data)
//	case confirm:
//		return models.Responses{{
//			ChatID: data.ChatID(),
//			Text:   "Сохранено!",
//			Menu:   menu.NewInline(buttons.Home, buttons.Back),
//		}}, nil
//	}
//
//	u, err := r.Get(ctx, data.UserID())
//	if err != nil {
//		return nil, fmt.Errorf("getting user from cache: %w", err)
//	}
//
//	if u == nil {
//		//TODO юзер не сохранился
//		return nil, nil
//	}
//
//	text := `Имя: ` + u.FirstName + "\n" +
//		`Фамилия: ` + u.LastName + "\n" +
//		`Телефон: ` + u.Phone + "\n" +
//		`Инфо: ` + u.Additional + "\n"
//
//	return models.Responses{{
//		ChatID: data.ChatID(),
//		Text:   "Проверьте свои свои данные\n" + text,
//		Menu: menu.NewInline(
//			buttons.ConstAuthSave.Inline().WithValue(ask, confirm),
//			buttons.ConstAuthEdit.Inline().WithValue(ask, edit),
//		),
//	}}, nil
//}
