package telegram

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/handlers"
	"github.com/k4sper1love/buskr/internal/transport/telegram/middleware"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
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
	r.b.bot.Handle("/help", r.handlers.auth.HandleHelp)

	// fsm message handlers
	r.b.bot.Handle(telebot.OnText, r.fsm.HandleText)
	r.b.bot.Handle(telebot.OnVideo, r.fsm.HandleVideo)
	r.b.bot.Handle(telebot.OnDocument, r.fsm.HandleVideo)
	r.b.bot.Handle(telebot.OnLocation, r.fsm.HandleLocationMsg)

	// auth / onboarding
	r.b.bot.Handle("\f"+keys.BtnAuthApply, r.handlers.auth.HandleApplyButton)

	// onboarding
	r.b.bot.Handle("\f"+keys.BtnOnboardNoiseSel, r.handlers.onboarding.HandleNoiseSelected)
	r.b.bot.Handle("\f"+keys.BtnOnboardCancel, r.handlers.onboarding.HandleCancelFlow)
	r.b.bot.Handle("\f"+keys.BtnOnboardSkip, r.handlers.onboarding.HandleSkipStep)

	// admin moderation (approve/reject apply + noise upgrade)
	r.b.bot.Handle("\f"+keys.BtnAdminAppAppr, r.handlers.admin.HandleAdminApproveApplication)
	r.b.bot.Handle("\f"+keys.BtnAdminAppRej, r.handlers.admin.HandleAdminRejectApplication)
	r.b.bot.Handle("\f"+keys.BtnAdminNoiseAppr, r.handlers.admin.HandleAdminApproveNoiseUpgrade)
	r.b.bot.Handle("\f"+keys.BtnAdminNoiseRej, r.handlers.admin.HandleAdminRejectNoiseUpgrade)
	r.b.bot.Handle("\f"+keys.BtnAdminLocSuggestAppr, r.handlers.admin.HandleAdminLocSuggestApprove)
	r.b.bot.Handle("\f"+keys.BtnAdminLocSuggestRej, r.handlers.admin.HandleAdminLocSuggestReject)

	// profile
	r.b.bot.Handle("\f"+keys.BtnProfileMain, r.handlers.profile.HandleProfile)
	r.b.bot.Handle("\f"+keys.BtnProfileEditName, r.handlers.profile.HandleEditName)
	r.b.bot.Handle("\f"+keys.BtnProfileNoiseUp, r.handlers.profile.HandleNoiseUpgrade)
	r.b.bot.Handle("\f"+keys.BtnProfileNoiseSel, r.handlers.profile.HandleNoiseUpgradeSelect)
	r.b.bot.Handle("\f"+keys.BtnProfileNoiseSkipReason, r.handlers.profile.HandleNoiseUpgradeSkipReason)

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
	r.b.bot.Handle("\f"+keys.BtnAdminLocTogVeteran, r.handlers.adminLoc.HandleAdminLocToggleVeteran)
	r.b.bot.Handle("\f"+keys.BtnAdminLocAdd, r.handlers.adminLoc.HandleAdminAddLoc)
	r.b.bot.Handle("\f"+keys.BtnAdminLocNoise, r.handlers.adminLoc.HandleAdminOnNoiseSelected)
	r.b.bot.Handle("\f"+keys.BtnAdminLocCancel, r.handlers.adminLoc.HandleAdminLocCancel)
	r.b.bot.Handle("\f"+keys.BtnAdminLocConfirm, r.handlers.adminLoc.HandleAdminLocConfirm)

	r.b.bot.Handle("\f"+keys.BtnAdminLocEdit, r.handlers.adminLoc.HandleAdminLocEdit)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditName, r.handlers.adminLoc.HandleAdminLocEditName)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditDesc, r.handlers.adminLoc.HandleAdminLocEditDesc)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditNoise, r.handlers.adminLoc.HandleAdminLocEditNoise)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditNoiseMenu, r.handlers.adminLoc.HandleAdminLocEditNoiseMenu)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditGeo, r.handlers.adminLoc.HandleAdminLocEditGeo)
	r.b.bot.Handle("\f"+keys.BtnAdminLocEditCancel, r.handlers.adminLoc.HandleAdminLocEditCancel)

	r.b.bot.Handle("\f"+keys.BtnAdminLocSchedule, r.handlers.adminLoc.HandleAdminLocSchedule)
	r.b.bot.Handle("\f"+keys.BtnAdminLocDel, r.handlers.adminLoc.HandleAdminLocDelete)
	r.b.bot.Handle("\f"+keys.BtnAdminLocDelConf, r.handlers.adminLoc.HandleAdminLocDeleteConfirm)

	// admin users
	r.b.bot.Handle("\f"+keys.BtnAdminUsers, r.handlers.adminUser.HandleAdminUsersList)
	r.b.bot.Handle("\f"+keys.BtnAdminUsersPage, r.handlers.adminUser.HandleAdminUsersList)
	r.b.bot.Handle("\f"+keys.BtnAdminUserSearchPage, r.handlers.adminUser.HandleAdminUserSearchPage)
	r.b.bot.Handle("\f"+keys.BtnAdminUserDetail, r.handlers.adminUser.HandleAdminUserDetail)
	r.b.bot.Handle("\f"+keys.BtnAdminUserSearchPrompt, r.handlers.adminUser.HandleAdminUserSearchPrompt)
	r.b.bot.Handle("\f"+keys.BtnAdminUserNoiseMenu, r.handlers.adminUser.HandleAdminUserNoiseMenu)
	r.b.bot.Handle("\f"+keys.BtnAdminUserNoiseSet, r.handlers.adminUser.HandleAdminUserNoiseSet)
	r.b.bot.Handle("\f"+keys.BtnAdminUserBan, r.handlers.adminUser.HandleAdminBan)
	r.b.bot.Handle("\f"+keys.BtnAdminUserUnban, r.handlers.adminUser.HandleAdminUnban)
	r.b.bot.Handle("\f"+keys.BtnAdminUserPromote, r.handlers.adminUser.HandleAdminPromote)
	r.b.bot.Handle("\f"+keys.BtnAdminUserDemote, r.handlers.adminUser.HandleAdminDemote)
	r.b.bot.Handle("\fnoop", func(c telebot.Context) error {
		return c.Respond(&telebot.CallbackResponse{})
	})

	// bookings list / details / cancel
	r.b.bot.Handle("\f"+keys.BtnBookList, r.handlers.booking.HandleBookingList)
	r.b.bot.Handle("\f"+keys.BtnBookDetails, r.handlers.booking.HandleBookingDetails)
	r.b.bot.Handle("\f"+keys.BtnBookCancelConf, r.handlers.booking.HandleBookingCancelConfirm)

	// booking schedule
	r.b.bot.Handle("\f"+keys.BtnBookSchedule, r.handlers.booking.HandleBookScheduleStart)
	r.b.bot.Handle("\f"+keys.BtnBookScheduleLocSel, r.handlers.booking.HandleBookScheduleLocSel)
	r.b.bot.Handle("\f"+keys.BtnBookScheduleDaySel, r.handlers.booking.HandleBookScheduleDaySel)

	// booking fsm flow
	r.b.bot.Handle("\f"+keys.BtnBookStart, r.handlers.booking.HandleBook)
	r.b.bot.Handle("\f"+keys.BtnBookDateSel, r.handlers.booking.HandleBookDateSelected)
	r.b.bot.Handle("\f"+keys.BtnBookLocSel, r.handlers.booking.HandleBookLocationSelected)
	r.b.bot.Handle("\f"+keys.BtnBookSlotSel, r.handlers.booking.HandleBookSlotSelected)
	r.b.bot.Handle("\f"+keys.BtnBookDurSel, r.handlers.booking.HandleBookDurationSelected)

	// suggest location flow
	r.b.bot.Handle("\f"+keys.BtnSuggestLocStart, r.handlers.booking.HandleSuggestLocStart)
	r.b.bot.Handle("\f"+keys.BtnSuggestLocNoise, r.handlers.booking.HandleSuggestLocNoise)
	r.b.bot.Handle("\f"+keys.BtnSuggestLocConfirm, r.handlers.booking.HandleSuggestLocConfirm)
	r.b.bot.Handle("\f"+keys.BtnSuggestLocCancel, r.handlers.booking.HandleSuggestLocCancel)

	// fsm backward navigation
	r.b.bot.Handle("\f"+keys.BtnBookBackToLocs, r.handlers.booking.HandleBookingBackToLocs)
	r.b.bot.Handle("\f"+keys.BtnBookBackToSlots, r.handlers.booking.HandleBookingBackToSlots)

	// hot slots
	r.b.bot.Handle("\f"+keys.BtnBookGrabHot, r.handlers.booking.HandleBookingGrabHotSlot)
	r.b.bot.Handle("\f"+keys.BtnBookDismissHot, func(c telebot.Context) error {
		_ = c.Respond(&telebot.CallbackResponse{})
		return c.Delete()
	})

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

		cancelRep, _ := r.handlers.booking.Uc().CancelFlow(ctx, u)

		lang := ""
		if s := c.Sender(); s != nil {
			lang = s.LanguageCode
		}
		cancelText := r.b.renderer.Translate(lang, cancelRep.Text)
		_ = c.Respond(&telebot.CallbackResponse{Text: cancelText})

		authRep, err := r.handlers.auth.Uc().Start(ctx, u, "")
		if err != nil {
			return err
		}
		if authRep.IsEmpty() {
			return nil
		}
		authRep.Kind = response.KindEdit
		return r.b.renderer.Render(c, authRep)
	})

	r.b.bot.Handle(telebot.OnAddedToGroup, r.handlers.group.HandleAddedToGroup)
	r.b.bot.Handle(telebot.OnUserJoined, r.handlers.group.HandleUserJoined)
}
