package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/infrastructure/postgres"
	redisInfra "github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"github.com/k4sper1love/buskr/internal/transport/telegram"
	"github.com/k4sper1love/buskr/internal/worker"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("Starting Buskr App...")

	botToken := "your_bot_token_here"
	pgConnStr := "your_postgres_dsn_here"
	redisAddr := "your_redis_addr_here"

	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Postgres is dead: %v", err)
	}
	log.Println("PostgreSQL connected successfully!")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis is dead: %v", err)
	}
	log.Println("Redis connected successfully!")

	stateStore := redisInfra.NewStateStore(rdb)
	userRepo := postgres.NewUserRepository(db)
	locRepo := postgres.NewLocationRepository(db)
	bookingRepo := postgres.NewBookingRepository(db)

	userService := user.NewService(userRepo)
	locService := location.NewService(locRepo)

	bookingService := booking.NewService(bookingRepo, userService, locService)

	bot, err := telegram.NewBot(botToken, userService, bookingService, locService, stateStore)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	scheduler := worker.NewScheduler(bot.GetTelebot(), bookingService, userService, locService, stateStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler.Start(ctx)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutting down gracefully...")
		cancel()
		bot.GetTelebot().Stop()
	}()

	bot.Start()

	log.Println("App stopped.")
}
