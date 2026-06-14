package booking

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/mdutil"
	"github.com/k4sper1love/buskr/internal/tz"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) List(ctx context.Context, u *user.User) (response.Reply, error) {
	list, err := uc.bookings.GetPendingOrActiveBookings(ctx, u.ID)
	if err != nil {
		return response.Reply{}, err
	}

	var rows [][]response.Button

	if len(list) == 0 {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextBookListMsgEmpty},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{
						Text: response.Text{Key: keys.TextCommonBtnMenu},
						Data: response.CallbackData{Unique: keys.BtnCommonMenu},
					}},
				},
			},
		}, nil
	}

	for _, bk := range list {
		locData, _ := uc.locs.GetByID(ctx, bk.LocationID)

		btnText := fmt.Sprintf("📍 %s | %s, %02d:00-%02d:00",
			locData.Name,
			tz.In(bk.StartTime).Format("02.01"),
			tz.In(bk.StartTime).Hour(),
			tz.In(bk.EndTime).Hour(),
		)

		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: btnText},
				Data: response.CallbackData{Unique: keys.BtnBookDetails, Args: []string{bk.ID}},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnMenu},
			Data: response.CallbackData{Unique: keys.BtnCommonMenu},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookListTitle},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) Details(ctx context.Context, u *user.User, bookingID string) (response.Reply, error) {
	bk, err := uc.bookings.GetByID(ctx, bookingID)
	if err != nil {
		return response.Reply{}, err
	}

	locData, _ := uc.locs.GetByID(ctx, bk.LocationID)
	now := tz.Now()

	statusKey := keys.TextBookDetLblUnknown
	switch bk.Status {
	case booking.StatusPending:
		statusKey = keys.TextBookDetLblPending
	case booking.StatusActive:
		statusKey = keys.TextBookDetLblActive
	case booking.StatusCompleted:
		statusKey = keys.TextBookDetLblCompleted
	case booking.StatusCancelled:
		statusKey = keys.TextBookDetLblCancelled
	case booking.StatusNoShow:
		statusKey = keys.TextBookDetLblNoshow
	}

	var rows [][]response.Button

	if bk.Status == booking.StatusPending {
		timeUntilStart := tz.In(bk.StartTime).Sub(now)
		if timeUntilStart <= 60*time.Minute && timeUntilStart >= -15*time.Minute {
			rows = append(rows, []response.Button{
				{
					Text: response.Text{Key: keys.TextBookDetBtnCheckin},
					Data: response.CallbackData{Unique: keys.BtnBookCheckin, Args: []string{bk.ID}},
				},
			})
		}
	}

	if locData != nil && locData.Coords.Lat != 0 && locData.Coords.Lon != 0 {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextBookCreateBtn2GIS},
				URL:  fmt.Sprintf("https://2gis.kz/geo/%f,%f", locData.Coords.Lon, locData.Coords.Lat),
			},
			{
				Text: response.Text{Key: keys.TextBookCreateBtnYandex},
				URL:  fmt.Sprintf("https://yandex.kz/maps/?text=%f,%f", locData.Coords.Lat, locData.Coords.Lon),
			},
		})
	}

	if now.Before(tz.In(bk.StartTime)) && bk.Status == booking.StatusPending {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextBookDetBtnCancel},
				Data: response.CallbackData{Unique: keys.BtnBookCancelConf, Args: []string{bk.ID}},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextBookDetBtnList},
			Data: response.CallbackData{Unique: keys.BtnBookList},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookDetTitle, Args: map[string]any{
			"name":   locData.Name,
			"date":   tz.In(bk.StartTime).Format("02.01"),
			"time":   fmt.Sprintf("%02d:00 - %02d:00", tz.In(bk.StartTime).Hour(), tz.In(bk.EndTime).Hour()),
			"status": statusKey,
		},
			SubKeyArgs: []string{"status"}, // statusKey is itself an i18n key
		},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) CancelConfirm(ctx context.Context, u *user.User, bookingID string) (response.Reply, error) {
	err := uc.bookings.CancelBooking(ctx, bookingID, u.ID)
	if err != nil {
		return response.Reply{}, err
	}

	return response.Reply{}, nil
}

func (uc *Usecase) CheckIn(ctx context.Context, u *user.User, bookingID string) (CheckInResult, error) {
	err := uc.bookings.CheckIn(ctx, bookingID, u.ID)
	if err != nil {
		return CheckInResult{}, err
	}

	return CheckInResult{
		SuccessSuffix: response.Text{Key: keys.TextBookDetMsgCheckinSfx},
		Callback:      response.Text{Key: keys.TextBookDetMsgCheckinCb},
	}, nil
}

