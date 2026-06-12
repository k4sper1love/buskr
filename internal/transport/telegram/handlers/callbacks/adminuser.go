package callbacks

import (
	"context"
	"strconv"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/adminuser"
	"gopkg.in/telebot.v3"
)

type AdminUser struct {
	uc       *adminuser.Usecase
	renderer *render.Renderer
}

func NewAdminUser(uc *adminuser.Usecase, renderer *render.Renderer) *AdminUser {
	return &AdminUser{
		uc:       uc,
		renderer: renderer,
	}
}

func (h *AdminUser) HandleAdminUsersList(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	page := 0
	args := c.Args()
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil {
			page = p
		}
	}

	rep, err := h.uc.UsersList(ctx, u, page)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUserDetail(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error: missing user ID", ShowAlert: true})
	}
	targetID := args[0]

	rep, err := h.uc.UserDetail(ctx, u, targetID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUserSearchPrompt(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.SearchPrompt(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUserSearchPage(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	page := 0
	args := c.Args()
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil {
			page = p
		}
	}

	rep, err := h.uc.SearchPage(ctx, u, page)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminBan(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	targetID := args[0]

	rep, err := h.uc.Ban(ctx, u, targetID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUnban(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	targetID := args[0]

	rep, err := h.uc.Unban(ctx, u, targetID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUserNoiseMenu(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error: missing user ID", ShowAlert: true})
	}
	targetID := args[0]

	rep, err := h.uc.UserNoiseMenu(ctx, u, targetID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminUser) HandleAdminUserNoiseSet(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error: missing arguments", ShowAlert: true})
	}
	targetID := args[0]
	noise := args[1]

	rep, err := h.uc.SetUserNoise(ctx, u, targetID, noise)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}
