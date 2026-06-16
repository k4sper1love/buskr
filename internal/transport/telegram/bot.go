package telegram

import (
	"errors"
	"log/slog"
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
	"github.com/k4sper1love/buskr/internal/usecase/keys"
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
	group      *handlers.GroupHandler
}

func NewBot(
	cfg *config.TelegramConfig,
	tr *i18n.Translator,
	state *redis.StateStore,
	users *user.Service,
	bookings *booking.Service,
	locs *location.Service,
	maxAdvanceDays int,
	webAppURL string,
	minHour int,
	maxHourWeekday int,
	maxHourWeekend int,
	maxDurationHours int,
	adminMaxDurationHours int,
) (*Bot, error) {
	// create bot instance
	pref := telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		OnError: func(err error, c telebot.Context) {
			// Handle access denied gracefully via callback alert
			if c != nil && c.Callback() != nil && errors.Is(err, user.ErrAccessDenied) {
				lang := cfg.DefaultLang
				if c.Sender() != nil && c.Sender().LanguageCode != "" {
					lang = c.Sender().LanguageCode
				}
				text := tr.T(lang, keys.TextCommonErrAccessDenied, nil)
				_ = c.Respond(&telebot.CallbackResponse{
					Text:      text,
					ShowAlert: true,
				})
				return
			}
			args := []any{"err", err}
			if c != nil && c.Sender() != nil {
				args = append(args, "telegram_id", c.Sender().ID)
			}
			slog.Error("unhandled handler error", args...)
		},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	// renderer for bot responses
	renderer := render.NewRenderer(tr, cfg.DefaultLang)

	// admin notifier
	notifier := notifier.NewNotifier(b, tr, cfg.AdminChatID, cfg.AdminLang, cfg.AdminThreadApplications, cfg.AdminThreadUpgrades, cfg.AdminThreadSuggestions)

	// usecases
	authUc := auth.NewUsecase(state, users, notifier, 1*time.Hour, b.Me.FirstName)
	bookingUc := bookingUsecase.NewUsecase(
		state,
		locs,
		bookings,
		users,
		notifier,
		1*time.Hour,
		maxAdvanceDays,
		webAppURL,
		b.Me.Username,
		minHour,
		maxHourWeekday,
		maxHourWeekend,
		maxDurationHours,
		adminMaxDurationHours,
	)
	adminUc := admin.NewUsecase(state, users, locs, 1*time.Hour, cfg.InviteTTLHours)
	adminlocUc := adminloc.NewUsecase(state, locs, bookings, users, 1*time.Hour, maxAdvanceDays, webAppURL, b.Me.Username)
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
		auth:       callbacks.NewAuth(authUc, bookingUc, adminlocUc, renderer, cfg.SupportContact),
		onboarding: callbacks.NewOnboarding(onboardingUc, authUc, renderer),
		booking:    callbacks.NewBooking(bookingUc, authUc, renderer),
		admin:      callbacks.NewAdmin(adminUc, renderer),
		adminLoc:   callbacks.NewAdminLoc(adminlocUc, renderer),
		adminUser:  callbacks.NewAdminUser(adminuserUc, renderer),
		profile:    callbacks.NewProfile(profileUc, renderer),
		group:      handlers.NewGroupHandler(cfg, renderer),
	}

	return &Bot{
		bot:      b,
		cfg:      cfg,
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

	slog.Info("buskr telegram bot starting", "username", b.bot.Me.Username)
	b.bot.Start()
}

func (b *Bot) GetTelebot() *telebot.Bot {
	return b.bot
}
