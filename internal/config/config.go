package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	Timezone string `env:"TIMEZONE" env-default:"Asia/Almaty"`

	// transport
	Telegram TelegramConfig

	// infra
	Database DatabaseConfig
	Redis    RedisConfig

	// services
	Booking BookingConfig
	Maps    MapsConfig
	Karma   KarmaConfig
}

type TelegramConfig struct {
	BotToken                string `env:"TELEGRAM_BOT_TOKEN" env-required:"true"`
	DefaultLang             string `env:"TELEGRAM_DEFAULT_LANG" env-default:"ru"`
	AdminLang               string `env:"TELEGRAM_ADMIN_LANG" env-default:"ru"`
	AdminChatID             int64  `env:"TELEGRAM_ADMIN_CHAT_ID" env-required:"true"`
	AdminThreadApplications int    `env:"TELEGRAM_ADMIN_THREAD_APPLICATIONS" env-default:"0"`
	AdminThreadUpgrades     int    `env:"TELEGRAM_ADMIN_THREAD_UPGRADES" env-default:"0"`
	PublicChatID            int64  `env:"TELEGRAM_PUBLIC_CHAT_ID" env-default:"0"`
	SupportContact          string `env:"TELEGRAM_SUPPORT_CONTACT" env-default:"@buskr_support"`
	InviteTTLHours          int    `env:"INVITE_TTL_HOURS" env-default:"24"`
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

type BookingConfig struct {
	MaxActive                int  `env:"BOOKING_MAX_ACTIVE" env-default:"2"`
	AdminMaxActive           int  `env:"BOOKING_ADMIN_MAX_ACTIVE" env-default:"-1"`
	MaxPerLocationAtDay      int  `env:"BOOKING_MAX_PER_LOCATION_AT_DAY" env-default:"1"`
	AdminMaxPerLocationAtDay int  `env:"BOOKING_ADMIN_MAX_PER_LOCATION_AT_DAY" env-default:"-1"`
	MaxAdvanceDays           int  `env:"BOOKING_MAX_ADVANCE_DAYS" env-default:"5"`
	AdjacencyRadius          int  `env:"BOOKING_ADJACENCY_RADIUS" env-default:"100"`
	EnableHotSlots           bool `env:"BOOKING_ENABLE_HOT_SLOTS" env-default:"true"`
	EnableNoShowCancel       bool `env:"BOOKING_ENABLE_NO_SHOW_CANCEL" env-default:"true"`
	AdminBypassNoiseLimits   bool `env:"BOOKING_ADMIN_BYPASS_NOISE_LIMITS" env-default:"true"`
	AdminBypassUserOverlap   bool `env:"BOOKING_ADMIN_BYPASS_USER_OVERLAP" env-default:"true"`
	AdminBypassNoisyNeighbor bool `env:"BOOKING_ADMIN_BYPASS_NOISY_NEIGHBOR" env-default:"true"`
	MinHour                  int  `env:"BOOKING_MIN_HOUR" env-default:"10"`
	MaxHourWeekday           int  `env:"BOOKING_MAX_HOUR_WEEKDAY" env-default:"22"`
	MaxHourWeekend           int  `env:"BOOKING_MAX_HOUR_WEEKEND" env-default:"23"`
	MaxDurationHours         int  `env:"BOOKING_MAX_DURATION_HOURS" env-default:"3"`
	AdminMaxDurationHours    int  `env:"BOOKING_ADMIN_MAX_DURATION_HOURS" env-default:"-1"`
}

type MapsConfig struct {
	GoogleAPIKey string `env:"GOOGLE_MAPS_API_KEY"`
	WebAppURL    string `env:"MAP_WEB_APP_URL" env-default:""`
}

type KarmaConfig struct {
	MaxKarma         int `env:"KARMA_MAX" env-default:"100"`
	MinKarma         int `env:"KARMA_MIN" env-default:"0"`
	KarmaReward      int `env:"KARMA_REWARD" env-default:"5"`
	KarmaPenalty     int `env:"KARMA_PENALTY" env-default:"20"`
	PerfectThreshold int `env:"KARMA_PERFECT_THRESHOLD" env-default:"100"`
	WarningThreshold int `env:"KARMA_WARNING_THRESHOLD" env-default:"50"`
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
