package adminuser

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
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetStats(ctx context.Context, userID string) (*user.UserStats, error)
	Update(ctx context.Context, user *user.User) error
	GetUsersPaginated(ctx context.Context, offset, limit int, sortBy string) ([]*user.User, int, error)
	FindByQuery(ctx context.Context, query string, offset, limit int) ([]*user.User, int, error)
	IsLowKarma(u *user.User) bool
}

type Usecase struct {
	state StateStore
	users Users
	ttl   time.Duration
}

func NewUsecase(state StateStore, users Users, ttl time.Duration) *Usecase {
	return &Usecase{
		state: state,
		users: users,
		ttl:   ttl,
	}
}
