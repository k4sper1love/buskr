package handlers

import (
	"fmt"
	"log/slog"

	"github.com/k4sper1love/buskr/internal/config"
	"github.com/k4sper1love/buskr/internal/mdutil"
	"github.com/k4sper1love/buskr/internal/transport/telegram/render"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
	"gopkg.in/telebot.v3"
)

type GroupHandler struct {
	cfg      *config.TelegramConfig
	renderer *render.Renderer
}

func NewGroupHandler(cfg *config.TelegramConfig, renderer *render.Renderer) *GroupHandler {
	return &GroupHandler{
		cfg:      cfg,
		renderer: renderer,
	}
}

func (h *GroupHandler) HandleAddedToGroup(c telebot.Context) error {
	chat := c.Chat()
	slog.Debug("bot added to group", "chat_id", chat.ID, "chat_title", chat.Title)
	if chat.ID != h.cfg.AdminChatID && chat.ID != h.cfg.PublicChatID {
		slog.Warn("unauthorized chat, leaving", "chat_id", chat.ID)
		msg := h.renderer.Translate(h.cfg.DefaultLang, response.Text{
			Key: keys.TextGroupUnathorizedChat,
		})
		_ = c.Send(msg)
		return c.Bot().Leave(chat)
	}
	return nil
}

func (h *GroupHandler) HandleUserJoined(c telebot.Context) error {
	slog.Debug("user joined event", "chat_id", c.Chat().ID, "chat_title", c.Chat().Title)

	if c.Chat().ID == h.cfg.AdminChatID {
		return nil
	}

	if c.Chat().ID != h.cfg.PublicChatID {
		slog.Warn("user joined unauthorized chat, leaving", "chat_id", c.Chat().ID)
		_ = c.Bot().Leave(c.Chat())
		return nil
	}

	joinedUser := c.Message().UserJoined
	usersJoined := c.Message().UsersJoined
	slog.Debug("user joined details", "singular", joinedUser != nil, "plural_count", len(usersJoined))

	var targetUser *telebot.User
	if joinedUser != nil {
		targetUser = joinedUser
	} else if len(usersJoined) > 0 {
		targetUser = &usersJoined[0]
	}

	if targetUser == nil || targetUser.ID == c.Bot().Me.ID {
		slog.Debug("skipping welcome: target is nil or bot itself")
		return nil
	}

	name := targetUser.FirstName
	if targetUser.LastName != "" {
		name += " " + targetUser.LastName
	}

	escapedName := mdutil.Escape(name)
	mention := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, targetUser.ID, escapedName)

	welcomeText := response.Text{
		Key: keys.TextGroupWelcomeUser,
		Args: map[string]any{
			"mention":  mention,
			"bot_link": "https://t.me/" + c.Bot().Me.Username,
		},
		NoEscapeArgs: []string{"mention", "bot_link"},
	}

	msg := h.renderer.Translate(h.cfg.DefaultLang, welcomeText)
	return c.Send(msg, telebot.ModeHTML)
}
