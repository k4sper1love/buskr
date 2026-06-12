package callbacks

import (
	"context"
	"strings"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/adminloc"
	"github.com/k4sper1love/buskr/internal/usecase/auth"
	"github.com/k4sper1love/buskr/internal/usecase/booking"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type Auth struct {
	uc         *auth.Usecase
	bookingUc  *booking.Usecase
	adminLocUc *adminloc.Usecase
	renderer   *render.Renderer
}

func NewAuth(uc *auth.Usecase, bookingUc *booking.Usecase, adminLocUc *adminloc.Usecase, renderer *render.Renderer) *Auth {
	return &Auth{
		uc:         uc,
		bookingUc:  bookingUc,
		adminLocUc: adminLocUc,
		renderer:   renderer,
	}
}

func (h *Auth) Uc() *auth.Usecase {
	return h.uc
}

func (h *Auth) HandleStart(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	payload := c.Message().Payload
	if strings.HasPrefix(payload, "admloc_") || strings.HasPrefix(payload, "loc_") {
		_ = c.Delete()

	}

	if strings.HasPrefix(payload, "admloc_") {
		locID := strings.TrimPrefix(payload, "admloc_")
		if u.Role == user.RoleAdmin {
			rep, err := h.adminLocUc.Details(ctx, u, locID)
			if err != nil {
				return err
			}
			if rep.IsEmpty() {
				return nil
			}
			rep.Kind = response.KindSend
			return h.renderer.Render(c, rep)
		}
	}

	if strings.HasPrefix(payload, "loc_") {
		locID := strings.TrimPrefix(payload, "loc_")
		rep, err := h.bookingUc.ProcessMapSelection(ctx, u, locID)
		if err != nil {
			return err
		}
		if rep.IsEmpty() {
			return nil
		}
		rep.Kind = response.KindSend
		return h.renderer.Render(c, rep)
	}

	rep, err := h.uc.Start(ctx, u, payload)
	if err != nil {
		return err
	}

	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *Auth) HandleCallbackStart(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Start(ctx, u, "")
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	rep.Kind = response.KindEdit

	return h.renderer.Render(c, rep)
}

func (h *Auth) HandleCallbackStartSend(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Start(ctx, u, "")
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	// Remove inline keyboard from the older message to prevent stale presses
	_ = c.Edit(c.Message().Text, c.Message().Entities, &telebot.ReplyMarkup{})

	// Force Send rather than Edit
	rep.Kind = response.KindSend

	return h.renderer.Render(c, rep)
}

func (h *Auth) HandleApplyButton(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.ApplyButton(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}
