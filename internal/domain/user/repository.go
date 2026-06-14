package user

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	GetActiveUsers(ctx context.Context) ([]*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error

	GetStats(ctx context.Context, userID string) (*UserStats, error)
	UpdateStats(ctx context.Context, stats *UserStats) error

	CreateInvite(ctx context.Context, invite *Invite) error
	GetInviteByToken(ctx context.Context, token string) (*Invite, error)
	UpdateInvite(ctx context.Context, invite *Invite) error
	GetInviteByUsedByID(ctx context.Context, userID string) (*Invite, error)

	GetUsersPaginated(ctx context.Context, offset, limit int, sortBy string) ([]*User, int, error)
	FindByQuery(ctx context.Context, query string, offset, limit int) ([]*User, int, error)
}
