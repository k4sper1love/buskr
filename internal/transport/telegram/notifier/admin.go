package notifier

import (
	"context"
	"strings"

	"github.com/k4sper1love/buskr/internal/mdutil"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/onboarding"
	"github.com/k4sper1love/buskr/internal/usecase/profile"
	"gopkg.in/telebot.v3"
)

func (n *Notifier) NewApplication(ctx context.Context, payload onboarding.ApplicationPayload) error {
	menu := &telebot.ReplyMarkup{}

	approveText := n.tr.T(n.adminLang, keys.TextAdminModBtnApprove, nil)
	rejectText := n.tr.T(n.adminLang, keys.TextAdminModBtnReject, nil)

	btnApprove := menu.Data(approveText, keys.BtnAdminAppAppr, payload.UserID, payload.NoiseLevel)
	btnReject := menu.Data(rejectText, keys.BtnAdminAppRej, payload.UserID)
	menu.Inline(menu.Row(btnApprove, btnReject))

	username := strings.TrimSpace(payload.TelegramUsername)
	if username != "" && !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	noiseTranslated := n.tr.T(n.adminLang, "common.lbl.noise_"+payload.NoiseLevel, nil)
	caption := n.tr.T(n.adminLang, keys.TextAdminModAppTitle, map[string]any{
		"username": mdutil.Escape(username),
		"name":     mdutil.Escape(payload.Name),
		"noise":    noiseTranslated,
	})

	var err error
	if payload.IsVideo {
		video := &telebot.Video{
			File:    telebot.File{FileID: payload.MediaData},
			Caption: caption,
		}
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, video, menu, telebot.ModeMarkdown)
	} else if payload.MediaData != "" {
		textMsg := n.tr.T(n.adminLang, keys.TextAdminModAppLink, map[string]any{
			"caption": caption,
			"link":    payload.MediaData,
		})
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, menu, telebot.ModeMarkdown)
	} else {
		textMsg := n.tr.T(n.adminLang, keys.TextAdminModAppNoMedia, map[string]any{
			"caption": caption,
		})
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, menu, telebot.ModeMarkdown)
	}

	return err
}

func (n *Notifier) NewNoiseUpgrade(ctx context.Context, payload profile.NoiseUpgradePayload) error {
	menu := &telebot.ReplyMarkup{}

	approveText := n.tr.T(n.adminLang, keys.TextAdminModBtnApprove, nil)
	rejectText := n.tr.T(n.adminLang, keys.TextAdminModBtnReject, nil)

	btnApprove := menu.Data(approveText, keys.BtnAdminNoiseAppr, payload.UserID, string(payload.RequestedNoise))
	btnReject := menu.Data(rejectText, keys.BtnAdminNoiseRej, payload.UserID)
	menu.Inline(menu.Row(btnApprove, btnReject))

	username := strings.TrimSpace(payload.TelegramUsername)
	if username != "" && !strings.HasPrefix(username, "@") {
		username = "@" + username
	}

	caption := n.tr.T(n.adminLang, keys.TextAdminModNoiseTitle, map[string]any{
		"username":        mdutil.Escape(username),
		"name":            mdutil.Escape(payload.Name),
		"current_noise":   n.tr.T(n.adminLang, "common.lbl.noise_"+string(payload.CurrentNoise), nil),
		"requested_noise": n.tr.T(n.adminLang, "common.lbl.noise_"+string(payload.RequestedNoise), nil),
		"karma":           payload.Karma,
	})

	_, err := n.bot.Send(&telebot.Chat{ID: n.adminChatID}, caption, menu, telebot.ModeMarkdown)
	return err
}
