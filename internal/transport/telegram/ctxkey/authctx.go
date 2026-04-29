package ctxkey

import (
	"errors"
	"fmt"

	"github.com/k4sper1love/buskr/internal/domain/user"

	"gopkg.in/telebot.v3"
)

type ctxKey int

const (
	userKey ctxKey = iota
)

func SetUser(c telebot.Context, user *user.User) {
	c.Set(fmt.Sprintf("%d", userKey), user)
}

func GetUser(c telebot.Context) (*user.User, error) {
	user, ok := c.Get(fmt.Sprintf("%d", userKey)).(*user.User)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}

	return user, nil
}
