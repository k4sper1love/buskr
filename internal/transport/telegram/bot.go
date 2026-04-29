package telegram

import (
	"context"
	"log"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/i18n"
	"github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers/callbacks"
	"github.com/k4sper1love/buskr/internal/transport/telegram/middleware"
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

const (
	adminChatID = -1234567890
	adminLang   = "ru"
)

type Bot struct {
	bot      *telebot.Bot
	state    *redis.StateStore
	renderer *render.Renderer
	users    *user.Service

	// usecases
	auth       *auth.Usecase
	booking    *bookingUsecase.Usecase
	adminUc    *admin.Usecase
	adminloc   *adminloc.Usecase
	adminuser  *adminuser.Usecase
	profile    *profile.Usecase
	onboarding *onboarding.Usecase
}

func NewBot(
	token string,
	users *user.Service,
	bookings *booking.Service,
	locs *location.Service,
	state *redis.StateStore,
) (*Bot, error) {

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	tr := i18n.NewTranslator()
	renderer := render.NewRenderer(tr, "ru")

	notifier := notifier.NewNotifier(b, tr, adminChatID, adminLang)

	authUc := auth.NewUsecase(state, users, 1*time.Hour)
	bookingUc := bookingUsecase.NewUsecase(state, locs, bookings, 1*time.Hour)
	adminUc := admin.NewUsecase(state, users, 1*time.Hour, 1)
	adminlocUc := adminloc.NewUsecase(state, locs, 1*time.Hour)
	adminuserUc := adminuser.NewUsecase(state, users, 1*time.Hour)
	profileUc := profile.NewUsecase(state, users, notifier, 1*time.Hour)
	onboardingUc := onboarding.NewUsecase(state, users, notifier, 1*time.Hour)

	return &Bot{
		bot:        b,
		state:      state,
		renderer:   renderer,
		auth:       authUc,
		booking:    bookingUc,
		adminUc:    adminUc,
		adminloc:   adminlocUc,
		adminuser:  adminuserUc,
		profile:    profileUc,
		onboarding: onboardingUc,
		users:      users,
	}, nil
}

func (b *Bot) GetTelebot() *telebot.Bot {
	return b.bot
}

func (b *Bot) Start() {
	// fsm handlers (text/video/location messages)
	fsm := &handlers.FSMHandlers{
		State:      b.state,
		Renderer:   b.renderer,
		Onboarding: b.onboarding,
		Booking:    b.booking,
		AdminLoc:   b.adminloc,
		AdminUser:  b.adminuser,
		Profile:    b.profile,
		Auth:       b.auth,
	}

	// callback handlers
	authCb := callbacks.NewAuthCallbacks(b.auth, b.renderer)
	bookingCb := callbacks.NewBookingCallbacks(b.booking, b.renderer)
	adminCb := callbacks.NewAdminCallbacks(b.adminUc, b.renderer)
	adminLocCb := callbacks.NewAdminLocCallbacks(b.adminloc, b.renderer)
	adminUserCb := callbacks.NewAdminUserCallbacks(b.adminuser, b.renderer)
	profileCb := callbacks.NewProfileCallbacks(b.profile, b.renderer)

	// middleware
	b.bot.Use(func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c.Set("ctx", ctx)
			return next(c)
		}
	})
	b.bot.Use(middleware.AuthMiddleware(b.users))

	// commands
	b.bot.Handle("/start", authCb.HandleStart)

	// fsm message handlers
	b.bot.Handle(telebot.OnText, fsm.HandleText)
	b.bot.Handle(telebot.OnVideo, fsm.HandleVideo)
	b.bot.Handle(telebot.OnDocument, fsm.HandleVideo)
	b.bot.Handle(telebot.OnLocation, fsm.HandleLocationMsg)

	// auth / onboarding
	b.bot.Handle("\f"+keys.BtnAuthApply, authCb.HandleApplyButton)

	// admin moderation (approve/reject apply + noise upgrade)
	b.bot.Handle("\f"+keys.BtnAdminAppAppr, adminCb.HandleAdminApproveApplication)
	b.bot.Handle("\f"+keys.BtnAdminAppRej, adminCb.HandleAdminRejectApplication)
	b.bot.Handle("\f"+keys.BtnAdminNoiseAppr, adminCb.HandleAdminApproveNoiseUpgrade)
	b.bot.Handle("\f"+keys.BtnAdminNoiseRej, adminCb.HandleAdminRejectNoiseUpgrade)

	// profile
	b.bot.Handle("\f"+keys.BtnProfileMain, profileCb.HandleProfile)
	b.bot.Handle("\f"+keys.BtnProfileEditName, profileCb.HandleEditName)
	b.bot.Handle("\f"+keys.BtnProfileNoiseUp, profileCb.HandleNoiseUpgrade)
	b.bot.Handle("\f"+keys.BtnProfileNoiseSel, profileCb.HandleNoiseUpgradeSelect)

	// admin panel
	b.bot.Handle("\f"+keys.BtnAdminPanel, adminCb.HandleAdminPanel)
	b.bot.Handle("\f"+keys.BtnAdminPanelSend, adminCb.HandleAdminPanelSend)

	// admin invites
	b.bot.Handle("\f"+keys.BtnAdminInvites, adminCb.HandleAdminInvites)
	b.bot.Handle("\f"+keys.BtnAdminInvGen, adminCb.HandleAdminGenerateInvite)

	// admin locations
	b.bot.Handle("\f"+keys.BtnAdminLocs, adminLocCb.HandleAdminLocs)
	b.bot.Handle("\f"+keys.BtnAdminLocDet, adminLocCb.HandleAdminLocDetails)
	b.bot.Handle("\f"+keys.BtnAdminLocTog, adminLocCb.HandleAdminLocToggle)
	b.bot.Handle("\f"+keys.BtnAdminLocAdd, adminLocCb.HandleAdminAddLoc)
	b.bot.Handle("\f"+keys.BtnAdminLocNoise, adminLocCb.HandleAdminOnNoiseSelected)
	b.bot.Handle("\f"+keys.BtnAdminLocCancel, adminLocCb.HandleAdminLocCancel)

	// admin users
	b.bot.Handle("\f"+keys.BtnAdminUsers, adminUserCb.HandleAdminSearch)
	b.bot.Handle("\f"+keys.BtnAdminUserBan, adminUserCb.HandleAdminBan)
	b.bot.Handle("\f"+keys.BtnAdminUserUnban, adminUserCb.HandleAdminUnban)

	// bookings list / details / cancel
	b.bot.Handle("\f"+keys.BtnBookList, bookingCb.HandleBookingList)
	b.bot.Handle("\f"+keys.BtnBookDetails, bookingCb.HandleBookingDetails)
	b.bot.Handle("\f"+keys.BtnBookCancelConf, bookingCb.HandleBookingCancelConfirm)

	// booking fsm flow
	b.bot.Handle("\f"+keys.BtnBookStart, bookingCb.HandleBook)
	b.bot.Handle("\f"+keys.BtnBookDateSel, bookingCb.HandleBookDateSelected)
	b.bot.Handle("\f"+keys.BtnBookLocSel, bookingCb.HandleBookLocationSelected)
	b.bot.Handle("\f"+keys.BtnBookSlotSel, bookingCb.HandleBookSlotSelected)
	b.bot.Handle("\f"+keys.BtnBookDurSel, bookingCb.HandleBookDurationSelected)

	// fsm backward navigation
	b.bot.Handle("\f"+keys.BtnBookBackToLocs, bookingCb.HandleBookingBackToLocs)
	b.bot.Handle("\f"+keys.BtnBookBackToSlots, bookingCb.HandleBookingBackToSlots)

	// hot slots
	b.bot.Handle("\f"+keys.BtnBookGrabHot, bookingCb.HandleBookingGrabHotSlot)

	// check-in
	b.bot.Handle("\f"+keys.BtnBookCheckin, bookingCb.HandleBookingCheckIn)

	// start over — callback, so we use HandleCallbackStart
	b.bot.Handle("\f"+keys.BtnCommonMenu, authCb.HandleCallbackStart)
	b.bot.Handle("\f"+keys.BtnCommonMenuSend, authCb.HandleCallbackStartSend)

	// cancel booking flow (not an existing booking, but the creation process)
	b.bot.Handle("\f"+keys.BtnBookCancel, func(c telebot.Context) error {
		ctx := c.Get("ctx").(context.Context)
		u, err := ctxkey.GetUser(c)
		if err != nil {
			return err
		}
		rep, err := bookingCb.Uc().CancelFlow(ctx, u)
		if err != nil {
			return err
		}
		c.Respond(&telebot.CallbackResponse{})
		return b.renderer.Render(c, rep)
	})

	log.Printf("Buskr Telegram Bot is starting as @%s...", b.bot.Me.Username)
	b.bot.Start()
}
