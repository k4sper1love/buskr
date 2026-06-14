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
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type Auth struct {
	uc             *auth.Usecase
	bookingUc      *booking.Usecase
	adminLocUc     *adminloc.Usecase
	renderer       *render.Renderer
	supportContact string
}

func NewAuth(uc *auth.Usecase, bookingUc *booking.Usecase, adminLocUc *adminloc.Usecase, renderer *render.Renderer, supportContact string) *Auth {
	return &Auth{
		uc:             uc,
		bookingUc:      bookingUc,
		adminLocUc:     adminLocUc,
		renderer:       renderer,
		supportContact: supportContact,
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
	if strings.HasPrefix(payload, "admloc_") || strings.HasPrefix(payload, "loc_") || strings.HasPrefix(payload, "locsch_") {
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

		// Retrieve stored message ID of the location selection interface and delete it
		if msgID, err := h.bookingUc.GetBookingMessageID(ctx, u.TelegramID); err == nil && msgID != 0 {
			_ = c.Bot().Delete(&telebot.Message{
				ID:   msgID,
				Chat: &telebot.Chat{ID: u.TelegramID},
			})
			_ = h.bookingUc.ClearBookingMessageID(ctx, u.TelegramID)
		}

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

	if strings.HasPrefix(payload, "locsch_") {
		locID := strings.TrimPrefix(payload, "locsch_")

		if msgID, err := h.bookingUc.GetBookingMessageID(ctx, u.TelegramID); err == nil && msgID != 0 {
			_ = c.Bot().Delete(&telebot.Message{
				ID:   msgID,
				Chat: &telebot.Chat{ID: u.TelegramID},
			})
			_ = h.bookingUc.ClearBookingMessageID(ctx, u.TelegramID)
		}

		rep, err := h.bookingUc.ScheduleForUser(ctx, u, locID, "")
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

func (h *Auth) HandleHelp(c telebot.Context) error {
	rep := response.Reply{
		Kind: response.KindSend,
		Text: response.Text{
			Key: keys.TextHelpTitle,
			Args: map[string]any{
				"support": h.supportContact,
			},
		},
	}
	return h.renderer.Render(c, rep)
}

