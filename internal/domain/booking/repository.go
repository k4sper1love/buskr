package booking

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, b *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	HasOverlap(ctx context.Context, locID string, start, end time.Time) (bool, error)
	Update(ctx context.Context, b *Booking) error
	GetByLocationAndTime(ctx context.Context, locID string, start, end time.Time) ([]*Booking, error)
	GetUpcoming(ctx context.Context, from, to time.Time) ([]*Booking, error)
	GetPendingBefore(ctx context.Context, deadline time.Time) ([]*Booking, error)
	GetActiveEndedBefore(ctx context.Context, deadline time.Time) ([]*Booking, error)
	GetPendingOrActiveByUser(ctx context.Context, userID string) ([]*Booking, error)
}
