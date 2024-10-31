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
	Repo     Repo
}

type App struct {
	LogLevel zerolog.Level `default:"info" envconfig:"LOG_LEVEL"`
}

type Telegram struct {
	Token string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	Debug bool   `default:"true"                 envconfig:"TELEGRAM_BOT_DEBUG_MODE"`
}

type Repo struct {
	InMemory
	Postgres
}

type InMemory struct {
	MaxItems int `envconfig:"REPO_IN_MEMORY_MAX_ITEMS" default:"100"`
}

type Postgres struct {
	Host     string `envconfig:"REPO_POSTGRES_HOST"     required:"true"`
	Port     string `envconfig:"REPO_POSTGRES_PORT"     required:"true"`
	Username string `envconfig:"REPO_POSTGRES_USERNAME" required:"true"`
	Database string `envconfig:"REPO_POSTGRES_DATABASE" required:"true"`
	Password string `envconfig:"REPO_POSTGRES_PASSWORD" required:"true"`
	SSL      string `envconfig:"REPO_POSTGRES_SSL"      default:"disable"`

	OperationTimeout time.Duration `envconfig:"REPO_POSTGRES_OPERATION_TIMEOUT" default:"60s"`
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
