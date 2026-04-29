package admin

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/response"
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
	ApproveApplication(ctx context.Context, userID string, category user.NoiseLevel) error
	RejectApplication(ctx context.Context, userID string) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	GenerateInvite(ctx context.Context, noiseLevel user.NoiseLevel, expirationHours int) (*user.Invite, error)
}

type UserNotification struct {
	TelegramID int64
	Reply      response.Reply
}

type ModerationResult struct {
	AdminEditSuffix response.Text
	NotifyUser      *UserNotification
	CallbackText    response.Text
}

type Usecase struct {
	state     StateStore
	users     Users
	ttl       time.Duration
	inviteTTL int
}

func NewUsecase(state StateStore, users Users, ttl time.Duration, inviteTTL int) *Usecase {
	return &Usecase{
		state:     state,
		users:     users,
		ttl:       ttl,
		inviteTTL: inviteTTL,
	}
}
