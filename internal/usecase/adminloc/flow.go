package adminloc

import (
	"context"
	"fmt"
	"strings"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

const (
	nearbyRadius = 100 // meters
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

	case keys.StateAdminLocEditName:
		var locID string
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, &locID)

		if locID == "" {
			return response.Reply{}, nil
		}

		loc, err := uc.locs.GetByID(ctx, locID)
		if err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
		}

		if err := uc.locs.UpdateLocation(ctx, locID, text, loc.Description, loc.Coords.Lat, loc.Coords.Lon, loc.MaxNoise); err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextAdminLocsEditMsgErr}}, nil
		}

		_ = uc.state.ClearState(ctx, actor.TelegramID)
		_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocEditID)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsEditMsgSuccess},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}}}},
				},
			},
		}, nil

	case keys.StateAdminLocEditDesc:
		var locID string
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, &locID)

		if locID == "" {
			return response.Reply{}, nil
		}

		loc, err := uc.locs.GetByID(ctx, locID)
		if err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
		}

		if err := uc.locs.UpdateLocation(ctx, locID, loc.Name, text, loc.Coords.Lat, loc.Coords.Lon, loc.MaxNoise); err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextAdminLocsEditMsgErr}}, nil
		}

		_ = uc.state.ClearState(ctx, actor.TelegramID)
		_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocEditID)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsEditMsgSuccess},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{},
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

	switch st {
	case keys.StateAdminLocGeo:
		var name, desc, noise string
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocName, &name)
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocDesc, &desc)
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocNoise, &noise)

		_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocLat, lat, uc.ttl)
		_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocLon, lon, uc.ttl)

		nearby, err := uc.locs.FindNearby(ctx, float64(lat), float64(lon), nearbyRadius)
		if err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
		}

		if len(nearby) > 0 {
			var lines []string
			for _, l := range nearby {
				icon := "🟢"
				if l.Status == location.StatusInactive {
					icon = "🔴"
				}
				lines = append(lines, fmt.Sprintf("%s %s", icon, l.Name))
			}

			return response.Reply{
				Kind: response.KindSend,
				Text: response.Text{
					Key: keys.TextAdminLocsNearbyWarn,
					Args: map[string]any{
						"radius": 100,
						"list":   strings.Join(lines, "\n"),
					},
				},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{
							{Text: response.Text{Key: keys.TextAdminLocsNearbyConfirm}, Data: response.CallbackData{Unique: keys.BtnAdminLocConfirm}},
							{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocCancel}},
						},
					},
				},
			}, nil
		}

		return uc.doCreate(ctx, actor, name, desc, noise, float64(lat), float64(lon))

	case keys.BtnAdminLocEditGeo:
		var locID string
		_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, &locID)
		if locID == "" {
			return response.Reply{}, nil
		}

		loc, err := uc.locs.GetByID(ctx, locID)
		if err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextCommonErrGeneral}}, nil
		}

		if err := uc.locs.UpdateLocation(ctx, locID, loc.Name, loc.Description, float64(lat), float64(lon), loc.MaxNoise); err != nil {
			return response.Reply{Kind: response.KindSend, Text: response.Text{Key: keys.TextAdminLocsEditMsgErr}}, nil
		}

		_ = uc.state.ClearState(ctx, actor.TelegramID)
		_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocEditID)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsEditMsgSuccess},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminLocDet, Args: []string{locID}}}},
				},
			},
		}, nil
	default:
		return response.Reply{}, nil
	}
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

func (uc *Usecase) EditNameStart(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, locID, uc.ttl)
	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocEditName, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsEditStepName},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditCancel, Args: []string{locID}}}},
			},
		},
	}, nil
}

func (uc *Usecase) EditDescStart(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}
	_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, locID, uc.ttl)
	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocEditDesc, uc.ttl)
	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsEditStepDesc},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditCancel, Args: []string{locID}}}},
			},
		},
	}, nil
}

func (uc *Usecase) EditGeoStart(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}
	_ = uc.state.SetData(ctx, actor.TelegramID, keys.DataAdminLocEditID, locID, uc.ttl)
	_ = uc.state.SetState(ctx, actor.TelegramID, keys.StateAdminLocEditGeo, uc.ttl)
	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextAdminLocsEditStepGeo},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{{Text: response.Text{Key: keys.TextCommonBtnCancel}, Data: response.CallbackData{Unique: keys.BtnAdminLocEditCancel, Args: []string{locID}}}},
			},
		},
	}, nil
}

func (uc *Usecase) CancelEditFlow(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	_ = uc.state.ClearState(ctx, actor.TelegramID)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocEditID)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocLat)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocLon)

	return uc.Details(ctx, actor, locID)
}

func (uc *Usecase) ConfirmCreate(ctx context.Context, actor *user.User, locID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	var name, desc, noise string
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocName, &name)
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocDesc, &desc)
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocNoise, &noise)
	var lat, lon float64
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocLat, &lat)
	_ = uc.state.GetData(ctx, actor.TelegramID, keys.DataAdminLocLon, &lon)

	return uc.doCreate(ctx, actor, name, desc, noise, lat, lon)
}

func (uc *Usecase) doCreate(ctx context.Context, actor *user.User, name, desc, noise string, lat, lon float64) (response.Reply, error) {
	_, err := uc.locs.CreateLocation(ctx, name, desc, float64(lat), float64(lon), location.NoiseLimit(noise))
	if err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextAdminLocsMsgAddErr},
		}, nil
	}

	_ = uc.state.ClearState(ctx, actor.TelegramID)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocName)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocDesc)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocNoise)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocLat)
	_ = uc.state.ClearData(ctx, actor.TelegramID, keys.DataAdminLocLon)

	return response.Reply{
		Kind: response.KindSend,
		Text: response.Text{Key: keys.TextAdminLocsMsgAddSuccess},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{{Text: response.Text{Key: keys.TextAdminLocsBtnList}, Data: response.CallbackData{Unique: keys.BtnAdminLocs}}},
			},
		},
	}, nil
}
