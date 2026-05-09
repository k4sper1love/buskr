package adminloc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

const (
	pageSize = 5
)

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

	total := len(locs)

	start := page * pageSize
	if start > total {
		start = total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	var rows [][]response.Button

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextAdminLocsBtnAdd},
			Data: response.CallbackData{Unique: keys.BtnAdminLocAdd},
		},
		{
			Text: response.Text{Key: keys.TextAdminLocsBtnMap},
			Data: response.CallbackData{Unique: keys.BtnAdminLocsMap},
		},
	})

	for _, loc := range locs[start:end] {
		icon := "🟢"
		if loc.Status == location.StatusInactive {
			icon = "🔴"
		}
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Fallback: fmt.Sprintf("%s %s", icon, loc.Name)},
				Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{loc.ID}},
			},
		})
	}

	if total > pageSize {
		totalPages := (total + pageSize - 1) / pageSize
		var nav []response.Button

		if page > 0 {
			nav = append(nav, response.Button{
				Text: response.Text{Fallback: "◀️"},
				Data: response.CallbackData{Unique: keys.BtnAdminLocs, Args: []string{strconv.Itoa(page - 1)}},
			})
		}
		nav = append(nav, response.Button{
			Text: response.Text{Fallback: fmt.Sprintf("%d / %d", page+1, totalPages)},
			Data: response.CallbackData{Unique: keys.BtnAdminLocs, Args: []string{strconv.Itoa(page)}},
		})
		if end < total {
			nav = append(nav, response.Button{
				Text: response.Text{Fallback: "▶️"},
				Data: response.CallbackData{Unique: keys.BtnAdminLocs, Args: []string{strconv.Itoa(page + 1)}},
			})
		}
		rows = append(rows, nav)
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

	actionKey := keys.TextAdminLocsBtnDisable
	if loc.Status == location.StatusInactive {
		actionKey = keys.TextAdminLocsBtnEnable
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminLocsDetails,
			Args: map[string]any{
				"name":      loc.Name,
				"desc":      loc.Description,
				"max_noise": loc.MaxNoise,
				"status":    loc.Status,
			},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: actionKey},
						Data: response.CallbackData{Unique: keys.BtnAdminLocTog, Args: []string{loc.ID}},
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
