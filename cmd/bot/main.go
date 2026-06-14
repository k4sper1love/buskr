package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/k4sper1love/buskr/internal/config"
	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/i18n"
	"github.com/k4sper1love/buskr/internal/infrastructure/postgres"
	redisInfra "github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"github.com/k4sper1love/buskr/internal/transport/telegram"
	"github.com/k4sper1love/buskr/internal/tz"
	"github.com/k4sper1love/buskr/internal/worker"
	"github.com/k4sper1love/buskr/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustLoad()

	logger.Setup(cfg.Env)
	slog.Info("starting bot", "env", cfg.Env)

	// timezone
	if err := tz.Init(cfg.Timezone); err != nil {
		slog.Error("failed to initialize timezone", "err", err)
		os.Exit(1)
	}

	// database
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		slog.Error("failed to open postgres", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("postgres ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("postgres connected")

	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Redis.Host,
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		DialTimeout: cfg.Redis.DialTimeout,
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		slog.Error("redis ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("redis connected")

	// translator
	tr := i18n.NewTranslator()

	// repositories
	stateStore := redisInfra.NewStateStore(rdb)
	userRepo := postgres.NewUserRepository(db)
	locRepo := postgres.NewLocationRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)

	// services
	userService := user.NewService(userRepo, user.Config{
		MaxKarma:         cfg.Karma.MaxKarma,
		MinKarma:         cfg.Karma.MinKarma,
		KarmaReward:      cfg.Karma.KarmaReward,
		KarmaPenalty:     cfg.Karma.KarmaPenalty,
		WarningThreshold: cfg.Karma.WarningThreshold,
	})
	locService := location.NewService(locRepo)
	bookingService := booking.NewService(bookingRepo, userService, locService, booking.Config{
		MaxActiveBookings:           cfg.Booking.MaxActive,
		AdminMaxActiveBookings:      cfg.Booking.AdminMaxActive,
		MaxBookingsPerLocation:      cfg.Booking.MaxPerLocationAtDay,
		AdminMaxBookingsPerLocation: cfg.Booking.AdminMaxPerLocationAtDay,
		MaxAdvanceDays:              cfg.Booking.MaxAdvanceDays,
		AdjacencyRadius:             cfg.Booking.AdjacencyRadius,
		AdminBypassNoiseLimits:      cfg.Booking.AdminBypassNoiseLimits,
		AdminBypassUserOverlap:      cfg.Booking.AdminBypassUserOverlap,
		AdminBypassNoisyNeighbor:    cfg.Booking.AdminBypassNoisyNeighbor,
	})

	// bot
	bot, err := telegram.NewBot(
		&cfg.Telegram,
		tr,
		stateStore,
		userService,
		bookingService,
		locService,
		cfg.Booking.MaxAdvanceDays,
		cfg.Maps.GoogleAPIKey,
		cfg.Maps.WebAppURL,
		cfg.Booking.MinHour,
		cfg.Booking.MaxHourWeekday,
		cfg.Booking.MaxHourWeekend,
		cfg.Booking.MaxDurationHours,
		cfg.Booking.AdminMaxDurationHours,
	)
	if err != nil {
		slog.Error("failed to initialize bot", "err", err)
		os.Exit(1)
	}

	// workers
	scheduler := worker.NewScheduler(
		bot.GetTelebot(),
		tr,
		bookingService,
		userService,
		locService,
		stateStore,
		cfg.Booking.EnableHotSlots,
		cfg.Booking.EnableNoShowCancel,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	// graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		slog.Info("shutting down gracefully")
		cancel()
		bot.GetTelebot().Stop()
	}()

	bot.Start()

	slog.Info("app stopped")
}
