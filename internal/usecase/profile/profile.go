package profile

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
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
	Reason           string
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

func (uc *Usecase) SaveProfileMessageID(ctx context.Context, tgID int64, msgID int) error {
	return uc.state.SetData(ctx, tgID, keys.DataProfileMsgID, msgID, uc.ttl)
}

func (uc *Usecase) GetProfileMessageID(ctx context.Context, tgID int64) (int, error) {
	var msgID int
	err := uc.state.GetData(ctx, tgID, keys.DataProfileMsgID, &msgID)
	return msgID, err
}

func (uc *Usecase) ClearProfileMessageID(ctx context.Context, tgID int64) error {
	return uc.state.ClearData(ctx, tgID, keys.DataProfileMsgID)
}
