package callbacks

import (
	"context"
	"strconv"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/adminloc"
	"gopkg.in/telebot.v3"
)

type AdminLoc struct {
	uc       *adminloc.Usecase
	renderer *render.Renderer
}

func NewAdminLoc(uc *adminloc.Usecase, renderer *render.Renderer) *AdminLoc {
	return &AdminLoc{
		uc:       uc,
		renderer: renderer,
	}
}

func (h *AdminLoc) HandleAdminAddLoc(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.AddStart(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminOnNoiseSelected(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	noise := args[0]

	rep, err := h.uc.OnNoiseSelected(ctx, u, noise)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminOnLocation(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	latStr := args[0]
	lonStr := args[1]

	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		return err
	}
	lon, err := strconv.ParseFloat(lonStr, 32)
	if err != nil {
		return err
	}

	rep, err := h.uc.OnLocation(ctx, u, float32(lat), float32(lon))
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocs(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.List(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocDetails(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	locID := args[0]

	rep, err := h.uc.Details(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocToggle(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	locID := args[0]

	rep, err := h.uc.ToggleStatus(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocCancel(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.CancelFlow(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}
