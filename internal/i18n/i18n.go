package i18n

import (
	"embed"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

type Translator struct {
	bundle *i18n.Bundle
}

func NewTranslator() *Translator {
	bundle := i18n.NewBundle(language.Russian)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFileFS(localesFS, "locales/ru.json")
	if err != nil {
		slog.Error("failed to load ru.json", "err", err)
		os.Exit(1)
	}

	_, err = bundle.LoadMessageFileFS(localesFS, "locales/en.json")
	if err != nil {
		slog.Error("failed to load en.json", "err", err)
		os.Exit(1)
	}

	return &Translator{bundle: bundle}
}

func (t *Translator) T(lang, key string, args map[string]any) string {
	localizer := i18n.NewLocalizer(t.bundle, lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	})
	if err != nil {
		return key
	}

	return msg
}
