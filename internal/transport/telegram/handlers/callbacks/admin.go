package callbacks

import (
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/admin"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type AdminCallbacks struct {
	uc       *admin.Usecase
	renderer *render.Renderer
}

func NewAdminCallbacks(uc *admin.Usecase, renderer *render.Renderer) *AdminCallbacks {
	return &AdminCallbacks{
		uc:       uc,
		renderer: renderer,
	}
}

func (h *AdminCallbacks) HandleAdminPanel(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.AdminPanel(cc.ctx, cc.user)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminCallbacks) HandleAdminPanelSend(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.AdminPanel(cc.ctx, cc.user)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	// Remove the inline keyboard from the older message to prevent stale presses
	_ = c.Edit(c.Message().Text, c.Message().Entities, &telebot.ReplyMarkup{})

	rep.Kind = response.KindSend
	return h.renderer.Render(c, rep)
}

func (h *AdminCallbacks) HandleAdminInvites(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Invites(cc.ctx, cc.user)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminCallbacks) HandleAdminGenerateInvite(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: не передан уровень шума", ShowAlert: true})
	}
	noise := args[0]
	botUsername := c.Bot().Me.Username

	rep, err := h.uc.GenerateInvite(cc.ctx, cc.user, botUsername, noise)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminCallbacks) HandleAdminApproveApplication(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: неверные аргументы", ShowAlert: true})
	}
	targetUserID := args[0]
	category := user.NoiseLevel(args[1])

	rep, err := h.uc.ApproveApplication(cc.ctx, cc.user, targetUserID, category)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при одобрении", ShowAlert: true})
	}

	suffix := h.renderer.Translate("", rep.AdminEditSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix
	_, _ = c.Bot().Edit(c.Message(), updatedText, telebot.ModeMarkdown)

	if rep.NotifyUser != nil {
		sendReplyToUser(c.Bot(), h.renderer, rep.NotifyUser.TelegramID, rep.NotifyUser.Reply)
	}

	callbackText := h.renderer.Translate("", rep.CallbackText)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *AdminCallbacks) HandleAdminRejectApplication(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: не передан ID пользователя", ShowAlert: true})
	}
	targetUserID := args[0]

	rep, err := h.uc.RejectApplication(cc.ctx, cc.user, targetUserID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при отклонении", ShowAlert: true})
	}

	suffix := h.renderer.Translate("", rep.AdminEditSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix
	_, _ = c.Bot().Edit(c.Message(), updatedText, telebot.ModeMarkdown)

	if rep.NotifyUser != nil {
		sendReplyToUser(c.Bot(), h.renderer, rep.NotifyUser.TelegramID, rep.NotifyUser.Reply)
	}

	callbackText := h.renderer.Translate("", rep.CallbackText)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *AdminCallbacks) HandleAdminApproveNoiseUpgrade(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: неверные аргументы", ShowAlert: true})
	}
	targetUserID := args[0]
	category := user.NoiseLevel(args[1])

	rep, err := h.uc.ApproveNoiseUpgrade(cc.ctx, cc.user, targetUserID, category)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при одобрении апгрейда", ShowAlert: true})
	}

	suffix := h.renderer.Translate("", rep.AdminEditSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix
	_, _ = c.Bot().Edit(c.Message(), updatedText, telebot.ModeMarkdown)

	if rep.NotifyUser != nil {
		sendReplyToUser(c.Bot(), h.renderer, rep.NotifyUser.TelegramID, rep.NotifyUser.Reply)
	}

	callbackText := h.renderer.Translate("", rep.CallbackText)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *AdminCallbacks) HandleAdminRejectNoiseUpgrade(c telebot.Context) error {
	cc, err := extractCtx(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка: не передан ID пользователя", ShowAlert: true})
	}
	targetUserID := args[0]

	rep, err := h.uc.RejectNoiseUpgrade(cc.ctx, cc.user, targetUserID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Ошибка при отклонении апгрейда", ShowAlert: true})
	}

	suffix := h.renderer.Translate("", rep.AdminEditSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix
	_, _ = c.Bot().Edit(c.Message(), updatedText, telebot.ModeMarkdown)

	if rep.NotifyUser != nil {
		sendReplyToUser(c.Bot(), h.renderer, rep.NotifyUser.TelegramID, rep.NotifyUser.Reply)
	}

	callbackText := h.renderer.Translate("", rep.CallbackText)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}
