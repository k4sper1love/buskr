package profile

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) NoiseUpgrade(ctx context.Context, u *user.User) (response.Reply, error) {
	var rows [][]response.Button

	if u.NoiseLevel != user.NoiseLight {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextCommonLblNoiseLight},
				Data: response.CallbackData{
					Unique: keys.BtnProfileNoiseSel,
					Args:   []string{string(user.NoiseLight)},
				},
			},
		})
	}
	if u.NoiseLevel != user.NoiseMedium {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextCommonLblNoiseMedium},
				Data: response.CallbackData{
					Unique: keys.BtnProfileNoiseSel,
					Args:   []string{string(user.NoiseMedium)},
				},
			},
		})
	}
	if u.NoiseLevel != user.NoiseHard {
		rows = append(rows, []response.Button{
			{
				Text: response.Text{Key: keys.TextCommonLblNoiseHard},
				Data: response.CallbackData{
					Unique: keys.BtnProfileNoiseSel,
					Args:   []string{string(user.NoiseHard)},
				},
			},
		})
	}

	rows = append(rows, []response.Button{
		{
			Text: response.Text{Key: keys.TextCommonBtnBack},
			Data: response.CallbackData{
				Unique: keys.BtnProfileMain,
			},
		},
	})

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileMainBtnUpgNoise},
		Keyboard: response.Keyboard{
			InlineRows: rows,
		},
	}, nil
}

func (uc *Usecase) NoiseUpgradeSelected(ctx context.Context, u *user.User, requestedNoise user.NoiseLevel) (response.Reply, error) {
	if err := uc.admin.NewNoiseUpgrade(ctx, NoiseUpgradePayload{
		TelegramUsername: u.Username,
		UserID:           u.ID,
		Name:             u.Name,
		CurrentNoise:     u.NoiseLevel,
		RequestedNoise:   requestedNoise,
		Karma:            u.Karma,
	}); err != nil {
		return response.Reply{}, err
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileUpgMsgSuccess},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{
							Unique: keys.BtnProfileMain,
						},
					},
				},
			},
		},
	}, nil
}
