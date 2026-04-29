package worker

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"gopkg.in/telebot.v3"
)

type Scheduler struct {
	bot      *telebot.Bot
	bookings *booking.Service
	users    *user.Service
	locs     *location.Service
	state    *redis.StateStore
}

func NewScheduler(bot *telebot.Bot, bookings *booking.Service, users *user.Service, locs *location.Service, state *redis.StateStore) *Scheduler {
	return &Scheduler{
		bot:      bot,
		bookings: bookings,
		users:    users,
		locs:     locs,
		state:    state,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	log.Println("Background scheduler started...")

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Scheduler stopped gracefully")
				return
			case <-ticker.C:
				s.processReminders(ctx)
				s.processCheckins(ctx)
				s.processCompletions(ctx)
			}
		}
	}()
}

func (s *Scheduler) processReminders(ctx context.Context) {
	upcoming, err := s.bookings.GetUpcomingForReminder(ctx, 60*time.Minute)
	if err != nil {
		log.Printf("Scheduler Reminder error: %v", err)
		return
	}

	for _, b := range upcoming {
		cacheKey := fmt.Sprintf("reminder_sent:%s", b.ID)
		var alreadySent bool
		_ = s.state.GetData(ctx, 0, cacheKey, &alreadySent)
		if alreadySent {
			continue
		}

		targetUser, err := s.users.GetByID(ctx, b.UserID)
		if err != nil {
			continue
		}

		menu := &telebot.ReplyMarkup{}
		btnCheckIn := menu.Data("📍 Я на месте", "btn_checkin", b.ID)
		menu.Inline(menu.Row(btnCheckIn))

		msg := fmt.Sprintf("🔔 **Напоминание!**\n\nВаше выступление начинается через час (в %02d:00).\n\n⚠️ Обязательно нажмите кнопку ниже, когда прибудете на точку, иначе бронь сгорит.", b.StartTime.Hour())

		_, err = s.bot.Send(&telebot.User{ID: targetUser.TelegramID}, msg, menu, telebot.ModeMarkdown)
		if err == nil {
			_ = s.state.SetData(ctx, 0, cacheKey, true, 2*time.Hour)
		}
	}
}

func (s *Scheduler) processCheckins(ctx context.Context) {
	expired, err := s.bookings.GetPendingForCheckinTimeout(ctx, 15*time.Minute)
	if err != nil {
		log.Printf("Scheduler Checkin error: %v", err)
		return
	}

	for _, b := range expired {
		err := s.bookings.MarkNoShow(ctx, b.ID)
		if err != nil {
			log.Printf("Failed to mark no-show for booking %s: %v", b.ID, err)
			continue
		}

		_ = s.users.RecordNoShow(ctx, b.UserID)

		targetUser, _ := s.users.GetByID(ctx, b.UserID)
		if targetUser != nil {
			msg := "❌ **Ваша бронь аннулирована.**\n\nВы не подтвердили присутствие в течение 15 минут после начала слота. Ваш рейтинг понижен."
			_, _ = s.bot.Send(&telebot.User{ID: targetUser.TelegramID}, msg, telebot.ModeMarkdown)
		}

		s.broadcastHotSpot(ctx, b)
	}
}

func (s *Scheduler) processCompletions(ctx context.Context) {
	finished, err := s.bookings.GetActiveForCompletion(ctx)
	if err != nil {
		log.Printf("Scheduler Completion error: %v", err)
		return
	}

	for _, b := range finished {
		err := s.bookings.CompleteBooking(ctx, b.ID)
		if err != nil {
			log.Printf("Failed to complete booking %s: %v", b.ID, err)
			continue
		}
	}
}

func (s *Scheduler) broadcastHotSpot(ctx context.Context, b *booking.Booking) {
	loc, err := s.locs.GetByID(ctx, b.LocationID)
	if err != nil || loc.Status != location.StatusActive {
		return
	}

	activeUsers, err := s.users.GetActiveUsers(ctx)
	if err != nil || len(activeUsers) == 0 {
		return
	}

	locTz, _ := time.LoadLocation("Asia/Almaty")
	startHour := b.StartTime.In(locTz).Hour()
	endHour := b.EndTime.In(locTz).Hour()
	duration := endHour - startHour

	msg := fmt.Sprintf(
		"🔥 **ГОРЯЩИЙ СЛОТ!**\n\n"+
			"Музыкант не явился на точку, она свободна прямо сейчас!\n\n"+
			"📍 **Локация:** %s\n"+
			"⏰ **Время:** %02d:00 - %02d:00\n"+
			"🔊 **Макс. шум:** %s\n\n"+
			"_Кто первый нажмет кнопку, тот и заберет слот._",
		loc.Name, startHour, endHour, loc.MaxNoise,
	)

	menu := &telebot.ReplyMarkup{}
	btnGrab := menu.Data("⚡ Забрать слот", "btn_grab_hot", loc.ID, strconv.Itoa(startHour), strconv.Itoa(duration))
	menu.Inline(menu.Row(btnGrab))

	for _, u := range activeUsers {
		if u.ID == b.UserID {
			continue
		}

		if u.NoiseLevel.Weight() > loc.MaxNoise.Weight() {
			continue
		}

		_, err := s.bot.Send(&telebot.User{ID: u.TelegramID}, msg, menu, telebot.ModeMarkdown)
		if err != nil {
			log.Printf("Failed to send hot spot to %d: %v", u.TelegramID, err)
		}
	}
}
