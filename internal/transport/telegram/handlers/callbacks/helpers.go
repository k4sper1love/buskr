package callbacks

import (
	"context"
	"fmt"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/transport/telegram/ctxkey"
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
