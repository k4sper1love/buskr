package adminuser

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) UsersList(ctx context.Context, actor *user.User, page int) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	const pageSize = 5
	if page < 0 {
		page = 0
	}

	offset := page * pageSize
	users, total, err := uc.users.GetUsersPaginated(ctx, offset, pageSize)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	if offset >= total && total > 0 {
		page = totalPages - 1
		offset = page * pageSize
		users, total, err = uc.users.GetUsersPaginated(ctx, offset, pageSize)
		if err != nil {
			return response.Reply{
				Kind: response.KindEdit,
				Text: response.Text{Key: keys.TextCommonErrGeneral},
			}, nil
		}
	}

	var inlineRows [][]response.Button

	for _, u := range users {
		display := u.Name
		if u.Username != "" {
			display += fmt.Sprintf(" (@%s)", u.Username)
		} else {
			display += fmt.Sprintf(" (ID: %d)", u.TelegramID)
		}

		inlineRows = append(inlineRows, []response.Button{
			{
				Text: response.Text{Fallback: display},
				Data: response.CallbackData{Unique: keys.BtnAdminUserDetail, Args: []string{u.ID}},
			},
		})
	}

	// Pagination row
	var navRow []response.Button
	if page > 0 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: "⬅️"},
			Data: response.CallbackData{Unique: keys.BtnAdminUsersPage, Args: []string{strconv.Itoa(page - 1)}},
		})
	}
	if totalPages > 1 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: fmt.Sprintf("%d / %d", page+1, totalPages)},
			Data: response.CallbackData{Unique: "noop"},
		})
	}
	if page < totalPages-1 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: "➡️"},
			Data: response.CallbackData{Unique: keys.BtnAdminUsersPage, Args: []string{strconv.Itoa(page + 1)}},
		})
	}
	if len(navRow) > 0 {
		inlineRows = append(inlineRows, navRow)
	}

	// Search button
	inlineRows = append(inlineRows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAdminUsersBtnSearch},
			Data: response.CallbackData{Unique: keys.BtnAdminUserSearchPrompt},
		},
	})

	// Back button
	inlineRows = append(inlineRows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnAdminPanel},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key:  keys.TextAdminUsersListTitle,
			Args: map[string]any{"total": total},
		},
		Keyboard: response.Keyboard{
			InlineRows: inlineRows,
		},
	}, nil
}

func (uc *Usecase) SearchPrompt(ctx context.Context, actor *user.User) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminUserSearch, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminUsersSearchQueryPrompt},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{Unique: keys.BtnAdminUsers},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) UserDetail(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}},
					},
				},
			},
		}, nil
	}

	stats, err := uc.users.GetStats(ctx, target.ID)
	if err != nil {
		stats = &user.UserStats{UserID: target.ID}
	}

	usernameDisplay := target.Username
	if usernameDisplay != "" {
		if !strings.HasPrefix(usernameDisplay, "@") {
			usernameDisplay = "@" + usernameDisplay
		}
	} else {
		usernameDisplay = "—"
	}

	var keyboardRows [][]response.Button

	keyboardRows = append(keyboardRows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAdminUsersBtnChangeNoise},
			Data: response.CallbackData{Unique: keys.BtnAdminUserNoiseMenu, Args: []string{target.ID}},
		},
	})

	if target.ID != actor.ID {
		roleText := keys.TextAdminUsersBtnPromote
		roleUnique := keys.BtnAdminUserPromote
		if target.Role == user.RoleAdmin {
			roleText = keys.TextAdminUsersBtnDemote
			roleUnique = keys.BtnAdminUserDemote
		}
		keyboardRows = append(keyboardRows, []response.Button{
			{
				Text: response.Text{Key: roleText},
				Data: response.CallbackData{Unique: roleUnique, Args: []string{target.ID}},
			},
		})
	}

	if target.Role != user.RoleAdmin {
		actionText := keys.TextAdminUsersBtnBan
		actionUnique := keys.BtnAdminUserBan
		if target.Status == user.StatusBanned {
			actionText = keys.TextAdminUsersBtnUnban
			actionUnique = keys.BtnAdminUserUnban
		}
		keyboardRows = append(keyboardRows, []response.Button{
			{
				Text: response.Text{Key: actionText},
				Data: response.CallbackData{Unique: actionUnique, Args: []string{target.ID}},
			},
		})
	}

	keyboardRows = append(keyboardRows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnAdminUsers},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminUsersSearchResult,
			Args: map[string]any{
				"name":                target.Name,
				"username":            usernameDisplay,
				"telegram_id":         target.TelegramID,
				"role":                "common.lbl.role_" + string(target.Role),
				"status":              "common.lbl.status_" + string(target.Status),
				"noise":               "common.lbl.noise_" + string(target.NoiseLevel),
				"karma":               target.Karma,
				"total_bookings":      stats.TotalBookings,
				"successful_checkins": stats.SuccessfulCheckins,
				"no_shows":            stats.NoShows,
			},
			SubKeyArgs: []string{"role", "status", "noise"},
		},
		Keyboard: response.Keyboard{
			InlineRows: keyboardRows,
		},
	}, nil
}

