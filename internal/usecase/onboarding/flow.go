package onboarding

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) OnText(ctx context.Context, u *user.User, text string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	switch st {
	case keys.StateOnboardName:
		_ = uc.state.SetData(ctx, u.TelegramID, keys.DataOnboardName, text, uc.ttl)
		_ = uc.state.SetState(ctx, u.TelegramID, keys.StateOnboardNoise, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextOnboardStep2PromptNoise},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{
						{
							Text: response.Text{Key: keys.TextCommonLblNoiseLight},
							Data: response.CallbackData{
								Unique: keys.BtnOnboardNoiseSel,
								Args:   []string{string(user.NoiseLight)},
							},
						},
						{

							Text: response.Text{Key: keys.TextCommonLblNoiseMedium},
							Data: response.CallbackData{
								Unique: keys.BtnOnboardNoiseSel,
								Args:   []string{string(user.NoiseMedium)},
							},
						},
						{
							Text: response.Text{Key: keys.TextCommonLblNoiseHard},
							Data: response.CallbackData{
								Unique: keys.BtnOnboardNoiseSel,
								Args:   []string{string(user.NoiseHard)},
							},
						},
					},
					{
						{
							Text: response.Text{Key: keys.TextCommonBtnCancel},
							Data: response.CallbackData{
								Unique: keys.BtnOnboardCancel,
							},
						},
					},
				},
			},
		}, nil

	case keys.StateOnboardMedia:
		return uc.finish(ctx, u, text, false)

	default:
		return response.Reply{}, nil
	}
}

func (uc *Usecase) OnVideo(ctx context.Context, u *user.User, fileID string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateOnboardMedia {
		return response.Reply{}, nil
	}

	return uc.finish(ctx, u, fileID, true)
}

func (uc *Usecase) NoiseSelected(ctx context.Context, u *user.User, noise user.NoiseLevel) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateOnboardNoise {
		return response.Reply{}, errors.New("invalid state")
	}
	_ = uc.state.SetData(ctx, u.TelegramID, keys.DataOnboardNoiseLevel, noise, uc.ttl)
	_ = uc.state.SetState(ctx, u.TelegramID, keys.StateOnboardMedia, uc.ttl)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextOnboardStep3PromptMedia},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextOnboardBtnSkip},
						Data: response.CallbackData{
							Unique: keys.BtnOnboardSkip,
						},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnCancel},
						Data: response.CallbackData{
							Unique: keys.BtnOnboardCancel,
						},
					},
				},
			},
		},
	}, nil
}

func (uc *Usecase) SkipMedia(ctx context.Context, u *user.User) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	if st != keys.StateOnboardMedia {
		return response.Reply{}, errors.New("invalid state")
	}

	return uc.finish(ctx, u, "", false)
}

func (uc *Usecase) CancelFlow(ctx context.Context, u *user.User) (response.Reply, error) {
	_ = uc.state.ClearState(ctx, u.TelegramID)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataOnboardNoiseLevel)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataOnboardName)

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextOnboardMsgCancel},
	}, nil
}

func (uc *Usecase) finish(ctx context.Context, u *user.User, mediaData string, isVideo bool) (response.Reply, error) {
	var name, noiseLevel string
	_ = uc.state.GetData(ctx, u.TelegramID, keys.DataOnboardName, &name)
	_ = uc.state.GetData(ctx, u.TelegramID, keys.DataOnboardNoiseLevel, &noiseLevel)

	u.Name = name
	u.UpdatedAt = time.Now()

	if err := uc.users.Update(ctx, u); err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextOnboardMsgErr},
		}, err
	}

	if err := uc.users.SubmitApplication(ctx, u.ID); err != nil {
		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextOnboardMsgErr},
		}, nil
	}

	_ = uc.state.ClearState(ctx, u.TelegramID)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataOnboardName)
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataOnboardNoiseLevel)

	errApp := uc.admin.NewApplication(ctx, ApplicationPayload{
		TelegramUsername: strings.TrimSpace(u.Username),
		UserID:           u.ID,
		Name:             name,
		NoiseLevel:       noiseLevel,
		MediaData:        mediaData,
		IsVideo:          isVideo,
	})
	if errApp != nil {
		slog.Error("failed to send admin notification for new application", "err", errApp)
	}

	return response.Reply{
		Kind: response.KindSend,
		Text: response.Text{Key: keys.TextOnboardMsgSuccess},
	}, nil
}
