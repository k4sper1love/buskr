package admin

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) AdminPanel(ctx context.Context, u *user.User) (response.Reply, error) {
	if u.Role != user.RoleAdmin {
		return response.Reply{}, user.ErrAccessDenied
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminPanelTitle},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{{Text: response.Text{Key: keys.TextAdminPanelBtnInvites}, Data: response.CallbackData{Unique: keys.BtnAdminInvites}}},
				{{Text: response.Text{Key: keys.TextAdminPanelBtnLocs}, Data: response.CallbackData{Unique: keys.BtnAdminLocs}}},
				{{Text: response.Text{Key: keys.TextAdminPanelBtnUsers}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				{{Text: response.Text{Key: keys.TextCommonBtnMenu}, Data: response.CallbackData{Unique: keys.BtnCommonMenu}}},
			},
		},
	}, nil
}