func (uc *Usecase) GrabHotSlot(ctx context.Context, u *user.User, locID string, startHour int, durationHours int) (response.Reply, error) {
	now := tz.Now()

	startTime := now
	endTime := time.Date(now.Year(), now.Month(), now.Day(), startHour+durationHours, 0, 0, 0, tz.Location())

	_, err := uc.bookings.BookSlot(ctx, u.ID, locID, startTime, endTime)
	if err != nil {
		return response.Reply{}, err
	}

	locData, _ := uc.locs.GetByID(ctx, locID)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextBookDetMsgGrabTitle,
			Args: map[string]any{
				"name":     locData.Name,
				"end_time": fmt.Sprintf("%02d:00", tz.In(endTime).Hour()),
			},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnMenu},
						Data: response.CallbackData{Unique: keys.BtnBookList},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) ScheduleStart(ctx context.Context, u *user.User) (response.Reply, error) {
	if u.Status != user.StatusActive {
		return response.Reply{}, fmt.Errorf("user is not active")
	}

	locations, err := uc.locs.GetLocationsForMusicians(ctx)
	if err != nil {
		return response.Reply{}, err
	}

	var webAppLocs []webAppLocation
	for _, loc := range locations {
		webAppLocs = append(webAppLocs, webAppLocation{
			ID:       loc.ID,
			Name:     loc.Name,
			Desc:     loc.Description,
			Lat:      loc.Coords.Lat,
			Lon:      loc.Coords.Lon,
			MaxNoise: string(loc.MaxNoise),
		})
	}

	var webAppURL string
	if uc.webAppURL != "" && len(webAppLocs) > 0 {
		jsonData, err := json.Marshal(webAppLocs)
		if err == nil {
			encoded := base64.URLEncoding.EncodeToString(jsonData)
			webAppURL = fmt.Sprintf("%s?v=%d#bot=%s&mode=schedule&locs=%s", uc.webAppURL, time.Now().Unix(), uc.botUsername, encoded)
		}
	}

	var lastLoc *location.Location
	lastBooking, err := uc.bookings.GetLastBookingByUser(ctx, u.ID)
	if err == nil && lastBooking != nil {
		for _, loc := range locations {
			if loc.ID == lastBooking.LocationID {
				lastLoc = loc
				break
			}
		}
	}

	var rows [][]response.Button
	if webAppURL != "" {
		rows = append(rows, []response.Button{
			{
				Text:      response.Text{Key: keys.TextCommonBtnOpenMap},
				WebAppURL: webAppURL,
			},
		})
	}

	if lastLoc != nil {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{
					Key:  keys.TextCommonBtnLastLoc,
					Args: map[string]any{"name": lastLoc.Name},
				},
				Data: response.CallbackData{
					Unique: keys.BtnBookScheduleLocSel,
					Args:   []string{lastLoc.ID},
				},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAuthActiveBtnSuggestLoc},
			Data: response.CallbackData{Unique: keys.BtnSuggestLocStart},
		},
	})

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnCommonMenu},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextBookSchedulePromptLoc},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) ScheduleForUser(ctx context.Context, u *user.User, locID string, dateStr string) (response.Reply, error) {
	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{}, err
	}

	targetDate := tz.Now()
	if dateStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", dateStr, tz.Location()); err == nil {
			targetDate = parsed
		}
	}

	bookings, err := uc.bookings.GetScheduleForLocation(ctx, locID, targetDate)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	now := tz.Now()
	var rows [][]response.Button

	for i := 0; i < uc.maxAdvanceDays; i++ {
		tDate := now.AddDate(0, 0, i)
		dateValue := tDate.Format("2006-01-02")
		dateHuman := tDate.Format("02.01")

		isActive := tDate.Year() == targetDate.Year() && tDate.YearDay() == targetDate.YearDay()

		var key string
		args := map[string]any{"date": dateHuman}
		var subKeyArgs []string

		switch i {
		case 0:
			if isActive {
				key = keys.TextAdminLocsScheduleLblTodayActive
			} else {
				key = keys.TextBookCreateLblToday
			}
		case 1:
			if isActive {
				key = keys.TextAdminLocsScheduleLblTomorrowActive
			} else {
				key = keys.TextBookCreateLblTomorrow
			}
		default:
			weekdayKeys := []string{
				keys.TextCommonWeekdaySun,
				keys.TextCommonWeekdayMon,
				keys.TextCommongWeekdayTue,
				keys.TextCommonWeekdayWed,
				keys.TextCommonWeekdayThu,
				keys.TextCommonWeekdayFri,
				keys.TextCommonWeekdaySat,
			}
			args["weekday"] = weekdayKeys[tDate.Weekday()]
			subKeyArgs = []string{"weekday"}
			if isActive {
				key = keys.TextAdminLocsScheduleLblOtherActive
			} else {
				key = keys.TextBookCreateLblOther
			}
		}

		rows = append(rows, []response.Button{
			{
				Text: response.Text{
					Key:        key,
					Args:       args,
					SubKeyArgs: subKeyArgs,
				},
				Data: response.CallbackData{
					Unique: keys.BtnBookScheduleDaySel,
					Args:   []string{locID, dateValue},
				},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnBookSchedule},
		},
	})

	keyboard := response.Keyboard{
		InlineRows: rows,
	}

	var active []*booking.Booking
	for _, b := range bookings {
		if b.Status == booking.StatusPending || b.Status == booking.StatusActive {
			active = append(active, b)
		}
	}

	if len(active) == 0 {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{
				Key: keys.TextBookScheduleEmptyForDay,
				Args: map[string]any{
					"name": loc.Name,
					"date": targetDate.Format("02.01.2006"),
				},
			},
			Keyboard: keyboard,
		}, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 <b>%s</b>\n%s\n\n", loc.Name, targetDate.Format("02.01.2006")))
	for _, b := range active {
		icon := "⏳"
		if b.Status == booking.StatusActive {
			icon = "✅"
		}

		var userPart string
		uData, err := uc.users.GetByID(ctx, b.UserID)
		if err == nil && uData != nil {
			var displayName string
			if uData.Name != "" {
				displayName = uData.Name
			} else {
				displayName = "User"
			}
			if uData.Username != "" {
				displayName += fmt.Sprintf(" (@%s)", uData.Username)
			}
			userPart = " - " + mdutil.Escape(displayName)
		}

		sb.WriteString(fmt.Sprintf("%s <code>%s – %s</code> %s\n",
			icon,
			tz.In(b.StartTime).Format("15:04"),
			tz.In(b.EndTime).Format("15:04"),
			userPart,
		))
	}

	return response.Reply{
		Kind:     response.KindEdit,
		Text:     response.Text{Fallback: sb.String()},
		Keyboard: keyboard,
	}, nil
}
