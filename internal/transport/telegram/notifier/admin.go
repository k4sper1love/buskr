package notifier

import (
	"context"
	"strings"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
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

	username := strings.TrimPrefix(strings.TrimSpace(payload.TelegramUsername), "@")

	noiseTranslated := n.tr.T(n.adminLang, "common.lbl.noise_"+payload.NoiseLevel, nil)
	caption := n.tr.T(n.adminLang, keys.TextAdminModAppTitle, map[string]any{
		"username": mdutil.Escape(username),
		"name":     mdutil.Escape(payload.Name),
		"noise":    noiseTranslated,
	})

	opts := &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: menu,
	}
	if n.adminThreadApplications > 0 {
		opts.ThreadID = n.adminThreadApplications
	}

	var err error
	if payload.IsVideo {
		video := &telebot.Video{
			File:    telebot.File{FileID: payload.MediaData},
			Caption: caption,
		}
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, video, opts)
	} else if payload.MediaData != "" {
		textMsg := n.tr.T(n.adminLang, keys.TextAdminModAppLink, map[string]any{
			"caption": caption,
			"link":    payload.MediaData,
		})
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, opts)
	} else {
		textMsg := n.tr.T(n.adminLang, keys.TextAdminModAppNoMedia, map[string]any{
			"caption": caption,
		})
		_, err = n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, opts)
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

	// Strip @ if present — the locale template already has "@{{.username}}"
	username := strings.TrimPrefix(strings.TrimSpace(payload.TelegramUsername), "@")

	caption := n.tr.T(n.adminLang, keys.TextAdminModNoiseTitle, map[string]any{
		"username":        mdutil.Escape(username),
		"name":            mdutil.Escape(payload.Name),
		"current_noise":   n.tr.T(n.adminLang, "common.lbl.noise_"+string(payload.CurrentNoise), nil),
		"requested_noise": n.tr.T(n.adminLang, "common.lbl.noise_"+string(payload.RequestedNoise), nil),
		"karma":           payload.Karma,
		"reason":          mdutil.Escape(payload.Reason),
	})

	opts := &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: menu,
	}
	if n.adminThreadUpgrades > 0 {
		opts.ThreadID = n.adminThreadUpgrades
	}

	_, err := n.bot.Send(&telebot.Chat{ID: n.adminChatID}, caption, opts)
	return err
}

func (n *Notifier) NewInviteRegistration(ctx context.Context, username, name, noiseLevel, invitedBy string) error {
	noiseTranslated := n.tr.T(n.adminLang, "common.lbl.noise_"+noiseLevel, nil)
	textMsg := n.tr.T(n.adminLang, keys.TextAdminModInviteRegTitle, map[string]any{
		"username":   mdutil.Escape(username),
		"name":       mdutil.Escape(name),
		"noise":      noiseTranslated,
		"invited_by": mdutil.Escape(invitedBy),
	})

	opts := &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	}
	if n.adminThreadApplications > 0 {
		opts.ThreadID = n.adminThreadApplications
	}

	_, err := n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, opts)
	return err
}

func (n *Notifier) NewLocationSuggestion(ctx context.Context, loc *location.Location, suggestedBy *user.User) error {
	menu := &telebot.ReplyMarkup{}

	approveText := n.tr.T(n.adminLang, keys.TextAdminModBtnApprove, nil)
	rejectText := n.tr.T(n.adminLang, keys.TextAdminModBtnReject, nil)

	btnApprove := menu.Data(approveText, keys.BtnAdminLocSuggestAppr, loc.ID)
	btnReject := menu.Data(rejectText, keys.BtnAdminLocSuggestRej, loc.ID)
	menu.Inline(menu.Row(btnApprove, btnReject))

	username := strings.TrimPrefix(strings.TrimSpace(suggestedBy.Username), "@")
	var suggestedByStr string
	if username != "" {
		suggestedByStr = mdutil.Escape(suggestedBy.Name) + " (@" + mdutil.Escape(username) + ")"
	} else {
		suggestedByStr = mdutil.Escape(suggestedBy.Name)
	}

	noiseTranslated := n.tr.T(n.adminLang, "common.lbl.noise_"+string(loc.MaxNoise), nil)
	textMsg := n.tr.T(n.adminLang, keys.TextAdminModLocSuggestTitle, map[string]any{
		"name":         mdutil.Escape(loc.Name),
		"desc":         mdutil.Escape(loc.Description),
		"noise":        noiseTranslated,
		"lat":          loc.Coords.Lat,
		"lon":          loc.Coords.Lon,
		"suggested_by": suggestedByStr,
	})

	opts := &telebot.SendOptions{
		ParseMode:   telebot.ModeHTML,
		ReplyMarkup: menu,
	}
	threadID := n.adminThreadSuggestions
	if threadID == 0 {
		threadID = n.adminThreadApplications
	}
	if threadID > 0 {
		opts.ThreadID = threadID
	}

	_, err := n.bot.Send(&telebot.Chat{ID: n.adminChatID}, textMsg, opts)
	return err
}
