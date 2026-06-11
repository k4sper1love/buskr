package booking

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"

	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) Book(ctx context.Context, u *user.User) (response.Reply, error) {
	if u.Status != user.StatusActive {
		return response.Reply{}, errors.New("user is not active")
	}

	err := uc.state.SetState(ctx, u.TelegramID, keys.StateBookDate, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	// clear any previous stale booking data just in case
	uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingDate)
	uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingLoc)

	loc, _ := time.LoadLocation("Asia/Almaty")
	now := time.Now().In(loc)

	var rows [][]response.Button
	for i := 0; i < uc.maxAdvanceDays; i++ {
		targetDate := now.AddDate(0, 0, i)
		dateValue := targetDate.Format("2006-01-02")
		dateHuman := targetDate.Format("02.01")

		key := keys.TextBookCreateLblOther
		switch i {
		case 0:
			key = keys.TextBookCreateLblToday
		case 1:
			key = keys.TextBookCreateLblTomorrow
		}

		rows = append(rows, []response.Button{
			{
				Text: response.Text{
					Key: key,
					Args: map[string]any{
						"date": dateHuman,
					},
				},
				Data: response.CallbackData{
					Unique: keys.BtnBookDateSel,
					Args:   []string{dateValue},
				},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnCancel},
			Data: response.CallbackData{
				Unique: keys.BtnBookCancel,
			},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookCreatePromptDate},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

type webAppLocation struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Desc string  `json:"desc"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func (uc *Usecase) DateSelected(ctx context.Context, u *user.User, date string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateBookDate {
		return response.Reply{}, errors.New("invalid state")
	}

	err = uc.state.SetData(ctx, u.TelegramID, keys.DataBookingDate, date, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	err = uc.state.SetState(ctx, u.TelegramID, keys.StateBookLoc, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	locations, err := uc.locs.GetLocationsForMusicians(ctx)
	if err != nil {
		return response.Reply{}, err
	}

	var webAppLocs []webAppLocation
	for _, loc := range locations {
		if u.NoiseLevel.Weight() > loc.MaxNoise.Weight() {
			continue
		}

		webAppLocs = append(webAppLocs, webAppLocation{
			ID:   loc.ID,
			Name: loc.Name,
			Desc: loc.Description,
			Lat:  loc.Coords.Lat,
			Lon:  loc.Coords.Lon,
		})
	}

	if len(webAppLocs) == 0 {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextBookCreateMsgNoLocs},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextCommonBtnCancel},
							Data: response.CallbackData{Unique: keys.BtnBookCancel},
						},
					},
				},
			},
		}, nil
	}

	var webAppURL string
	if uc.webAppURL != "" {
		jsonData, err := json.Marshal(webAppLocs)
		if err == nil {
			encoded := base64.URLEncoding.EncodeToString(jsonData)
			webAppURL = fmt.Sprintf("%s?v=%d#bot=%s&locs=%s", uc.webAppURL, time.Now().Unix(), uc.botUsername, encoded)
		}
	}

	// Fetch user's last booked location for a quick shortcut
	var lastLoc *location.Location
	lastBooking, err := uc.bookings.GetLastBookingByUser(ctx, u.ID)
	if err == nil && lastBooking != nil {
		for _, loc := range locations {
			if loc.ID == lastBooking.LocationID && u.NoiseLevel.Weight() <= loc.MaxNoise.Weight() {
				lastLoc = loc
				break
			}
		}
	}

	var rows [][]response.Button
	if webAppURL != "" {
		rows = append(rows, []response.Button{
			{
				Text:      response.Text{Key: keys.TextBookCreateBtnOpenMap},
				WebAppURL: webAppURL,
			},
		})
	}

	if lastLoc != nil {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{
					Fallback: fmt.Sprintf("🕒 Последнее место: %s", lastLoc.Name),
				},
				Data: response.CallbackData{
					Unique: keys.BtnBookLocSel,
					Args:   []string{lastLoc.ID},
				},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextBookCreateBtnDates},
			Data: response.CallbackData{Unique: keys.BtnBookStart},
		},
		{
			Text: response.Text{Key: keys.TextCommonBtnCancel},
			Data: response.CallbackData{Unique: keys.BtnBookCancel},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookCreatePromptLoc},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) LocationSelected(ctx context.Context, u *user.User, locID string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateBookLoc {
		return response.Reply{}, errors.New("invalid state")
	}

	var selectedDateStr string
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingDate, &selectedDateStr)
	if err != nil || selectedDateStr == "" {
		return response.Reply{}, errors.New("no selected date")
	}

	err = uc.state.SetData(ctx, u.TelegramID, keys.DataBookingLoc, locID, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	err = uc.state.SetState(ctx, u.TelegramID, keys.StateBookSlot, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	locTz, _ := time.LoadLocation("Asia/Almaty")
	targetDate, err := time.ParseInLocation("2006-01-02", selectedDateStr, locTz)
	if err != nil {
		return response.Reply{}, err
	}

	schedule, err := uc.bookings.GetScheduleForLocation(ctx, locID, targetDate)
	if err != nil {
		return response.Reply{}, err
	}

	var rows [][]response.Button
	var currentRow []response.Button

	endHour := 22
	if targetDate.Weekday() == time.Saturday || targetDate.Weekday() == time.Sunday {
		endHour = 23
	}

	now := time.Now().In(locTz)
	isToday := targetDate.YearDay() == now.YearDay() && targetDate.Year() == now.Year()

	for hour := 10; hour < endHour; hour++ {
		if isToday && hour <= now.Hour() {
			continue // skip already passed hours on today
		}

		slotStart := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), hour, 0, 0, 0, locTz)
		slotEnd := slotStart.Add(1 * time.Hour)

		isOccupied := false
		for _, bookingObj := range schedule {
			if bookingObj.Status != booking.StatusPending && bookingObj.Status != booking.StatusActive {
				continue
			}

			if bookingObj.StartTime.Before(slotEnd) && bookingObj.EndTime.After(slotStart) {
				isOccupied = true
				break
			}
		}

		if !isOccupied {
			currentRow = append(currentRow, response.Button{
				Text: response.Text{
					Fallback: fmt.Sprintf("%02d:00", hour),
				},
				Data: response.CallbackData{
					Unique: keys.BtnBookSlotSel,
					Args:   []string{fmt.Sprintf("%d", hour)},
				},
			})

			if len(currentRow) == 3 {
				rows = append(rows, currentRow)
				currentRow = nil
			}
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	if len(rows) == 0 {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextBookCreateMsgNoSlots},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnBookCancel}},
					},
				},
			},
		}, nil
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextBookCreateBtnLocs},
			Data: response.CallbackData{Unique: keys.BtnBookBackToLocs},
		},
		{
			Text: response.Text{Key: keys.TextCommonBtnCancel},
			Data: response.CallbackData{Unique: keys.BtnBookCancel},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookCreatePromptSlot},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) SlotSelected(ctx context.Context, u *user.User, hour int) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateBookSlot {
		return response.Reply{}, errors.New("invalid state")
	}

	var selectedDateStr, locID string
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingDate, &selectedDateStr)
	if err != nil || selectedDateStr == "" {
		return response.Reply{}, errors.New("no selected date")
	}
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingLoc, &locID)
	if err != nil || locID == "" {
		return response.Reply{}, errors.New("no selected location")
	}

	err = uc.state.SetData(ctx, u.TelegramID, keys.DataBookingSlot, hour, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	err = uc.state.SetState(ctx, u.TelegramID, keys.StateBookDuration, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	locTz, _ := time.LoadLocation("Asia/Almaty")
	targetDate, _ := time.ParseInLocation("2006-01-02", selectedDateStr, locTz)

	schedule, _ := uc.bookings.GetScheduleForLocation(ctx, locID, targetDate)

	endHourLimit := 22
	if targetDate.Weekday() == time.Saturday || targetDate.Weekday() == time.Sunday {
		endHourLimit = 23
	}

	maxAvailableHours := 0
	for h := 1; h <= 3; h++ {
		checkHour := hour + h - 1
		if checkHour >= endHourLimit {
			break
		}

		slotStart := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), checkHour, 0, 0, 0, locTz)
		slotEnd := slotStart.Add(1 * time.Hour)

		isOccupied := false
		for _, bookingObj := range schedule {
			if bookingObj.Status != booking.StatusPending && bookingObj.Status != booking.StatusActive {
				continue
			}

			if bookingObj.StartTime.Before(slotEnd) && bookingObj.EndTime.After(slotStart) {
				isOccupied = true
				break
			}
		}

		if isOccupied {
			break
		}
		maxAvailableHours++
	}

	var rows [][]response.Button
	var currentRow []response.Button

	for i := 1; i <= maxAvailableHours; i++ {
		currentRow = append(currentRow, response.Button{
			Text: response.Text{
				Key: keys.TextBookCreateLblHours,
				Args: map[string]any{
					"count": i,
				},
			},
			Data: response.CallbackData{
				Unique: keys.BtnBookDurSel,
				Args:   []string{fmt.Sprintf("%d", i)},
			},
		})
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextBookCreateBtnSlots},
			Data: response.CallbackData{Unique: keys.BtnBookBackToSlots},
		},
		{
			Text: response.Text{Key: keys.TextCommonBtnCancel},
			Data: response.CallbackData{Unique: keys.BtnBookCancel},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookCreatePromptDur},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) DurationSelected(ctx context.Context, u *user.User, duration int) (BookingResult, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return BookingResult{}, err
	}

	if st != keys.StateBookDuration {
		return BookingResult{}, errors.New("invalid state")
	}

	var dateStr, locID string
	var startHour int
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingDate, &dateStr)
	if err != nil || dateStr == "" {
		return BookingResult{}, errors.New("no selected date")
	}
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingLoc, &locID)
	if err != nil || locID == "" {
		return BookingResult{}, errors.New("no selected location")
	}
	err = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingSlot, &startHour)
	if err != nil || startHour == 0 {
		return BookingResult{}, errors.New("no selected hour")
	}

	locTz, _ := time.LoadLocation("Asia/Almaty")
	targetDate, _ := time.ParseInLocation("2006-01-02", dateStr, locTz)

	startTime := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), startHour, 0, 0, 0, locTz)
	endTime := startTime.Add(time.Duration(duration) * time.Hour)

	_, err = uc.bookings.BookSlot(ctx, u.ID, locID, startTime, endTime)
	if err != nil {
		if key := bookingErrKey(err); key != "" {
			return BookingResult{Reply: response.Reply{
				Kind: response.KindEdit,
				Text: response.Text{Key: key},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{
							{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnBookCancel}},
						},
					},
				},
			}}, nil
		}
		return BookingResult{}, err
	}

	uc.state.ClearState(ctx, u.TelegramID)
	uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingDate)
	uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingLoc)
	uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingSlot)

	locData, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return BookingResult{}, err
	}

	return BookingResult{
		Reply: response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextBookCreateMsgSuccess, Args: map[string]any{
				"name": locData.Name,
				"date": targetDate.Format("02.01"),
				"time": fmt.Sprintf("%02d:00 - %02d:00", startHour, startHour+duration),
			}},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextCommonBtnMenu},
							Data: response.CallbackData{Unique: keys.BtnCommonMenuSend},
						},
					},
				},
			},
		},
		Location: SendLocation{
			Latitude:  locData.Coords.Lat,
			Longitude: locData.Coords.Lon,
		},
		Callback: response.Text{
			Key: keys.TextBookCreateMsgSuccessCb,
		},
	}, nil
}

func (uc *Usecase) CancelFlow(ctx context.Context, u *user.User) (response.Reply, error) {
	_ = uc.state.ClearState(ctx, u.TelegramID)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingDate)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingLoc)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataBookingSlot)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookCreateMsgCancel},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnMenu},
						Data: response.CallbackData{Unique: keys.BtnCommonMenuSend},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) BackToLocations(ctx context.Context, u *user.User) (response.Reply, error) {
	// Revert the FSM to StateBookLoc step.
	// The user clicked "Back to Locations" from the Slot Selection screen.
	// The Date was already selected and stored in the FSM state.

	// We need to fetch the saved date from state and pass it back to DateSelected.
	// Reset the actual state to DateBookDate just to pretend we're stepping perfectly forward.
	_ = uc.state.SetState(ctx, u.TelegramID, keys.StateBookDate, uc.ttl)

	var dateStr string
	err := uc.state.GetData(ctx, u.TelegramID, keys.DataBookingDate, &dateStr)
	if err != nil || dateStr == "" {
		// If we lost state, just restart the book flow
		return uc.Book(ctx, u)
	}

	// This will render the Location Selection screen using the saved date
	return uc.DateSelected(ctx, u, dateStr)
}

func (uc *Usecase) BackToSlots(ctx context.Context, u *user.User) (response.Reply, error) {
	// Revert the FSM to StateBookSlot step.
	// The user clicked "Back to Slots" from the Duration Selection screen.

	// Reset the state perfectly backward
	_ = uc.state.SetState(ctx, u.TelegramID, keys.StateBookLoc, uc.ttl)

	var locID string
	err := uc.state.GetData(ctx, u.TelegramID, keys.DataBookingLoc, &locID)
	if err != nil || locID == "" {
		// If we lost state, fallback to locations step
		return uc.BackToLocations(ctx, u)
	}

	// This will render the Slot Selection screen using the saved location
	return uc.LocationSelected(ctx, u, locID)
}

func (uc *Usecase) ProcessMapSelection(ctx context.Context, u *user.User, locID string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateBookLoc {
		return uc.Book(ctx, u)
	}

	locations, err := uc.locs.GetLocationsForMusicians(ctx)
	if err != nil {
		return response.Reply{}, err
	}

	var found bool
	for _, loc := range locations {
		if loc.ID == locID && u.NoiseLevel.Weight() <= loc.MaxNoise.Weight() {
			found = true
			break
		}
	}

	if !found {
		var selectedDateStr string
		_ = uc.state.GetData(ctx, u.TelegramID, keys.DataBookingDate, &selectedDateStr)
		if selectedDateStr == "" {
			return uc.Book(ctx, u)
		}

		_ = uc.state.SetState(ctx, u.TelegramID, keys.StateBookDate, uc.ttl)

		reply, err := uc.DateSelected(ctx, u, selectedDateStr)
		if err != nil {
			return response.Reply{}, err
		}
		reply.Text.Fallback = "⚠️ Please select a valid location.\n\n" + reply.Text.Fallback
		return reply, nil
	}

	return uc.LocationSelected(ctx, u, locID)
}

func bookingErrKey(err error) string {
	switch {
	case errors.Is(err, booking.ErrSlotTaken), errors.Is(err, booking.ErrTimeOverlap):
		return keys.TextBookErrSlotTaken
	case errors.Is(err, booking.ErrNoisyNeighbor):
		return keys.TextBookErrNoisyNeighbor
	case errors.Is(err, booking.ErrNoiseExceeded):
		return keys.TextBookErrNoiseExceeded
	case errors.Is(err, booking.ErrMaxActiveBookings):
		return keys.TextBookErrMaxActive
	case errors.Is(err, booking.ErrMaxBookingsPerLocation):
		return keys.TextBookErrMaxPerLoc
	case errors.Is(err, booking.ErrTooFarInFuture):
		return keys.TextBookErrTooFarInFuture
	default:
		return ""
	}
}
