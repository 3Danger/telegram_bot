package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

var config *Config

func New() (*Config, error) {
	config = new(Config)

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getting current working directory: %v", err))
	}

	envFilePath := filepath.Join(cwd, ".env")

	if err = godotenv.Load(envFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env file: %v", err)
	}

	if err = envconfig.Process("", config); err != nil {
		return nil, fmt.Errorf("processing env vars: %w", err)
	}

	return config, nil
}

type Config struct {
	App      App
	Telegram Telegram
	Postgres Postgres
}

type App struct {
	LogLevel zerolog.Level `default:"info" envconfig:"LOG_LEVEL"`
}

type Telegram struct {
	Token string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	Debug bool   `default:"true"                 envconfig:"TELEGRAM_BOT_DEBUG_MODE"`
}

type Postgres struct {
	Host             string        `envconfig:"POSTGRES_HOST"     required:"true"`
	Port             string        `envconfig:"POSTGRES_PORT"     required:"true"`
	Username         string        `envconfig:"POSTGRES_USERNAME" required:"true"`
	Database         string        `envconfig:"POSTGRES_DATABASE" required:"true"`
	Password         string        `envconfig:"POSTGRES_PASSWORD" required:"true"`
	SSL              string        `default:"disable"             envconfig:"POSTGRES_SSL"`
	OperationTimeout time.Duration `default:"60s"                 envconfig:"POSTGRES_OPERATION_TIMEOUT"`
}

func (p Postgres) DSN() string {
	dsn := strings.Builder{}
	dsn.WriteString("postgres://")

	if p.Username != "" {
		dsn.WriteString(p.Username)

		if p.Password != "" {
			dsn.WriteString(fmt.Sprintf(":%s", p.Password))
		}

		dsn.WriteString("@")
	}

	dsn.WriteString(fmt.Sprintf("%s/", net.JoinHostPort(p.Host, p.Port)))

	if p.Database != "" {
		dsn.WriteString(p.Database)
	}

	dsn.WriteString(fmt.Sprintf("?sslmode=%s", p.SSL))

	return dsn.String()
}
