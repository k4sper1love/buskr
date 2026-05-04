package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	Telegram TelegramConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type TelegramConfig struct {
	BotToken    string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	DefaultLang string `env:"TELEGRAM_DEFAULT_LANG" env-default:"ru"`
	AdminChatID int64  `env:"TELEGRAM_ADMIN_CHAT_ID" env-required:"true"`
	AdminLang   string `env:"TELEGRAM_ADMIN_LANG" env-default:"ru"`
}

type DatabaseConfig struct {
	DSN string `env:"DATABASE_DSN" env-required:"true"`
}

type RedisConfig struct {
	Host        string        `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password    string        `env:"REDIS_PASSWORD"`
	DB          int           `env:"REDIS_DB" env-default:"0"`
	DialTimeout time.Duration `env:"REDIS_DIAL_TIMEOUT" env-default:"5s"`
}

// singleton pattern
var instance *Config

func MustLoad() *Config {
	if instance != nil {
		return instance
	}

	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	instance = &cfg

	return instance
}
