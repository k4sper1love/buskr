package telegram

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers"
	"github.com/k4sper1love/buskr/internal/transport/telegram/middleware"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"gopkg.in/telebot.v3"
)

type Router struct {
	b        *Bot
	fsm      *handlers.FSM
	handlers *Handlers
}

func NewRouter(b *Bot, fsm *handlers.FSM, handlers *Handlers) *Router {
	return &Router{
		b:        b,
		fsm:      fsm,
		handlers: handlers,
	}
}

func (r *Router) RegisterRoutes() {
	r.b.bot.Use(func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c.Set("ctx", ctx)
			return next(c)
		}
	})
	r.b.bot.Use(middleware.AuthMiddleware(r.b.users))

	// commands
	r.b.bot.Handle("/start", r.handlers.auth.HandleStart)

	// fsm message handlers
	r.b.bot.Handle(telebot.OnText, r.fsm.HandleText)
	r.b.bot.Handle(telebot.OnVideo, r.fsm.HandleVideo)
	r.b.bot.Handle(telebot.OnDocument, r.fsm.HandleVideo)
	r.b.bot.Handle(telebot.OnLocation, r.fsm.HandleLocationMsg)

	// auth / onboarding
	r.b.bot.Handle("\f"+keys.BtnAuthApply, r.handlers.auth.HandleApplyButton)

	// admin moderation (approve/reject apply + noise upgrade)
	r.b.bot.Handle("\f"+keys.BtnAdminAppAppr, r.handlers.admin.HandleAdminApproveApplication)
	r.b.bot.Handle("\f"+keys.BtnAdminAppRej, r.handlers.admin.HandleAdminRejectApplication)
	r.b.bot.Handle("\f"+keys.BtnAdminNoiseAppr, r.handlers.admin.HandleAdminApproveNoiseUpgrade)
	r.b.bot.Handle("\f"+keys.BtnAdminNoiseRej, r.handlers.admin.HandleAdminRejectNoiseUpgrade)

	// profile
	r.b.bot.Handle("\f"+keys.BtnProfileMain, r.handlers.profile.HandleProfile)
	r.b.bot.Handle("\f"+keys.BtnProfileEditName, r.handlers.profile.HandleEditName)
	r.b.bot.Handle("\f"+keys.BtnProfileNoiseUp, r.handlers.profile.HandleNoiseUpgrade)
	r.b.bot.Handle("\f"+keys.BtnProfileNoiseSel, r.handlers.profile.HandleNoiseUpgradeSelect)

	// admin panel
	r.b.bot.Handle("\f"+keys.BtnAdminPanel, r.handlers.admin.HandleAdminPanel)
	r.b.bot.Handle("\f"+keys.BtnAdminPanelSend, r.handlers.admin.HandleAdminPanelSend)

	// admin invites
	r.b.bot.Handle("\f"+keys.BtnAdminInvites, r.handlers.admin.HandleAdminInvites)
	r.b.bot.Handle("\f"+keys.BtnAdminInvGen, r.handlers.admin.HandleAdminGenerateInvite)

	// admin locations
	r.b.bot.Handle("\f"+keys.BtnAdminLocs, r.handlers.adminLoc.HandleAdminLocs)
	r.b.bot.Handle("\f"+keys.BtnAdminLocDet, r.handlers.adminLoc.HandleAdminLocDetails)
	r.b.bot.Handle("\f"+keys.BtnAdminLocTog, r.handlers.adminLoc.HandleAdminLocToggle)
	r.b.bot.Handle("\f"+keys.BtnAdminLocAdd, r.handlers.adminLoc.HandleAdminAddLoc)
	r.b.bot.Handle("\f"+keys.BtnAdminLocNoise, r.handlers.adminLoc.HandleAdminOnNoiseSelected)
	r.b.bot.Handle("\f"+keys.BtnAdminLocCancel, r.handlers.adminLoc.HandleAdminLocCancel)

	// admin users
	r.b.bot.Handle("\f"+keys.BtnAdminUsers, r.handlers.adminUser.HandleAdminSearch)
	r.b.bot.Handle("\f"+keys.BtnAdminUserBan, r.handlers.adminUser.HandleAdminBan)
	r.b.bot.Handle("\f"+keys.BtnAdminUserUnban, r.handlers.adminUser.HandleAdminUnban)

	// bookings list / details / cancel
	r.b.bot.Handle("\f"+keys.BtnBookList, r.handlers.booking.HandleBookingList)
	r.b.bot.Handle("\f"+keys.BtnBookDetails, r.handlers.booking.HandleBookingDetails)
	r.b.bot.Handle("\f"+keys.BtnBookCancelConf, r.handlers.booking.HandleBookingCancelConfirm)

	// booking fsm flow
	r.b.bot.Handle("\f"+keys.BtnBookStart, r.handlers.booking.HandleBook)
	r.b.bot.Handle("\f"+keys.BtnBookDateSel, r.handlers.booking.HandleBookDateSelected)
	r.b.bot.Handle("\f"+keys.BtnBookLocSel, r.handlers.booking.HandleBookLocationSelected)
	r.b.bot.Handle("\f"+keys.BtnBookSlotSel, r.handlers.booking.HandleBookSlotSelected)
	r.b.bot.Handle("\f"+keys.BtnBookDurSel, r.handlers.booking.HandleBookDurationSelected)

	// fsm backward navigation
	r.b.bot.Handle("\f"+keys.BtnBookBackToLocs, r.handlers.booking.HandleBookingBackToLocs)
	r.b.bot.Handle("\f"+keys.BtnBookBackToSlots, r.handlers.booking.HandleBookingBackToSlots)

	// hot slots
	r.b.bot.Handle("\f"+keys.BtnBookGrabHot, r.handlers.booking.HandleBookingGrabHotSlot)

	// check-in
	r.b.bot.Handle("\f"+keys.BtnBookCheckin, r.handlers.booking.HandleBookingCheckIn)

	// start over — callback, so we use HandleCallbackStart
	r.b.bot.Handle("\f"+keys.BtnCommonMenu, r.handlers.auth.HandleCallbackStart)
	r.b.bot.Handle("\f"+keys.BtnCommonMenuSend, r.handlers.auth.HandleCallbackStartSend)

	// cancel booking flow (not an existing booking, but the creation process)
	r.b.bot.Handle("\f"+keys.BtnBookCancel, func(c telebot.Context) error {
		ctx := c.Get("ctx").(context.Context)
		u, err := ctxkey.GetUser(c)
		if err != nil {
			return err
		}
		rep, err := r.handlers.booking.Uc().CancelFlow(ctx, u)
		if err != nil {
			return err
		}
		c.Respond(&telebot.CallbackResponse{})
		return r.b.renderer.Render(c, rep)
	})
}
