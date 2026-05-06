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

		// profile
		keys.StateProfileEditName: h.Profile.OnText,

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

	rep, err := h.AdminLoc.OnLocation(ctx, u, msg.Location.Lat, msg.Location.Lng)
	if err != nil {
		return err
	}
	if rep.IsEmpty() {
		return nil
	}

	return h.Renderer.Render(c, rep)
}
