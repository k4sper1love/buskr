package handlers

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/infrastructure/redis"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/adminloc"
	"github.com/k4sper1love/buskr/internal/usecase/adminuser"
	"github.com/k4sper1love/buskr/internal/usecase/auth"
	"github.com/k4sper1love/buskr/internal/usecase/booking"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/onboarding"
	"github.com/k4sper1love/buskr/internal/usecase/profile"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type FSM struct {
	State      *redis.StateStore
	Renderer   *render.Renderer
	Onboarding *onboarding.Usecase
	Booking    *booking.Usecase
	AdminLoc   *adminloc.Usecase
	AdminUser  *adminuser.Usecase
	Profile    *profile.Usecase
	Auth       *auth.Usecase
}

type textStep func(ctx context.Context, u *user.User, text string) (response.Reply, error)

func (h *FSM) HandleText(c telebot.Context) error {
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}
	ctx := c.Get("ctx").(context.Context)

	state, _ := h.State.GetState(ctx, c.Sender().ID)

	steps := map[string]textStep{
		// auth
		keys.StateAuthInvitedName: h.Auth.OnText,

		// onboarding
		keys.StateOnboardName:  h.Onboarding.OnText,
		keys.StateOnboardMedia: h.Onboarding.OnText,

		// admin location
		keys.StateAdminLocName: h.AdminLoc.OnText,
		keys.StateAdminLocDesc: h.AdminLoc.OnText,

		keys.StateAdminLocEditName: h.AdminLoc.OnText,
		keys.StateAdminLocEditDesc: h.AdminLoc.OnText,

		// profile
		keys.StateProfileEditName:    h.Profile.OnText,
		keys.StateProfileNoiseReason: h.Profile.OnText,

		// admin user
		keys.StateAdminUserSearch: h.AdminUser.OnText,
	}

	step, ok := steps[state]
	if !ok {
		// no active FSM state
		return nil
	}

	rep, err := step(ctx, u, c.Text())
	if err != nil {
		return err
	}
	if rep.IsEmpty() {
		return nil
	}

	if state == keys.StateProfileEditName {
		_ = c.Delete()

		if msgID, err := h.Profile.GetProfileMessageID(ctx, u.TelegramID); err == nil && msgID != 0 {
			_ = c.Bot().Delete(&telebot.Message{
				ID:   msgID,
				Chat: &telebot.Chat{ID: u.TelegramID},
			})
			_ = h.Profile.ClearProfileMessageID(ctx, u.TelegramID)
		}

		success := response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextProfileEditNameMsgSuccess},
		}
		_ = h.Renderer.Render(c, success)
	}

	if state == keys.StateProfileNoiseReason {
		_ = c.Delete()

		if msgID, err := h.Profile.GetProfileMessageID(ctx, u.TelegramID); err == nil && msgID != 0 {
			_ = c.Bot().Delete(&telebot.Message{
				ID:   msgID,
				Chat: &telebot.Chat{ID: u.TelegramID},
			})
			_ = h.Profile.ClearProfileMessageID(ctx, u.TelegramID)
		}
		rep.Kind = response.KindSend
	}

	if state == keys.StateAdminLocName || state == keys.StateAdminLocDesc {
		_ = c.Delete()

		msgID, err := h.AdminLoc.GetAdminLocMessageID(ctx, u.TelegramID)
		if err == nil && msgID != 0 {
			lang := ""
			if s := c.Sender(); s != nil {
				lang = s.LanguageCode
			}
			return h.Renderer.RenderToMessage(c.Bot(), u.TelegramID, msgID, lang, rep)
		}
	}

	if state == keys.StateAdminLocEditName || state == keys.StateAdminLocEditDesc {
		success := response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsEditMsgSuccess},
		}
		_ = h.Renderer.Render(c, success)
	}

	return h.Renderer.Render(c, rep)
}

func (h *FSM) HandleVideo(c telebot.Context) error {
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}
	ctx := c.Get("ctx").(context.Context)

	state, _ := h.State.GetState(ctx, c.Sender().ID)

	if state == keys.StateOnboardMedia {
		rep, err := h.Onboarding.OnVideo(ctx, u, c.Message().Video.FileID)
		if err != nil {
			return err
		}
		if rep.IsEmpty() {
			return nil
		}
		return h.Renderer.Render(c, rep)
	}

	return nil
}

func (h *FSM) HandleLocationMsg(c telebot.Context) error {
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}
	ctx := c.Get("ctx").(context.Context)

	msg := c.Message()
	if msg == nil || msg.Location == nil {
		return nil
	}

	state, _ := h.State.GetState(ctx, c.Sender().ID)

	rep, err := h.AdminLoc.OnLocation(ctx, u, msg.Location.Lat, msg.Location.Lng)
	if err != nil {
		return err
	}
	if rep.IsEmpty() {
		return nil
	}

	if state == keys.StateAdminLocGeo {
		_ = c.Delete()

		if msgID, err := h.AdminLoc.GetAdminLocMessageID(ctx, u.TelegramID); err == nil && msgID != 0 {
			_ = c.Bot().Delete(&telebot.Message{
				ID:   msgID,
				Chat: &telebot.Chat{ID: u.TelegramID},
			})
			_ = h.AdminLoc.ClearAdminLocMessageID(ctx, u.TelegramID)
		}

		rep.Kind = response.KindSend
	}

	if state == keys.StateAdminLocEditGeo {
		success := response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsEditMsgSuccess},
		}
		_ = h.Renderer.Render(c, success)
	}

	return h.Renderer.Render(c, rep)
}
