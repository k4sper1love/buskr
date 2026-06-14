package notifier

import (
	"gopkg.in/telebot.v3"
)


type Translator interface {
	T(lang, key string, args map[string]any) string
}

type Notifier struct {
	bot                     *telebot.Bot
	tr                      Translator
	adminChatID             int64
	adminLang               string
	adminThreadApplications int
	adminThreadUpgrades     int
}

func NewNotifier(
	bot *telebot.Bot,
	tr Translator,
	adminChatID int64,
	adminLang string,
	adminThreadApplications int,
	adminThreadUpgrades int,
) *Notifier {
	return &Notifier{
		bot:                     bot,
		tr:                      tr,
		adminChatID:             adminChatID,
		adminLang:               adminLang,
		adminThreadApplications: adminThreadApplications,
		adminThreadUpgrades:     adminThreadUpgrades,
	}
}
