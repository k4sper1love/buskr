package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) Invites(ctx context.Context, actor *user.User) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, errors.New("access denied")
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminInvTitle},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonLblNoiseLight},
						Data: response.CallbackData{
							Unique: keys.BtnAdminInvGen,
							Args:   []string{string(user.NoiseLight)},
						},
					},
					{
						Text: response.Text{Key: keys.TextCommonLblNoiseMedium},
						Data: response.CallbackData{
							Unique: keys.BtnAdminInvGen,
							Args:   []string{string(user.NoiseMedium)},
						},
					},
					{
						Text: response.Text{Key: keys.TextCommonLblNoiseHard},
						Data: response.CallbackData{
							Unique: keys.BtnAdminInvGen,
							Args:   []string{string(user.NoiseHard)},
						},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{
							Unique: keys.BtnAdminPanel,
						},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) GenerateInvite(ctx context.Context, actor *user.User, botUsername string, noise string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, errors.New("access denied")
	}

	invite, err := uc.users.GenerateInvite(ctx, user.NoiseLevel(noise), uc.inviteTTL)
	if err != nil {
		return response.Reply{}, err
	}

	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, invite.Token)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{
			Key: keys.TextAdminInvCreated,
			Args: map[string]any{
				"noise": "common.lbl.noise_" + noise,
				"link":  inviteLink,
			},
			SubKeyArgs: []string{"noise"},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{
							Unique: keys.BtnAdminPanelSend,
						},
					},
				},
			},
		},
	}, nil
}
