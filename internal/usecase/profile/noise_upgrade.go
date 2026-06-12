package profile

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) NoiseUpgrade(ctx context.Context, u *user.User) (response.Reply, error) {
	// Если мы зашли сюда (например, вернулись назад), сбрасываем стейт ввода причины
	state, _ := uc.state.GetState(ctx, u.TelegramID)
	if state == keys.StateProfileNoiseReason {
		_ = uc.state.ClearState(ctx, u.TelegramID)
		_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataProfileRequestedNoise)
	}

	// Защита от спама: проверяем, есть ли уже активная заявка в Redis
	var isPending bool
	_ = uc.state.GetData(ctx, u.TelegramID, "noise_upgrade_pending", &isPending)
	if isPending {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextProfileUpgAlreadyPending},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextCommonBtnBack},
							Data: response.CallbackData{Unique: keys.BtnProfileMain},
						},
					},
				},
			},
		}, nil
	}

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
	// Сохраняем выбранный уровень шума в Redis
	err := uc.state.SetData(ctx, u.TelegramID, keys.DataProfileRequestedNoise, string(requestedNoise), uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	// Переводим пользователя в стейт ожидания причины
	err = uc.state.SetState(ctx, u.TelegramID, keys.StateProfileNoiseReason, uc.ttl)
	if err != nil {
		return response.Reply{}, err
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileUpgPromptReason},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextProfileUpgBtnSkip},
						Data: response.CallbackData{
							Unique: keys.BtnProfileNoiseSkipReason,
						},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnBack},
						Data: response.CallbackData{
							Unique: keys.BtnProfileNoiseUp,
						},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) NoiseUpgradeSubmit(ctx context.Context, u *user.User, reason string) (response.Reply, error) {
	var requestedNoiseStr string
	err := uc.state.GetData(ctx, u.TelegramID, keys.DataProfileRequestedNoise, &requestedNoiseStr)
	if err != nil || requestedNoiseStr == "" {
		_ = uc.state.ClearState(ctx, u.TelegramID)
		return uc.NoiseUpgrade(ctx, u)
	}

	requestedNoise := user.NoiseLevel(requestedNoiseStr)

	// Очищаем стейт FSM и сохраненные данные выбора шума
	_ = uc.state.ClearState(ctx, u.TelegramID)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataProfileRequestedNoise)

	// Защита от спама: ставим флаг в Redis с TTL 48 часов
	err = uc.state.SetData(ctx, u.TelegramID, "noise_upgrade_pending", true, 48*time.Hour)
	if err != nil {
		return response.Reply{}, err
	}

	if err := uc.admin.NewNoiseUpgrade(ctx, NoiseUpgradePayload{
		TelegramUsername: u.Username,
		UserID:           u.ID,
		Name:             u.Name,
		CurrentNoise:     u.NoiseLevel,
		RequestedNoise:   requestedNoise,
		Karma:            u.Karma,
		Reason:           reason,
	}); err != nil {
		_ = uc.state.ClearData(ctx, u.TelegramID, "noise_upgrade_pending")
		return response.Reply{}, err
	}

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileUpgMsgSuccess},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnMenu},
						Data: response.CallbackData{
							Unique: keys.BtnCommonMenuSend,
						},
					},
				},
			},
		},
	}, nil
}
