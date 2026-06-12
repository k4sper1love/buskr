package admin

import (
	"context"
	"errors"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) ApproveApplication(ctx context.Context, actor *user.User, targetUserID string, category user.NoiseLevel) (ModerationResult, error) {
	if actor.Role != user.RoleAdmin {
		return ModerationResult{}, errors.New("access denied")
	}

	err := uc.users.ApproveApplication(ctx, targetUserID, category)
	if err != nil {
		return ModerationResult{}, err
	}

	targetUser, err := uc.users.GetByID(ctx, targetUserID)
	if err != nil {
		targetUser = nil // approved, but cant prepare notification
	}

	res := ModerationResult{
		AdminEditSuffix: response.Text{
			Key: keys.TextAdminModMsgApprSfx,
			Args: map[string]any{
				"actor": actor.Username,
			},
		},
		CallbackText: response.Text{
			Key: keys.TextAdminModMsgApprCb,
		},
	}

	if targetUser != nil {
		res.NotifyUser = &UserNotification{
			TelegramID: targetUser.TelegramID,
			Reply: response.Reply{
				Kind: response.KindSend,
				Text: response.Text{
					Key: keys.TextAdminModMsgApprNotify,
				},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{{
							Text: response.Text{Key: keys.TextCommonBtnMenu},
							Data: response.CallbackData{Unique: keys.BtnCommonMenu},
						}},
					},
				},
			},
		}
	}

	return res, nil
}

func (uc *Usecase) RejectApplication(ctx context.Context, actor *user.User, targetUserID string) (ModerationResult, error) {
	if actor.Role != user.RoleAdmin {
		return ModerationResult{}, errors.New("access denied")
	}

	err := uc.users.RejectApplication(ctx, targetUserID)
	if err != nil {
		return ModerationResult{}, err
	}

	targetUser, err := uc.users.GetByID(ctx, targetUserID)
	if err != nil {
		targetUser = nil // rejected, but cant prepare notification
	}

	res := ModerationResult{
		AdminEditSuffix: response.Text{
			Key: keys.TextAdminModMsgRejSfx,
			Args: map[string]any{
				"actor": actor.Username,
			},
		},
		CallbackText: response.Text{
			Key: keys.TextAdminModMsgRejCb,
		},
	}

	if targetUser != nil {
		res.NotifyUser = &UserNotification{
			TelegramID: targetUser.TelegramID,
			Reply: response.Reply{
				Kind: response.KindSend,
				Text: response.Text{
					Key: keys.TextAdminModMsgRejNotify,
				},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{{
							Text: response.Text{Key: keys.TextCommonBtnMenu},
							Data: response.CallbackData{Unique: keys.BtnCommonMenu},
						}},
					},
				},
			},
		}
	}

	return res, nil
}

func (uc *Usecase) ApproveNoiseUpgrade(ctx context.Context, actor *user.User, targetUserID string, category user.NoiseLevel) (ModerationResult, error) {
	if actor.Role != user.RoleAdmin {
		return ModerationResult{}, errors.New("access denied")
	}

	targetUser, err := uc.users.GetByID(ctx, targetUserID)
	if err != nil {
		return ModerationResult{}, err // not found
	}

	targetUser.NoiseLevel = category
	targetUser.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, targetUser)
	if err != nil {
		return ModerationResult{}, err // update failed
	}

	_ = uc.state.ClearData(ctx, targetUser.TelegramID, "noise_upgrade_pending")

	res := ModerationResult{
		AdminEditSuffix: response.Text{
			Key: keys.TextAdminModMsgUpgSfx,
			Args: map[string]any{
				"actor":    actor.Username,
				"category": "common.lbl.noise_" + string(category),
			},
			SubKeyArgs: []string{"category"},
		},
		CallbackText: response.Text{
			Key: keys.TextAdminModMsgUpgCb,
		},
		NotifyUser: &UserNotification{
			TelegramID: targetUser.TelegramID,
			Reply: response.Reply{
				Kind: response.KindSend,
				Text: response.Text{
					Key: keys.TextAdminModMsgUpgNotify,
					Args: map[string]any{
						"actor":    actor.Username,
						"category": "common.lbl.noise_" + string(category),
					},
					SubKeyArgs: []string{"category"},
				},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{{
							Text: response.Text{Key: keys.TextCommonBtnMenu},
							Data: response.CallbackData{Unique: keys.BtnCommonMenu},
						}},
					},
				},
			},
		},
	}

	return res, nil
}

func (uc *Usecase) RejectNoiseUpgrade(ctx context.Context, actor *user.User, targetUserID string) (ModerationResult, error) {
	if actor.Role != user.RoleAdmin {
		return ModerationResult{}, errors.New("access denied")
	}

	targetUser, err := uc.users.GetByID(ctx, targetUserID)
	if err != nil {
		return ModerationResult{}, err // not found
	}

	_ = uc.state.ClearData(ctx, targetUser.TelegramID, "noise_upgrade_pending")

	res := ModerationResult{
		AdminEditSuffix: response.Text{
			Key: keys.TextAdminModMsgRejSfx,
			Args: map[string]any{
				"actor": actor.Username,
			},
		},
		CallbackText: response.Text{
			Key: keys.TextAdminModMsgRejCb,
		},
		NotifyUser: &UserNotification{
			TelegramID: targetUser.TelegramID,
			Reply: response.Reply{
				Kind: response.KindSend,
				Text: response.Text{
					Key: keys.TextAdminModMsgRejNotify,
				},
				Keyboard: response.Keyboard{
					InlineRows: [][]response.Button{
						{{
							Text: response.Text{Key: keys.TextCommonBtnMenu},
							Data: response.CallbackData{Unique: keys.BtnCommonMenu},
						}},
					},
				},
			},
		},
	}

	return res, nil
}
