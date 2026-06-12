package callbacks

import (
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/onboarding"
	"gopkg.in/telebot.v3"
)

type Onboarding struct {
	uc       *onboarding.Usecase
	renderer *render.Renderer
}

func NewOnboarding(uc *onboarding.Usecase, renderer *render.Renderer) *Onboarding {
	return &Onboarding{
		uc:       uc,
		renderer: renderer,
	}
}

func (h *Onboarding) HandleNoiseSelected(c telebot.Context) error {
	cbCtx, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	noise := args[0]

	rep, err := h.uc.NoiseSelected(cbCtx.ctx, cbCtx.user, user.NoiseLevel(noise))
	if err != nil {
		return nil
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *Onboarding) HandleCancelFlow(c telebot.Context) error {
	cbCtx, err := extractCtx(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.CancelFlow(cbCtx.ctx, cbCtx.user)
	if err != nil {
		return nil
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *Onboarding) HandleSkipStep(c telebot.Context) error {
	cbCtx, err := extractCtx(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.SkipMedia(cbCtx.ctx, cbCtx.user)
	if err != nil {
		return nil
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}
