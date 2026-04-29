package notifier

import (
	"gopkg.in/telebot.v3"
)


type Translator interface {
	T(lang, key string, args map[string]any) string
}

type Notifier struct {
	bot         *telebot.Bot
	tr          Translator
	adminChatID int64
	adminLang   string
}

func NewNotifier(bot *telebot.Bot, tr Translator, adminChatID int64, adminLang string) *Notifier {
	return &Notifier{bot: bot, tr: tr, adminChatID: adminChatID, adminLang: adminLang}
}
