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
		Text: response.Text{Key: keys.TextProfileEditNamePrompt},
	}, nil
}

func (uc *Usecase) OnText(ctx context.Context, user *user.User, text string) (response.Reply, error) {
	state, _ := uc.state.GetState(ctx, user.TelegramID)
	if state == keys.StateProfileNoiseReason {
		return uc.NoiseUpgradeSubmit(ctx, user, text)
	}

	user.Name = text
	user.UpdatedAt = time.Now()

	if err := uc.users.Update(ctx, user); err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextProfileEditNameMsgErr},
		}, err
	}

	_ = uc.state.ClearState(ctx, user.TelegramID)

	rep, err := uc.Profile(ctx, user)
	rep.Kind = response.KindSend
	return rep, err
}
