package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/user"
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
