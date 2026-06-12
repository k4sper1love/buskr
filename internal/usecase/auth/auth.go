package auth

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
)


type StateStore interface {
	SetState(ctx context.Context, userID int64, state string, ttl time.Duration) error
	GetState(ctx context.Context, userID int64) (string, error)
	ClearState(ctx context.Context, userID int64) error

	SetData(ctx context.Context, userID int64, key string, data any, ttl time.Duration) error
	GetData(ctx context.Context, userID int64, key string, dest any) error
	ClearData(ctx context.Context, userID int64, key string) error
}

type Users interface {
	ProcessInvite(ctx context.Context, userID string, payload string) error
	Update(ctx context.Context, user *user.User) error
}

type Usecase struct {
	state   StateStore
	users   Users
	ttl     time.Duration
	botName string
}

func NewUsecase(state StateStore, users Users, ttl time.Duration, botName string) *Usecase {
	return &Usecase{
		state:   state,
		users:   users,
		ttl:     ttl,
		botName: botName,
	}
}
