package callbacks

import (
	"context"
	"fmt"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type callbackCtx struct {
	ctx  context.Context
	user *user.User
}

func extractCtx(c telebot.Context) (callbackCtx, error) {
	ctx, ok := c.Get("ctx").(context.Context)
	if !ok {
		return callbackCtx{}, fmt.Errorf("ctx not found in telebot context")
	}

	u, err := ctxkey.GetUser(c)
	if err != nil {
		return callbackCtx{}, err
	}

	return callbackCtx{ctx: ctx, user: u}, nil
}

func sendReplyToUser(bot *telebot.Bot, r *render.Renderer, tgID int64, rep response.Reply) {
	text := r.Translate("", rep.Text)
	var opts []any
	if rep.Keyboard.Remove {
		opts = append(opts, telebot.RemoveKeyboard)

	} else if len(rep.Keyboard.InlineRows) > 0 {
		menu := &telebot.ReplyMarkup{}
		var rows []telebot.Row
		for _, row := range rep.Keyboard.InlineRows {
			var btns []telebot.Btn
			for _, btn := range row {
				t := r.Translate("", btn.Text)
				btns = append(btns, menu.Data(t, btn.Data.Unique, btn.Data.Args...))
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
		menu := &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true}
		var rows []telebot.Row
		for _, row := range rep.Keyboard.ReplyRows {
			var btns []telebot.Btn
			for _, t := range row {
				btns = append(btns, menu.Text(r.Translate("", t)))
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

	opts = append(opts, telebot.ModeMarkdown)
	_, _ = bot.Send(&telebot.User{ID: tgID}, text, opts...)
}
