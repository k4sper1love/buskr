package callbacks

import (
	"context"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/auth"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type Auth struct {
	uc       *auth.Usecase
	renderer *render.Renderer
}

func NewAuth(uc *auth.Usecase, renderer *render.Renderer) *Auth {
	return &Auth{
		uc:       uc,
		renderer: renderer,
	}
}

// HandleStart handles the /start command — no c.Respond() since it's not a callback.
func (h *Auth) HandleStart(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Start(ctx, u, c.Message().Payload)
	if err != nil {
		return err
	}

	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

// HandleCallbackStart handles btn_start_over callback — edits the current message.
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

	// Always edit the current message — this is a callback, not a fresh /start
	rep.Kind = response.KindEdit

	return h.renderer.Render(c, rep)
}

// HandleCallbackStartSend handles btn_common_menu_send callback — sends a new message instead of editing.
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

// HandleApplyButton handles btn_apply callback.
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
