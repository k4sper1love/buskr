package render

import (
	"bytes"
	"fmt"

	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type Translator interface {
	T(lang, key string, args map[string]any) string
}

type Renderer struct {
	tr          Translator
	defaultLang string
}

func NewRenderer(tr Translator, defaultLang string) *Renderer {
	return &Renderer{tr: tr, defaultLang: defaultLang}
}

func (r *Renderer) Render(c telebot.Context, rep response.Reply) error {
	lang := r.defaultLang
	if s := c.Sender(); s != nil && s.LanguageCode != "" {
		lang = s.LanguageCode
	}

	text := r.Translate(lang, rep.Text)

	var opts []any

	if rep.Keyboard.Remove {
		opts = append(opts, telebot.RemoveKeyboard)
	} else if len(rep.Keyboard.InlineRows) > 0 {
		// ← Inline-клавиатура
		menu := &telebot.ReplyMarkup{}
		var rows []telebot.Row
		for _, row := range rep.Keyboard.InlineRows {
			var btns []telebot.Btn
			for _, btn := range row {
				text := r.Translate(lang, btn.Text)
				switch {
				case btn.WebAppURL != "":
					btns = append(btns, menu.WebApp(text, &telebot.WebApp{URL: btn.WebAppURL}))
				case btn.URL != "":
					btns = append(btns, menu.URL(text, btn.URL))
				case btn.Data.Unique != "":
					btns = append(btns, menu.Data(text, btn.Data.Unique, btn.Data.Args...))
				default:
					btns = append(btns, menu.Text(text))
				}
			}
			if len(btns) > 0 {
				rows = append(rows, menu.Row(btns...))
			}
		}
		if len(rows) > 0 {
			menu.Inline(rows...)
			opts = append(opts, menu)
		}
	} else if len(rep.Keyboard.ReplyRows) > 0 {
		// ← Reply-клавиатура
		menu := &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
		var rows []telebot.Row

		for _, row := range rep.Keyboard.ReplyRows {
			var btns []telebot.Btn
			for _, t := range row {
				btns = append(btns, menu.Text(r.Translate(lang, t)))
			}
			if len(btns) > 0 {
				rows = append(rows, menu.Row(btns...))
			}
		}
		if len(rows) > 0 {
			menu.Reply(rows...)
			opts = append(opts, menu)
		}
	}

	// Always render with Markdown support
	opts = append(opts, telebot.ModeMarkdown)

	switch rep.Kind {
	case response.KindSend:
		return c.Send(text, opts...)
	case response.KindEdit:
		return c.Edit(text, opts...)
	case response.KindSendImage:
		photo := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(rep.Image)),
			Caption: text,
		}
		return c.Send(photo, opts...)
	default:
		return fmt.Errorf("unknown reply kind: %d", rep.Kind)
	}
}

func (r *Renderer) Translate(lang string, text response.Text) string {
	if lang == "" {
		lang = r.defaultLang
	}

	if text.Key == "" {
		if text.Fallback != "" {
			return text.Fallback
		}
		return ""
	}

	if r.tr == nil {
		if text.Fallback != "" {
			return text.Fallback
		}
		return text.Key
	}

	text = r.resolveSubKeys(lang, text)

	s := r.tr.T(lang, text.Key, text.Args)
	if s == "" && text.Fallback != "" {
		return text.Fallback
	}
	return s
}

// resolveSubKeys translates any Args values that are themselves i18n keys,
// as declared in text.SubKeyArgs. This enables nested translations.
func (r *Renderer) resolveSubKeys(lang string, text response.Text) response.Text {
	if len(text.SubKeyArgs) == 0 || len(text.Args) == 0 || r.tr == nil {
		return text
	}

	resolvedArgs := make(map[string]any, len(text.Args))
	for k, v := range text.Args {
		resolvedArgs[k] = v
	}

	for _, argKey := range text.SubKeyArgs {
		if val, ok := resolvedArgs[argKey].(string); ok && val != "" {
			resolvedArgs[argKey] = r.tr.T(lang, val, nil)
		}
	}

	text.Args = resolvedArgs
	return text
}
