package config

type Config struct {
	Telegram Telegram
}

type Telegram struct {
	Token string `envconfig:"TELEGRAM_BOT_TOKEN"`
}
