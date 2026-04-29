package adminuser

import (
	"context"
	"strconv"
	"strings"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
)

func (uc *Usecase) SearchStart(ctx context.Context, actor *user.User) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminUserSearch, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminUsersPromptSearch},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminPanel}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) OnText(ctx context.Context, actor *user.User, text string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.ClearState(ctx, actor.TelegramID)

	input := strings.TrimSpace(text)

	var target *user.User
	var err error
	if strings.HasPrefix(input, "@") {
		target, err = uc.users.GetByUsername(ctx, strings.TrimPrefix(input, "@"))
	} else {
		var id int64
		id, err = strconv.ParseInt(input, 10, 64)
		if err != nil {
			return response.Reply{
				Kind: response.KindEdit,
				Text: response.Text{Key: keys.TextAdminUsersMsgInvalid},
			}, nil
		}
		target, err = uc.users.GetByTelegramID(ctx, id)
	}

	if err != nil || target == nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
		}, nil
	}


	actionText := "admin.users.search.ban"
	actionUnique := keys.BtnAdminUserBan
	if target.Status == user.StatusBanned {
		actionText = "admin.users.search.unban"
		actionUnique = keys.BtnAdminUserUnban
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminUsersSearchResult,
			Args: map[string]any{
				"name": target.Name,
				"id":   target.ID,
			},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: actionText},
						Data: response.CallbackData{Unique: actionUnique, Args: []string{target.ID}},
					},
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{Unique: keys.BtnAdminUsers},
					},
				},
			},
		},
	}, nil
}
