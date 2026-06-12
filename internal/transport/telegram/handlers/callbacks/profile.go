package callbacks

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/profile"
	"gopkg.in/telebot.v3"
)

type Profile struct {
	uc       *profile.Usecase
	renderer *render.Renderer
}

func NewProfile(uc *profile.Usecase, renderer *render.Renderer) *Profile {
	return &Profile{
		uc:       uc,
		renderer: renderer,
	}
}

func (h *Profile) HandleProfile(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Profile(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *Profile) HandleNoiseUpgrade(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.NoiseUpgrade(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *Profile) HandleNoiseUpgradeSelect(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка", ShowAlert: true})
	}
	requestedNoise := args[0]

	rep, err := h.uc.NoiseUpgradeSelected(ctx, u, user.NoiseLevel(requestedNoise))
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	_ = h.uc.SaveProfileMessageID(ctx, u.TelegramID, c.Message().ID)
	return nil
}

func (h *Profile) HandleEditName(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.EditNameStart(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	_ = h.uc.SaveProfileMessageID(ctx, u.TelegramID, c.Message().ID)
	return nil
}

func (h *Profile) HandleNoiseUpgradeSkipReason(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.NoiseUpgradeSubmit(ctx, u, "")
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	_ = h.uc.ClearProfileMessageID(ctx, u.TelegramID)

	return h.renderer.Render(c, rep)
}
