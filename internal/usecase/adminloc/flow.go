package adminloc

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) AddStart(ctx context.Context, actor *user.User) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocName, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsAddStep1},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocCancel}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) OnText(ctx context.Context, actor *user.User, text string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	st, err := uc.state.GetState(ctx, actor.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	switch st {
	case keys.StateAdminLocName:
		_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocName, text, uc.ttl)
		_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocDesc, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsAddStep2},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocCancel}},
					},
				},
			},
		}, nil

	case keys.StateAdminLocDesc:
		_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocDesc, text, uc.ttl)
		_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocNoise, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsAddStep3},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{Text: response.Text{Key: keys.TextCommonLblNoiseLight}, Data: response.CallbackData{Unique: keys.BtnAdminLocNoise, Args: []string{string(user.NoiseLight)}}},
						{Text: response.Text{Key: keys.TextCommonLblNoiseMedium}, Data: response.CallbackData{Unique: keys.BtnAdminLocNoise, Args: []string{string(user.NoiseMedium)}}},
						{Text: response.Text{Key: keys.TextCommonLblNoiseHard}, Data: response.CallbackData{Unique: keys.BtnAdminLocNoise, Args: []string{string(user.NoiseHard)}}},
					},
					{
						{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocCancel}},
					},
				},
			},
		}, nil

	default:
		return response.Reply{}, nil
	}
}

func (uc *Usecase) OnNoiseSelected(ctx context.Context, actor *user.User, noise string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocNoise, noise, uc.ttl)
	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocGeo, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsAddStep4},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocCancel}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) OnLocation(ctx context.Context, actor *user.User, lat, lon float32) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	st, err := uc.state.GetState(ctx, actor.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}
	if st != keys.StateAdminLocGeo {
		return response.Reply{}, nil
	}

	var name, desc, noise string
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocName, &name)
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocDesc, &desc)
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocNoise, &noise)

	_, err = uc.locs.CreateLocation(ctx, name, desc, float64(lat), float64(lon), location.NoiseLimit(noise))
	if err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsMsgAddErr},
		}, nil // why error is not returned?
	}

	_ = uc.state.ClearState(ctx, actor.TelegramID)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocName)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocDesc)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocNoise)

	return response.Reply{
		Kind: response.KindSend,
		Text: response.Text{Key: keys.TextAdminLocsMsgAddSuccess},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextAdminLocsBtnList}, Data: response.CallbackData{Unique: keys.BtnAdminLocs}},
				},
			},
		},
	}, nil
}

func (uc *Usecase) CancelFlow(ctx context.Context, actor *user.User) (response.Reply, error) {
	_ = uc.state.ClearState(ctx, actor.TelegramID)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocName)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocDesc)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocNoise)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsMsgCancel},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocs}},
				},
			},
		},
	}, nil
}
