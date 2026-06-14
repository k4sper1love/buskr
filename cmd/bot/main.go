package main

import (
	"context"
	"database/sql"
	"log"
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
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustLoad()

	log.Print("Starting bot in env:", cfg.Env)

	// timezone
	if err := tz.Init(cfg.Timezone); err != nil {
		log.Fatalf("Failed to initialize timezone: %v", err)
	}

	// database
	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Postgres is dead: %v", err)
	}
	log.Println("PostgreSQL connected successfully!")

	// redis
	rdb := redis.NewClient(&redis.Options{
		Addr:        cfg.Redis.Host,
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		DialTimeout: cfg.Redis.DialTimeout,
	})
	defer rdb.Close()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis is dead: %v", err)
	}
	log.Println("Redis connected successfully!")

	// translator
	tr := i18n.NewTranslator()

	// repositories
	stateStore := redisInfra.NewStateStore(rdb)
	userRepo := postgres.NewUserRepository(db)
	locRepo := postgres.NewLocationRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)

	// servives
	userService := user.NewService(userRepo, user.Config{
		MaxKarma:         cfg.Karma.MaxKarma,
		MinKarma:         cfg.Karma.MinKarma,
		KarmaReward:      cfg.Karma.KarmaReward,
		KarmaPenalty:     cfg.Karma.KarmaPenalty,
		WarningThreshold: cfg.Karma.WarningThreshold,
	})
	locService := location.NewService(locRepo)
	bookingService := booking.NewService(bookingRepo, userService, locService, booking.Config{
		MaxActiveBookings:      cfg.Booking.MaxActive,
		AdminMaxActiveBookings: cfg.Booking.AdminMaxActive,
		MaxBookingsPerLocation: cfg.Booking.MaxPerLocationAtDay,
		MaxAdvanceDays:         cfg.Booking.MaxAdvanceDays,
		AdjacencyRadius:        cfg.Booking.AdjacencyRadius,
	})

	// bot
	bot, err := telegram.NewBot(&cfg.Telegram, tr, stateStore, userService, bookingService, locService, cfg.Booking.MaxAdvanceDays, cfg.Maps.GoogleAPIKey, cfg.Maps.WebAppURL)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
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
		log.Println("Shutting down gracefully...")
		cancel()
		bot.GetTelebot().Stop()
	}()

	// start bot
	bot.Start()

	// when apllication was stopped
	log.Println("App stopped.")
}
