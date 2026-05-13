package telegram

import (
	"log"
	"time"

	"github.com/k4sper1love/buskr/internal/config"
	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/i18n"
	"github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers/callbacks"
	"github.com/k4sper1love/buskr/internal/transport/telegram/notifier"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/admin"
	"github.com/k4sper1love/buskr/internal/usecase/adminloc"
	"github.com/k4sper1love/buskr/internal/usecase/adminuser"
	"github.com/k4sper1love/buskr/internal/usecase/auth"
	bookingUsecase "github.com/k4sper1love/buskr/internal/usecase/booking"
	"github.com/k4sper1love/buskr/internal/usecase/onboarding"
	"github.com/k4sper1love/buskr/internal/usecase/profile"

	"gopkg.in/telebot.v3"
)

type Bot struct {
	bot      *telebot.Bot
	cfg      *config.TelegramConfig
	state    *redis.StateStore
	renderer *render.Renderer
	users    *user.Service

	fsm      *handlers.FSM
	handlers *Handlers
}

type Handlers struct {
	auth       *callbacks.Auth
	onboarding *callbacks.Onboarding
	booking    *callbacks.Booking
	admin      *callbacks.Admin
	adminLoc   *callbacks.AdminLoc
	adminUser  *callbacks.AdminUser
	profile    *callbacks.Profile
}

func NewBot(
	cfg *config.TelegramConfig,
	tr *i18n.Translator,
	state *redis.StateStore,
	users *user.Service,
	bookings *booking.Service,
	locs *location.Service,
	maxAdvanceDays int,
) (*Bot, error) {
	// create bot instance
	pref := telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	// renderer for bot responses
	renderer := render.NewRenderer(tr, cfg.DefaultLang)

	// admin notifier
	notifier := notifier.NewNotifier(b, tr, cfg.AdminChatID, cfg.AdminLang)

	// usecases
	authUc := auth.NewUsecase(state, users, 1*time.Hour)
	bookingUc := bookingUsecase.NewUsecase(state, locs, bookings, 1*time.Hour, maxAdvanceDays)
	adminUc := admin.NewUsecase(state, users, 1*time.Hour, 1)
	adminlocUc := adminloc.NewUsecase(state, locs, bookings, 1*time.Hour)
	adminuserUc := adminuser.NewUsecase(state, users, 1*time.Hour)
	profileUc := profile.NewUsecase(state, users, notifier, 1*time.Hour)
	onboardingUc := onboarding.NewUsecase(state, users, notifier, 1*time.Hour)

	// fsm handlers (text/video/location messages)
	fsm := &handlers.FSM{
		State:      state,
		Renderer:   renderer,
		Onboarding: onboardingUc,
		Booking:    bookingUc,
		AdminLoc:   adminlocUc,
		AdminUser:  adminuserUc,
		Profile:    profileUc,
		Auth:       authUc,
	}

	// handlers
	handlers := &Handlers{
		auth:       callbacks.NewAuth(authUc, renderer),
		onboarding: callbacks.NewOnboarding(onboardingUc, renderer),
		booking:    callbacks.NewBooking(bookingUc, renderer),
		admin:      callbacks.NewAdmin(adminUc, renderer),
		adminLoc:   callbacks.NewAdminLoc(adminlocUc, renderer),
		adminUser:  callbacks.NewAdminUser(adminuserUc, renderer),
		profile:    callbacks.NewProfile(profileUc, renderer),
	}

	return &Bot{
		bot:      b,
		state:    state,
		renderer: renderer,
		users:    users,
		fsm:      fsm,
		handlers: handlers,
	}, nil
}

func (b *Bot) Start() {
	// router
	router := NewRouter(b, b.fsm, b.handlers)
	router.RegisterRoutes()

	log.Printf("Buskr Telegram Bot is starting as @%s...", b.bot.Me.Username)
	b.bot.Start()
}

func (b *Bot) GetTelebot() *telebot.Bot {
	return b.bot
}
