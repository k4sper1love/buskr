package callbacks

import (
	"context"
	"strconv"

	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/booking"
	"gopkg.in/telebot.v3"
)

type BookingCallbacks struct {
	uc       *booking.Usecase
	renderer *render.Renderer
}

func NewBookingCallbacks(uc *booking.Usecase, renderer *render.Renderer) *BookingCallbacks {
	return &BookingCallbacks{
		uc:       uc,
		renderer: renderer,
	}
}

// Uc exposes the underlying usecase for one-off use in bot.go inline handlers.
func (h *BookingCallbacks) Uc() *booking.Usecase {
	return h.uc
}

func (h *BookingCallbacks) HandleBook(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.Book(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookDateSelected(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	date := args[0]

	rep, err := h.uc.DateSelected(ctx, u, date)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookLocationSelected(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	location := args[0]

	rep, err := h.uc.LocationSelected(ctx, u, location)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookSlotSelected(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	hourStr := args[0]

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return err
	}

	rep, err := h.uc.SlotSelected(ctx, u, hour)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookDurationSelected(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	durationStr := args[0]

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return err
	}

	rep, err := h.uc.DurationSelected(ctx, u, duration)
	if err != nil {
		return err
	}

	err = h.renderer.Render(c, rep.Reply)
	if err != nil {
		return err
	}

	locationMsg := &telebot.Location{Lat: float32(rep.Location.Latitude), Lng: float32(rep.Location.Longitude)}
	_, _ = c.Bot().Send(c.Recipient(), locationMsg)

	callbackText := h.renderer.Translate("", rep.Callback)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *BookingCallbacks) HandleBookingBackToLocs(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.BackToLocations(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookingBackToSlots(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.BackToSlots(ctx, u)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookingList(c telebot.Context) error {
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

func (h *BookingCallbacks) HandleBookingDetails(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	bookingID := args[0]

	rep, err := h.uc.Details(ctx, u, bookingID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookingCancelConfirm(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	bookingID := args[0]

	rep, err := h.uc.CancelConfirm(ctx, u, bookingID)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookingCheckIn(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	bookingID := args[0]

	rep, err := h.uc.CheckIn(ctx, u, bookingID)
	if err != nil {
		return err
	}

	suffix := h.renderer.Translate("", rep.SuccessSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix

	_, err = c.Bot().Edit(c.Message(), updatedText, telebot.ModeMarkdown)
	if err != nil {
		return err
	}

	callbackText := h.renderer.Translate("", rep.Callback)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *BookingCallbacks) HandleBookingGrabHotSlot(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	args := c.Args()
	if len(args) < 3 {
		return c.Respond(&telebot.CallbackResponse{Text: "System error", ShowAlert: true})
	}
	locID := args[0]
	startHourStr := args[1]
	durationHoursStr := args[2]

	startHour, err := strconv.Atoi(startHourStr)
	if err != nil {
		return err
	}
	durationHours, err := strconv.Atoi(durationHoursStr)
	if err != nil {
		return err
	}

	rep, err := h.uc.GrabHotSlot(ctx, u, locID, startHour, durationHours)
	if err != nil {
		return err
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}
	return h.renderer.Render(c, rep)
}

func (h *BookingCallbacks) HandleBookingCancelFlow(c telebot.Context) error {
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
