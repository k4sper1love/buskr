package callbacks

import (
	"context"
	"errors"
	"strconv"

	"github.com/k4sper1love/buskr/internal/config"
	bookingDomain "github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/auth"
	"github.com/k4sper1love/buskr/internal/usecase/booking"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type Booking struct {
	uc       *booking.Usecase
	authUc   *auth.Usecase
	renderer *render.Renderer
}

func NewBooking(uc *booking.Usecase, authUc *auth.Usecase, renderer *render.Renderer) *Booking {
	return &Booking{
		uc:       uc,
		authUc:   authUc,
		renderer: renderer,
	}
}

// Uc exposes the underlying usecase for one-off use in bot.go inline handlers.
func (h *Booking) Uc() *booking.Usecase {
	return h.uc
}

func (h *Booking) HandleBook(c telebot.Context) error {
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

func (h *Booking) HandleBookDateSelected(c telebot.Context) error {
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

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	_ = h.uc.SaveBookingMessageID(ctx, u.TelegramID, c.Message().ID)
	return nil
}

func (h *Booking) HandleBookLocationSelected(c telebot.Context) error {
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

func (h *Booking) HandleBookSlotSelected(c telebot.Context) error {
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

func (h *Booking) HandleBookDurationSelected(c telebot.Context) error {
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

	lang := ""
	if s := c.Sender(); s != nil {
		lang = s.LanguageCode
	}

	err = h.renderer.Render(c, rep.Reply)
	if err != nil {
		return err
	}

	if rep.Location.Latitude != 0 && rep.Location.Longitude != 0 {
		locationMsg := &telebot.Location{Lat: float32(rep.Location.Latitude), Lng: float32(rep.Location.Longitude)}
		_, _ = c.Bot().Send(c.Recipient(), locationMsg)
	}

	var callbackText string
	if rep.Callback.Key != "" || rep.Callback.Fallback != "" {
		callbackText = h.renderer.Translate(lang, rep.Callback)
	}
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *Booking) HandleBookingBackToLocs(c telebot.Context) error {
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

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	_ = h.uc.SaveBookingMessageID(ctx, u.TelegramID, c.Message().ID)
	return nil
}

func (h *Booking) HandleBookingBackToSlots(c telebot.Context) error {
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

func (h *Booking) HandleBookingList(c telebot.Context) error {
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

func (h *Booking) HandleBookingDetails(c telebot.Context) error {
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

func (h *Booking) HandleBookingCancelConfirm(c telebot.Context) error {
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

	_, err = h.uc.CancelConfirm(ctx, u, bookingID)
	if err != nil {
		return err
	}

	lang := ""
	if s := c.Sender(); s != nil {
		lang = s.LanguageCode
	}
	toastText := h.renderer.Translate(lang, response.Text{Key: keys.TextBookDetMsgCancelSuccess})
	_ = c.Respond(&telebot.CallbackResponse{Text: toastText})

	listRep, err := h.uc.List(ctx, u)
	if err != nil {
		return err
	}
	listRep.Kind = response.KindEdit

	return h.renderer.Render(c, listRep)
}

func (h *Booking) HandleBookingCheckIn(c telebot.Context) error {
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
		lang := ""
		if s := c.Sender(); s != nil {
			lang = s.LanguageCode
		}
		var alertText string
		if errors.Is(err, bookingDomain.ErrInvalidStatus) {
			alertText = h.renderer.Translate(lang, response.Text{
				Key:      keys.TextBookErrInvalidStatus,
				Fallback: "Чек-ин недоступен: время вышло или бронь отменена",
			})
			_ = c.Respond(&telebot.CallbackResponse{Text: alertText, ShowAlert: true})
			// Append error text and remove keyboard
			updatedText := c.Message().Text + "\n\n" + alertText
			_, _ = c.Bot().Edit(c.Message(), updatedText, telebot.ModeHTML, &telebot.ReplyMarkup{})
			return nil
		} else {
			alertText = h.renderer.Translate(lang, response.Text{
				Key:      keys.TextCommonErrGeneral,
				Fallback: "Что-то пошло не так. Повторите попытку.",
			})
		}
		return c.Respond(&telebot.CallbackResponse{Text: alertText, ShowAlert: true})
	}

	suffix := h.renderer.Translate("", rep.SuccessSuffix)
	updatedText := c.Message().Text + "\n\n" + suffix

	_, err = c.Bot().Edit(c.Message(), updatedText, telebot.ModeHTML, &telebot.ReplyMarkup{})
	if err != nil {
		return err
	}

	callbackText := h.renderer.Translate("", rep.Callback)
	return c.Respond(&telebot.CallbackResponse{Text: callbackText})
}

func (h *Booking) HandleBookingGrabHotSlot(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	if !config.MustLoad().Booking.EnableHotSlots {
		lang := ""
		if s := c.Sender(); s != nil {
			lang = s.LanguageCode
		}
		alertText := h.renderer.Translate(lang, response.Text{
			Key:      keys.TextBookErrHotSlotsDisabled,
			Fallback: "Горящие слоты временно отключены / Hot slots are temporarily disabled",
		})
		return c.Respond(&telebot.CallbackResponse{Text: alertText, ShowAlert: true})
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
		lang := ""
		if s := c.Sender(); s != nil {
			lang = s.LanguageCode
		}
		var alertKey string
		var fallbackText string
		switch {
		case errors.Is(err, bookingDomain.ErrSlotTaken):
			alertKey = keys.TextBookErrSlotTaken
			fallbackText = "Временной слот уже занят."
		case errors.Is(err, bookingDomain.ErrTimeOverlap):
			alertKey = keys.TextBookErrTimeOverlap
			fallbackText = "У вас уже есть бронирование на это время."
		case errors.Is(err, bookingDomain.ErrMaxActiveBookings):
			alertKey = keys.TextBookErrMaxActive
			fallbackText = "Достигнут лимит активных бронирований."
		case errors.Is(err, bookingDomain.ErrMaxBookingsPerLocation):
			alertKey = keys.TextBookErrMaxPerLoc
			fallbackText = "Лимит бронирований на этой точке исчерпан."
		case errors.Is(err, bookingDomain.ErrNoiseExceeded):
			alertKey = keys.TextBookErrNoiseExceeded
			fallbackText = "Превышен лимит шума для этой точки."
		case errors.Is(err, bookingDomain.ErrNoisyNeighbor):
			alertKey = keys.TextBookErrNoisyNeighbor
			fallbackText = "Рядом запланировано другое громкое выступление."
		case errors.Is(err, bookingDomain.ErrInvalidTime):
			alertKey = keys.TextBookErrInvalidTime
			fallbackText = "Неверное время выступления."
		default:
			alertKey = keys.TextCommonErrGeneral
			fallbackText = "Не удалось занять слот. Попробуйте еще раз."
		}
		alertText := h.renderer.Translate(lang, response.Text{Key: alertKey, Fallback: fallbackText})
		return c.Respond(&telebot.CallbackResponse{Text: alertText, ShowAlert: true})
	}

	c.Respond(&telebot.CallbackResponse{})
	if rep.IsEmpty() {
		return nil
	}

	err = h.renderer.Render(c, rep)
	if err != nil {
		return err
	}

	return nil
}

func (h *Booking) HandleBookingCancelFlow(c telebot.Context) error {
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
	toastText := h.renderer.Translate(lang, response.Text{Key: keys.TextBookCreateMsgCancel})
	_ = c.Respond(&telebot.CallbackResponse{Text: toastText})

	startRep, err := h.authUc.Start(ctx, u, "")
	if err != nil {
		return err
	}
	startRep.Kind = response.KindEdit

	return h.renderer.Render(c, startRep)
}

func (h *Booking) HandleBookScheduleStart(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	rep, err := h.uc.ScheduleStart(ctx, u)
	if err != nil {
		return err
	}

	_ = h.uc.SaveBookingMessageID(ctx, u.TelegramID, c.Message().ID)

	_ = c.Respond(&telebot.CallbackResponse{})
	return h.renderer.Render(c, rep)
}

func (h *Booking) HandleBookScheduleLocSel(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	locID := c.Args()[0]
	rep, err := h.uc.ScheduleForUser(ctx, u, locID, "")
	if err != nil {
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{})
	return h.renderer.Render(c, rep)
}

func (h *Booking) HandleBookScheduleDaySel(c telebot.Context) error {
	ctx := c.Get("ctx").(context.Context)
	u, err := ctxkey.GetUser(c)
	if err != nil {
		return err
	}

	locID := c.Args()[0]
	dateStr := c.Args()[1]
	rep, err := h.uc.ScheduleForUser(ctx, u, locID, dateStr)
	if err != nil {
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{})
	return h.renderer.Render(c, rep)
}
