package middleware

import (
	"context"
	"log/slog"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
	"gopkg.in/telebot.v3"
)

func AuthMiddleware(users *user.Service) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			sender := c.Sender()
			if sender == nil {
				return next(c) // skip if there's no sender (e.g., some channel events)
			}

			// Do not register/load users for group messages (only allow private chats and callback clicks in groups)
			chat := c.Chat()
			if chat != nil && chat.Type != telebot.ChatPrivate && c.Callback() == nil {
				return next(c)
			}

			u, err := users.GetOrCreateUser(c.Get("ctx").(context.Context), sender.ID, sender.Username)
			if err != nil {
				slog.Error("failed to get or create user", "telegram_id", sender.ID, "err", err)
				return c.Send("System error occurred. Please try again later.")
			}

			if u.Status == user.StatusBanned {
				return c.Send("Your account has been banned due to policy violations.")
			}

			ctxkey.SetUser(c, u)
			return next(c)
		}
	}
}
