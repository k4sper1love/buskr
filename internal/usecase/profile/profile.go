package profile

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
	Update(ctx context.Context, user *user.User) error
	GetStats(ctx context.Context, userID string) (*user.UserStats, error)
}

type AdminNotifier interface {
	NewNoiseUpgrade(ctx context.Context, payload NoiseUpgradePayload) error
}

type NoiseUpgradePayload struct {
	TelegramUsername string
	UserID           string
	Name             string
	CurrentNoise     user.NoiseLevel
	RequestedNoise   user.NoiseLevel
	Karma            int
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
