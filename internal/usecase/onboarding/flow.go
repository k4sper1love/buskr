package onboarding

import (
	"context"
	"strings"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
)

func (uc *Usecase) OnText(ctx context.Context, u *user.User, text string) (response.Reply, error) {
	st, err := uc.state.GetState(ctx, u.TelegramID)
	if err != nil {
		return response.Reply{}, err
	}

	switch st {
	case keys.StateOnboardName:
		_ = uc.state.SetData(ctx, u.TelegramID, keys.DataOnboardName, text, uc.ttl)
		_ = uc.state.SetState(ctx, u.TelegramID, keys.StateOnboardFormat, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextOnboardStep2PromptFormat},
			Keyboard: response.Keyboard{
				ReplyRows: [][]response.Text{
					{
						{Key: keys.TextOnboardBtnSolo},
						{Key: keys.TextOnboardBtnGroup},
					},
				},
			},
		}, nil

	case keys.StateOnboardFormat:
		category := user.NoiseLight
		if strings.Contains(text, "Группа") {
			category = user.NoiseHard
		}

		_ = uc.state.SetData(ctx, u.TelegramID, keys.DataOnboardCat, category, uc.ttl)
		_ = uc.state.SetState(ctx, u.TelegramID, keys.StateOnboardMedia, uc.ttl)

		return response.Reply{
			Kind: response.KindSend,
			Text: response.Text{Key: keys.TextOnboardStep3PromptMedia},
			Keyboard: response.Keyboard{
				Remove: true,
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

func (uc *Usecase) finish(ctx context.Context, u *user.User, mediaData string, isVideo bool) (response.Reply, error) {
	var name, category string
	_ = uc.state.GetData(ctx, u.TelegramID, keys.DataOnboardName, &name)
	_ = uc.state.GetData(ctx, u.TelegramID, keys.DataOnboardCat, &category)

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
	_ = uc.state.ClearData(ctx, u.TelegramID, keys.DataOnboardCat)

	_ = uc.admin.NewApplication(ctx, ApplicationPayload{
		TelegramUsername: strings.TrimSpace(u.Username),
		UserID:           u.ID,
		Name:             name,
		Category:         category,
		MediaData:        mediaData,
		IsVideo:          isVideo,
	})

	return response.Reply{
		Kind: response.KindSend,
		Text: response.Text{Key: keys.TextOnboardMsgSuccess},
	}, nil
}
