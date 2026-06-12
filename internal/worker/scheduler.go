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
	"github.com/k4sper1love/buskr/internal/tz"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"gopkg.in/telebot.v3"
)

type Translator interface {
	T(lang, key string, args map[string]any) string
}

type Scheduler struct {
	bot      *telebot.Bot
	tr       Translator
	bookings *booking.Service
	users    *user.Service
	locs     *location.Service
	state    *redis.StateStore
}

func NewScheduler(bot *telebot.Bot, tr Translator, bookings *booking.Service, users *user.Service, locs *location.Service, state *redis.StateStore) *Scheduler {
	return &Scheduler{
		bot:      bot,
		tr:       tr,
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

		btnText := s.tr.T("ru", keys.TextWorkerReminderBtn, nil)
		btnCheckIn := menu.Data(btnText, keys.BtnBookCheckin, b.ID)
		menu.Inline(menu.Row(btnCheckIn))

		msg := s.tr.T("ru", keys.TextWorkerReminderMsg, map[string]any{
			"time": tz.In(b.StartTime).Format("15:04"),
		})

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
			msg := s.tr.T("ru", keys.TextWorkerCheckinFailMsg, nil)
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

	startHour := tz.In(b.StartTime).Hour()
	endHour := tz.In(b.EndTime).Hour()
	duration := endHour - startHour

	noiseText := s.tr.T("ru", keys.TextCommonLblNoise+"_"+string(loc.MaxNoise), nil)
	msg := s.tr.T("ru", keys.TextWorkerHotSpotMsg, map[string]any{
		"loc":   loc.Name,
		"start": fmt.Sprintf("%02d", startHour),
		"end":   fmt.Sprintf("%02d", endHour),
		"noise": noiseText,
	})

	menu := &telebot.ReplyMarkup{}
	btnText := s.tr.T("ru", keys.TextWorkerHotSpotBtn, nil)
	btnGrab := menu.Data(btnText, keys.BtnBookGrabHot, loc.ID, strconv.Itoa(startHour), strconv.Itoa(duration))
	menu.Inline(menu.Row(btnGrab))

	for _, u := range activeUsers {
		if u.ID == b.UserID {
			continue
		}

		if !booking.IsNoiseCompatible(u.NoiseLevel, loc.MaxNoise) {
			continue
		}

		_, err := s.bot.Send(&telebot.User{ID: u.TelegramID}, msg, menu, telebot.ModeMarkdown)
		if err != nil {
			log.Printf("Failed to send hot spot to %d: %v", u.TelegramID, err)
		}
	}
}
