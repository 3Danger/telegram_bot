package config

import "time"

type Config struct {
	Telegram Telegram
	Postgres Postgres
}

type Telegram struct {
	Token string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	Debug bool   `envconfig:"TELEGRAM_BOT_DEBUG_MODE" default:"true" required:"false"`
}

type Postgres struct {
	Host             string        `envconfig:"POSTGRES_HOST" required:"true"`
	Username         string        `envconfig:"POSTGRES_USERNAME" required:"true"`
	Database         string        `envconfig:"POSTGRES_DATABASE" required:"true"`
	Password         string        `envconfig:"POSTGRES_PASSWORD" required:"true"`
	OperationTimeout time.Duration `envconfig:"POSTGRES_OPERATION_TIMEOUT" default:"60s"`
}
