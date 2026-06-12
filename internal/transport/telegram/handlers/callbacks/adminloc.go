package callbacks

import (
	"context"
	"strconv"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/adminloc"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
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

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	_ = h.uc.SaveAdminLocMessageID(ctx, u.TelegramID, c.Message().ID)
	return nil
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

	page := 0
	if args := c.Args(); len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p >= 0 {
			page = p
		}
	}

	rep, err := h.uc.List(ctx, u, page)
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

	_, err = h.uc.CancelFlow(ctx, u)
	if err != nil {
		return err
	}

	lang := ""
	if s := c.Sender(); s != nil {
		lang = s.LanguageCode
	}

	callbackText := h.renderer.Translate(lang, response.Text{Key: keys.TextAdminLocsMsgCancel})
	c.Respond(&telebot.CallbackResponse{Text: callbackText})

	_ = h.uc.ClearAdminLocMessageID(ctx, u.TelegramID)

	rep, err := h.uc.List(ctx, u, 0)
	if err != nil {
		return err
	}
	rep.Kind = response.KindEdit

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEdit(c telebot.Context) error {
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

	rep, err := h.uc.EditMenu(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditName(c telebot.Context) error {
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

	rep, err := h.uc.EditNameStart(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditDesc(c telebot.Context) error {
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

	rep, err := h.uc.EditDescStart(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditNoise(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 2 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	locID := args[0]
	noise := args[1]

	rep, err := h.uc.EditNoiseSelected(ctx, u, locID, noise)
	if err != nil {
		return err
	}

	callbackText := h.renderer.Translate("", response.Text{Key: keys.TextAdminLocsEditMsgSuccess})
	c.Respond(&telebot.CallbackResponse{Text: callbackText})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditGeo(c telebot.Context) error {
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

	rep, err := h.uc.EditGeoStart(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditCancel(c telebot.Context) error {
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

	rep, err := h.uc.CancelEditFlow(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocSchedule(c telebot.Context) error {
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

	var dateStr string
	if len(args) >= 2 {
		dateStr = args[1]
	}

	rep, err := h.uc.Schedule(ctx, u, locID, dateStr)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocConfirm(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.ConfirmCreate(ctx, u, "")
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	// Remove inline keyboard from the warning message
	_ = c.Edit(c.Message().Text, c.Message().Entities, &telebot.ReplyMarkup{})

	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocDelete(c telebot.Context) error {
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

	rep, err := h.uc.DeleteConfirm(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocDeleteConfirm(c telebot.Context) error {
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

	rep, err := h.uc.DeleteExecuted(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}

func (h *AdminLoc) HandleAdminLocEditNoiseMenu(c telebot.Context) error {
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

	rep, err := h.uc.EditNoiseMenu(ctx, u, locID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}
