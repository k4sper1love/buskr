package handlers

import (
	"fmt"
	"log"

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
	log.Printf("HandleAddedToGroup triggered in chat: %d (%s)", chat.ID, chat.Title)
	if chat.ID != h.cfg.AdminChatID && chat.ID != h.cfg.PublicChatID {
		log.Printf("Chat %d is unauthorized. Leaving...", chat.ID)
		msg := h.renderer.Translate(h.cfg.DefaultLang, response.Text{
			Key: keys.TextGroupUnathorizedChat,
		})
		_ = c.Send(msg)
		return c.Bot().Leave(chat)
	}
	return nil
}

func (h *GroupHandler) HandleUserJoined(c telebot.Context) error {
	log.Printf("HandleUserJoined triggered in chat: %d (%s)", c.Chat().ID, c.Chat().Title)
	
	if c.Chat().ID == h.cfg.AdminChatID {
		return nil
	}

	if c.Chat().ID != h.cfg.PublicChatID {
		log.Printf("User joined unauthorized chat: %d. Leaving...", c.Chat().ID)
		_ = c.Bot().Leave(c.Chat())
		return nil
	}

	joinedUser := c.Message().UserJoined
	usersJoined := c.Message().UsersJoined
	log.Printf("UserJoined (singular): %+v, UsersJoined (plural) count: %d", joinedUser, len(usersJoined))

	var targetUser *telebot.User
	if joinedUser != nil {
		targetUser = joinedUser
	} else if len(usersJoined) > 0 {
		targetUser = &usersJoined[0]
	}

	if targetUser == nil || targetUser.ID == c.Bot().Me.ID {
		log.Printf("Target user is nil or is the bot itself. Skipping welcome.")
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
