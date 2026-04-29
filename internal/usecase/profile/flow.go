package profile

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
)

func (uc *Usecase) EditNameStart(ctx context.Context, user *user.User) (response.Reply, error) {
	_ = uc.state.SetState(ctx, user.TelegramID, keys.StateProfileEditName, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileMainBtnEditName},
	}, nil
}

func (uc *Usecase) OnText(ctx context.Context, user *user.User, text string) (response.Reply, error) {
	user.Name = text
	user.UpdatedAt = time.Now()

	if err := uc.users.Update(ctx, user); err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextProfileEditNameMsgErr},
		}, err
	}

	_ = uc.state.ClearState(ctx, user.TelegramID)

	return response.Reply{
		Kind: response.KindSend,
		Text: response.Text{Key: keys.TextProfileEditNameMsgSuccess},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnProfileMain}},
				},
			},
		},
	}, nil
}
