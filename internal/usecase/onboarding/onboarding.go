package onboarding

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
	SubmitApplication(ctx context.Context, userID string) error
	Update(ctx context.Context, user *user.User) error
}

type AdminNotifier interface {
	NewApplication(ctx context.Context, payload ApplicationPayload) error
}

type ApplicationPayload struct {
	TelegramUsername string
	UserID           string
	Name             string
	Category         string
	MediaData        string
	IsVideo          bool
}

type Usecase struct {
	state StateStore
	users Users
	admin AdminNotifier
	ttl   time.Duration
}

func NewUsecase(state StateStore, users Users, admin AdminNotifier, ttl time.Duration) *Usecase {
	return &Usecase{
		state: state,
		users: users,
		admin: admin,
		ttl:   ttl,
	}
}
