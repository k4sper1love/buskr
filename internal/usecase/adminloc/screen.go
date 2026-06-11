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

	var mapBtn response.Button
	if webAppURL != "" {
		mapBtn = response.Button{
			Text:      response.Text{Key: keys.TextAdminLocsBtnMap},
			WebAppURL: webAppURL,
		}
	} else {
		mapBtn = response.Button{
			Text: response.Text{Key: keys.TextAdminLocsBtnMap},
			Data: response.CallbackData{Unique: keys.BtnAdminLocsMap},
		}
	}

	var rows [][]response.Button

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAdminLocsBtnAdd},
			Data: response.CallbackData{Unique: keys.BtnAdminLocAdd},
		},
	})

	rows = append(rows, []response.Button{
		mapBtn,
	})

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

	mapURL := fmt.Sprintf("%s?q=%f,%f", baseMapURL, loc.Coords.Lat, loc.Coords.Lon)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminLocsDetails,
			Args: map[string]any{
				"name":      loc.Name,
				"desc":      loc.Description,
				"max_noise": loc.MaxNoise,
				"status":    loc.Status,
				"lat":       fmt.Sprint("%.6f", loc.Coords.Lat),
				"lon":       fmt.Sprint("%.6f", loc.Coords.Lon),
			},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: toogleKey},
						Data: response.CallbackData{Unique: keys.BtnAdminLocTog, Args: []string{loc.ID}},
					},
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnEdit},
						Data: response.CallbackData{Unique: keys.BtnAdminLocEdit, Args: []string{loc.ID}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnSchedule},
						Data: response.CallbackData{Unique: keys.BtnAdminLocSchedule, Args: []string{loc.ID}},
					},
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnOpenMap},
						URL:  mapURL,
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextAdminLocsBtnDelete},
						Data: response.CallbackData{Unique: keys.BtnAdminLocDel, Args: []string{loc.ID}},
					},
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
					{Text: response.Text{Key: keys.TextCommonLblNoiseLight}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitLight)}}},
					{Text: response.Text{Key: keys.TextCommonLblNoiseMedium}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitMedium)}}},
					{Text: response.Text{Key: keys.TextCommonLblNoiseHard}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditNoise, Args: []string{locID, string(location.LimitHard)}}},
				},
				{
					{Text: response.Text{Key: keys.TextAdminLocsEditBtnGeo}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditGeo, Args: []string{locID}}},
				},
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}}},
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

func (uc *Usecase) Schedule(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	loc, err := uc.locs.GetByID(ctx, locID)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	locTz, _ := time.LoadLocation("Asia/Almaty")
	now := time.Now().In(locTz)

	bookings, err := uc.bookings.GetScheduleForLocation(ctx, locID, now)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	backKeyboard := response.Keyboard{
		InlineRows: [][]response.Button{{{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}},
		}}},
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
			Text:     response.Text{Key: keys.TextAdminLocsScheduleEmpty},
			Keyboard: backKeyboard,
		}, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📅 *%s*\n%s\n\n", loc.Name, now.Format("02.01.2006")))
	for _, b := range active {
		icon := "⏳"
		if b.Status == booking.StatusActive {
			icon = "✅"
		}
		sb.WriteString(fmt.Sprintf("%s `%s – %s`\n",
			icon,
			b.StartTime.In(locTz).Format("15:04"),
			b.EndTime.In(locTz).Format("15:04"),
		))
	}

	return response.Reply{
		Kind:     response.KindEdit,
		Text:     response.Text{Fallback: sb.String()},
		Keyboard: backKeyboard,
	}, nil
}

func (uc *Usecase) AllOnMap(ctx context.Context, actor *user.User) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	locs, err := uc.locs.GetLocationsForAdmin(ctx)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	if len(locs) == 0 {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextAdminLocsMapEmpty}}, nil
	}

	imgBytes, numberedList, err := uc.maps.Generate(locs)
	if err != nil {
		return response.Reply{Kind: response.KindEdit, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
	}

	caption := fmt.Sprintf("📍 %d\n\n%s", len(locs), numberedList)

	return response.Reply{
		Kind:  response.KindSendImage,
		Image: imgBytes,
		Text:  response.Text{Fallback: caption},
	}, nil
}

func (uc *Usecase) Delete(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	err := uc.locs.DeleteLocation(ctx, locID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Fallback: "❌ Невозможно удалить локацию, так как на неё уже есть бронирования. Вместо этого вы можете её отключить."},
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
