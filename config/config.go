package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App      App `yaml:"app"`
		Postgres Postgres
		Bot      Bot
	}

	App struct {
		Name    string `yaml:"name" env:"APP_NAME" env-default:"telegram-bot"`
		Version string `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	}

	Bot struct {
		Token string `env:"BOT_TOKEN" env-required:"true"`
	}

	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-required:"true"`
		Port     string `env:"POSTGRES_PORT" env-required:"true"`
		User     string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Database string `env:"POSTGRES_DATABASE" env-required:"true"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	// Read sensitive data from environment variables.
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
