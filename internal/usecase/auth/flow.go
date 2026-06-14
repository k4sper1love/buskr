package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/tz"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) Start(ctx context.Context, u *user.User, payload string) (response.Reply, error) {
	if payload != "" {
		if u.Status == user.StatusActive {
			return uc.Start(ctx, u, "")
		}

		err := uc.users.ProcessInvite(ctx, u.ID, payload)
		if err != nil {
			var key string
			switch {
			case errors.Is(err, user.ErrInviteNotFound):
				key = keys.TextAuthInvitedErrNotFound
			case errors.Is(err, user.ErrInviteAlreadyUsed):
				key = keys.TextAuthInvitedErrAlreadyUsed
			case errors.Is(err, user.ErrInviteExpired):
				key = keys.TextAuthInvitedErrExpired
			case errors.Is(err, user.ErrInvalidStatus):
				key = keys.TextAuthInvitedErrActive
			default:
				key = keys.TextAuthInvitedErrGeneric
			}

			return response.Reply{
				Kind: response.KindSend,
				Text: response.Text{Key: key},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{
							{
								Text: response.Text{Key: keys.TextCommonBtnBack},
								Data: response.CallbackData{Unique: keys.BtnCommonMenu},
							},
						},
					},
				},
			}, nil
		}

		_ = uc.state.SetState(ctx, u.TelegramID, keys.StateAuthInvitedName, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAuthInvitedPromptName},
		}, nil
	}

	switch u.Status {
	case user.StatusGuest:
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{
				Key: keys.TextAuthGuestTitle,
				Args: map[string]any{
					"bot_name": uc.botName,
				},
			},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextAuthGuestBtnApply},
							Data: response.CallbackData{Unique: keys.BtnAuthApply},
						},
					},
				},
			},
		}, nil

	case user.StatusPending:
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAuthPendingTitle},
		}, nil

	case user.StatusBanned:
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAuthBannedTitle},
		}, nil

	case user.StatusActive:
		var rows [][]response.Button

		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextAuthActiveBtnBook},
				Data: response.CallbackData{Unique: keys.BtnBookStart},
			},
		})

		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextCommonBtnSchedule},
				Data: response.CallbackData{Unique: keys.BtnBookSchedule},
			},
		})

		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextAuthActiveBtnBookings},
				Data: response.CallbackData{Unique: keys.BtnBookList},
			},
			{
				Text: response.Text{Key: keys.TextAuthActiveBtnProfile},
				Data: response.CallbackData{Unique: keys.BtnProfileMain},
			},
		})

		if u.Role == user.RoleAdmin {
			rows = append(rows, []response.Button{
				{
					Text: response.Text{Key: keys.TextAuthActiveBtnAdmin},
					Data: response.CallbackData{Unique: keys.BtnAdminPanel},
				},
			})
		}

		now := tz.Now()

		weekdayKeys := []string{
			keys.TextCommonWeekdaySunFull,
			keys.TextCommonWeekdayMonFull,
			keys.TextCommonWeekdayTueFull,
			keys.TextCommonWeekdayWedFull,
			keys.TextCommonWeekdayThuFull,
			keys.TextCommonWeekdayFriFull,
			keys.TextCommonWeekdaySatFull,
		}

		monthKeys := []string{
			"",
			keys.TextCommonMonthJan,
			keys.TextCommonMonthFeb,
			keys.TextCommonMonthMar,
			keys.TextCommonMonthApr,
			keys.TextCommonMonthMay,
			keys.TextCommonMonthJun,
			keys.TextCommonMonthJul,
			keys.TextCommonMonthAug,
			keys.TextCommonMonthSep,
			keys.TextCommonMonthOct,
			keys.TextCommonMonthNov,
			keys.TextCommonMonthDec,
		}

		name := u.Username
		if u.Name != "" {
			name = u.Name
		}

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{
				Key: keys.TextAuthActiveTitle,
				Args: map[string]any{
					"name":    name,
					"day":     now.Day(),
					"month":   monthKeys[now.Month()],
					"weekday": weekdayKeys[now.Weekday()],
				},
				SubKeyArgs: []string{"month", "weekday"},
			},
			Keyboard: response.Keyboard{
				InlineRows: rows,
			},
		}, nil
	}

	return response.Reply{}, nil
}

func (uc *Usecase) OnText(ctx context.Context, u *user.User, text string) (response.Reply, error) {
	u.Name = text
	u.UpdatedAt = time.Now()

	if err := uc.users.Update(ctx, u); err != nil {
		return response.Reply{}, err
	}

	_ = uc.state.ClearState(ctx, u.TelegramID)

	// Notify admins if registered via invite
	if invite, err := uc.users.GetInviteByUsedByID(ctx, u.ID); err == nil && invite != nil && invite.CreatedBy != "" {
		invitedBy := "Admin"
		if creator, err := uc.users.GetByID(ctx, invite.CreatedBy); err == nil {
			invitedBy = creator.Name
			if creator.Username != "" {
				invitedBy += " (@" + creator.Username + ")"
			}
		}

		username := u.Username
		if username == "" {
			username = "—"
		} else if !strings.HasPrefix(username, "@") {
			username = "@" + username
		}

		_ = uc.notifier.NewInviteRegistration(ctx, username, u.Name, string(u.NoiseLevel), invitedBy)
	}

	return uc.Start(ctx, u, "")
}

func (uc *Usecase) ApplyButton(ctx context.Context, u *user.User) (response.Reply, error) {
	if u.Status != user.StatusGuest {
		return response.Reply{}, nil
	}

	// hand off to onboarding flow — sets the first onboarding fsm state.
	_ = uc.state.SetState(ctx, u.TelegramID, keys.StateOnboardName, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextOnboardStep1PromptName},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnCancel},
						Data: response.CallbackData{
							Unique: keys.BtnOnboardCancel,
						},
					},
				},
			},
		},
	}, nil
}
