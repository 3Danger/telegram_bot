package constants

const (
	Home          = "🏠На главную"
	Auth          = "⚙️Регистрация"
	ImSupplier    = "📦Я продавец"
	ImCustomer    = "💝Я покупатель"
	Back          = "⬅️ Назад"
	AuthConfirm   = "✅ Подтвердить"
	AuthEditName  = "📝 Изменить ФИО"
	AuthEditPhone = "📱 Изменить телефон"

	CustomerShowItems = "📦Показать товары"
	SupplierShowItems = "📦Показать мои товары"
	SupplierPostItems = "➕Добавить товар"
)

func IsValid(route string) bool {
	_, ok := validRoutes[route]

	return ok
}

var validRoutes = map[string]struct{}{
	Home:              {},
	Auth:              {},
	ImSupplier:        {},
	ImCustomer:        {},
	Back:              {},
	AuthConfirm:       {},
	AuthEditName:      {},
	AuthEditPhone:     {},
	CustomerShowItems: {},
	SupplierShowItems: {},
}