func (uc *Usecase) OnText(ctx context.Context, actor *user.User, text string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.ClearState(ctx, actor.TelegramID)

	input := strings.TrimSpace(text)
	if strings.HasPrefix(input, "@") {
		input = strings.TrimPrefix(input, "@")
	}

	_ = uc.state.SetData(ctx, actor.TelegramID, "admin_search_query", input, uc.ttl)

	return uc.SearchResults(ctx, actor, input, 0, response.KindSend)
}

func (uc *Usecase) SearchPage(ctx context.Context, actor *user.User, page int) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	var query string
	err := uc.state.GetData(ctx, actor.TelegramID, "admin_search_query", &query)
	if err != nil || query == "" {
		return uc.UsersList(ctx, actor, 0)
	}

	return uc.SearchResults(ctx, actor, query, page, response.KindEdit)
}

func (uc *Usecase) SearchResults(ctx context.Context, actor *user.User, query string, page int, kind response.Kind) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	const pageSize = 5
	if page < 0 {
		page = 0
	}

	offset := page * pageSize
	results, total, err := uc.users.FindByQuery(ctx, query, offset, pageSize)
	if err != nil || len(results) == 0 {
		return response.Reply{
			Kind: kind,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextCommonBtnBack},
							Data: response.CallbackData{Unique: keys.BtnAdminUsers},
						},
					},
				},
			},
		}, nil
	}

	_ = uc.state.SetData(ctx, actor.TelegramID, "admin_search_query", query, uc.ttl)

	if len(results) == 1 && total == 1 {
		rep, err := uc.UserDetail(ctx, actor, results[0].ID)
		if err == nil {
			rep.Kind = kind
		}
		return rep, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	if offset >= total && total > 0 {
		page = totalPages - 1
		offset = page * pageSize
		results, total, err = uc.users.FindByQuery(ctx, query, offset, pageSize)
		if err != nil {
			return response.Reply{
				Kind: kind,
				Text: response.Text{Key: keys.TextCommonErrGeneral},
			}, nil
		}
	}

	var inlineRows [][]response.Button
	for _, u := range results {
		display := u.Name
		if u.Username != "" {
			display += fmt.Sprintf(" (@%s)", u.Username)
		} else {
			display += fmt.Sprintf(" (ID: %d)", u.TelegramID)
		}
		inlineRows = append(inlineRows, []response.Button{
			{
				Text: response.Text{Fallback: display},
				Data: response.CallbackData{Unique: keys.BtnAdminUserDetail, Args: []string{u.ID}},
			},
		})
	}

	// Pagination row
	var navRow []response.Button
	if page > 0 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: "⬅️"},
			Data: response.CallbackData{Unique: keys.BtnAdminUserSearchPage, Args: []string{strconv.Itoa(page - 1)}},
		})
	}
	if totalPages > 1 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: fmt.Sprintf("%d / %d", page+1, totalPages)},
			Data: response.CallbackData{Unique: "noop"},
		})
	}
	if page < totalPages-1 {
		navRow = append(navRow, response.Button{
			Text: response.Text{Fallback: "➡️"},
			Data: response.CallbackData{Unique: keys.BtnAdminUserSearchPage, Args: []string{strconv.Itoa(page + 1)}},
		})
	}
	if len(navRow) > 0 {
		inlineRows = append(inlineRows, navRow)
	}

	// Back to search prompt
	inlineRows = append(inlineRows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{Unique: keys.BtnAdminUserSearchPrompt},
		},
	})

	return response.Reply{
		Kind: kind,
		Text: response.Text{Key: keys.TextAdminUsersSearchResultsTitle},
		Keyboard: response.Keyboard{
			InlineRows: inlineRows,
		},
	}, nil
}

func (uc *Usecase) UserNoiseMenu(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key:  keys.TextAdminUsersNoiseMenuTitle,
			Args: map[string]any{"name": target.Name},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: "common.lbl.noise_light"},
						Data: response.CallbackData{Unique: keys.BtnAdminUserNoiseSet, Args: []string{target.ID, string(user.NoiseLight)}},
					},
				},
				{
					{
						Text: response.Text{Key: "common.lbl.noise_medium"},
						Data: response.CallbackData{Unique: keys.BtnAdminUserNoiseSet, Args: []string{target.ID, string(user.NoiseMedium)}},
					},
				},
				{
					{
						Text: response.Text{Key: "common.lbl.noise_hard"},
						Data: response.CallbackData{Unique: keys.BtnAdminUserNoiseSet, Args: []string{target.ID, string(user.NoiseHard)}},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{Unique: keys.BtnAdminUserDetail, Args: []string{target.ID}},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) SetUserNoise(ctx context.Context, actor *user.User, targetID string, noise string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	target.NoiseLevel = user.NoiseLevel(noise)
	target.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, target)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.UserDetail(ctx, actor, targetID)
}
