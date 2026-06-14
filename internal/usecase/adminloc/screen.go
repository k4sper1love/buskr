package adminloc

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

const (
	pageSize   = 5
	baseMapURL = "https://maps.google.com/"
)

type webAppLocation struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Desc   string  `json:"desc"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	Status string  `json:"status"`
}

func (uc *Usecase) List(ctx context.Context, actor *user.User, page int) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	locs, err := uc.locs.GetLocationsForAdmin(ctx)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	var webAppURL string
	if uc.webAppURL != "" && len(locs) > 0 {
		var webAppLocs []webAppLocation
		for _, loc := range locs {
			webAppLocs = append(webAppLocs, webAppLocation{
				ID:     loc.ID,
				Name:   loc.Name,
				Desc:   loc.Description,
				Lat:    loc.Coords.Lat,
				Lon:    loc.Coords.Lon,
				Status: string(loc.Status),
			})
		}
		jsonData, err := json.Marshal(webAppLocs)
		if err == nil {
			encoded := base64.URLEncoding.EncodeToString(jsonData)
			webAppURL = fmt.Sprintf("%s?v=%d#bot=%s&mode=admin&locs=%s", uc.webAppURL, time.Now().Unix(), uc.botUsername, encoded)
		}
	}

	var rows [][]response.Button

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAdminLocsBtnAdd},
			Data: response.CallbackData{Unique: keys.BtnAdminLocAdd},
		},
	})

	if webAppURL != "" {
		rows = append(rows, []response.Button{
			{
				Text:      response.Text{Key: keys.TextCommonBtnOpenMap},
				WebAppURL: webAppURL,
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnAdminPanel},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsTitle},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) Details(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminLocsMsgNotFound},
		}, nil
	}

	toogleKey := keys.TextAdminLocsBtnDisable
	if loc.Status == location.StatusInactive {
		toogleKey = keys.TextAdminLocsBtnEnable
	}

	map2GISURL := fmt.Sprintf("https://2gis.kz/geo/%f,%f", loc.Coords.Lon, loc.Coords.Lat)
	mapYandexURL := fmt.Sprintf("https://yandex.kz/maps/?text=%f,%f", loc.Coords.Lat, loc.Coords.Lon)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminLocsDetails,
			Args: map[string]any{
				"name":      loc.Name,
				"desc":      loc.Description,
				"max_noise": "common.lbl.noise_" + string(loc.MaxNoise),
				"status":    "common.lbl.loc_status_" + string(loc.Status),
				"lat":       fmt.Sprintf("%.6f", loc.Coords.Lat),
				"lon":       fmt.Sprintf("%.6f", loc.Coords.Lon),
			},
			SubKeyArgs: []string{"max_noise", "status"},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextBookCreateBtn2GIS},
						URL:  map2GISURL,
					},
					{
						Text: response.Text{Key: keys.TextBookCreateBtnYandex},
						URL:  mapYandexURL,
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnSchedule},
						Data: response.CallbackData{Unique: keys.BtnAdminLocSchedule, Args: []string{loc.ID}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnEdit},
						Data: response.CallbackData{Unique: keys.BtnAdminLocEdit, Args: []string{loc.ID}},
					},
					{
						Text: response.Text{Key: toogleKey},
						Data: response.CallbackData{Unique: keys.BtnAdminLocTog, Args: []string{loc.ID}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnDelete},
						Data: response.CallbackData{Unique: keys.BtnAdminLocDel, Args: []string{loc.ID}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnList},
						Data: response.CallbackData{Unique: keys.BtnAdminLocs},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) ToggleStatus(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	newStatus := location.StatusInactive
	if loc.Status == location.StatusInactive {
		newStatus = location.StatusActive
	}

	if err := uc.locs.ChangeStatus(ctx, locID, newStatus); err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.Details(ctx, actor, locID)
}

func (uc *Usecase) EditMenu(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key:  keys.TextAdminLocsEditTitle,
			Args: map[string]any{"name": loc.Name},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextAdminLocsEditBtnName}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditName, Args: []string{locID}}},
					{Text: response.Text{Key: keys.TextAdminLocsEditBtnDesc}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditDesc, Args: []string{locID}}},
				},
				{
					{Text: response.Text{Key: keys.TextAdminLocsEditBtnNoise}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoiseMenu, Args: []string{locID}}},
					{Text: response.Text{Key: keys.TextAdminLocsEditBtnGeo}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditGeo, Args: []string{locID}}},
				},
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) EditNoiseMenu(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key:  keys.TextAdminLocsEditNoiseTitle,
			Args: map[string]any{"name": loc.Name, "current": "common.lbl.noise_" + string(loc.MaxNoise)},
			SubKeyArgs: []string{"current"},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonLblNoiseLight}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitLight)}}},
				},
				{
					{Text: response.Text{Key: keys.TextCommonLblNoiseMedium}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitMedium)}}},
				},
				{
					{Text: response.Text{Key: keys.TextCommonLblNoiseHard}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitHard)}}},
				},
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocEdit, Args: []string{locID}}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) EditNoiseSelected(ctx context.Context, actor *user.User, locID, noise string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	if err := uc.locs.UpdateLocation(ctx, locID, loc.Name, loc.Description, loc.Coords.Lat, loc.Coords.Lon, location.NoiseLimit(noise)); err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextAdminLocsEditMsgErr}}, nil
	}

	return uc.Details(ctx, actor, locID)
}

func (uc *Usecase) Schedule(ctx context.Context, actor *user.User, locID string, dateStr string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	var targetDate time.Time
	if dateStr != "" {
		if t, err := time.ParseInLocation("2006-01-02", dateStr, tz.Location()); err == nil {
			targetDate = t
		} else {
			targetDate = tz.Now()
		}
	} else {
		targetDate = tz.Now()
	}

	bookings, err := uc.bookings.GetScheduleForLocation(ctx, locID, targetDate)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	// Generate keyboard with dates
	now := tz.Now()
	var inlineRows [][]response.Button

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

		inlineRows = append(inlineRows, []response.Button{
			{
				Text: response.Text{
					Key:        key,
					Args:       args,
					SubKeyArgs: subKeyArgs,
				},
				Data: response.CallbackData{
					Unique: keys.BtnAdminLocSchedule,
					Args:   []string{locID, dateValue},
				},
			},
		})
	}

	// Add Back button at the bottom
	inlineRows = append(inlineRows, []response.Button{{
		Text: response.Text{Key: keys.TextCommonBtnBack},
		Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}},
	}})

	keyboard := response.Keyboard{
		InlineRows: inlineRows,
	}

	var active []*booking.Booking
	for _, b := range bookings {
		if b.Status == booking.StatusPending || b.Status == booking.StatusActive {
			active = append(active, b)
		}
	}

	if len(active) == 0 {
		return response.Reply{
			Kind:     response.KindEdit,
			Text:     response.Text{
				Key:  keys.TextAdminLocsScheduleEmptyForDay,
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
		u, err := uc.users.GetByID(ctx, b.UserID)
		if err == nil && u != nil {
			var displayName string
			if u.Name != "" {
				displayName = u.Name
			} else {
				displayName = "User"
			}
			if u.Username != "" {
				displayName += fmt.Sprintf(" (@%s)", u.Username)
			}
			userPart = " — " + mdutil.Escape(displayName)
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

func (uc *Usecase) DeleteConfirm(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminLocsMsgNotFound},
		}, nil
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminLocsDelConfirmTitle,
			Args: map[string]any{
				"name": loc.Name,
			},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsDelConfirmBtn},
						Data: response.CallbackData{Unique: keys.BtnAdminLocDelConf, Args: []string{loc.ID}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnCancel},
						Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{loc.ID}},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) DeleteExecuted(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	err := uc.locs.DeleteLocation(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminLocsErrDeleteHasBookings},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{{{
					Text: response.Text{Key: keys.TextCommonBtnBack},
					Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}},
				}}},
			},
		}, nil
	}

	return uc.List(ctx, actor, 0)
}
