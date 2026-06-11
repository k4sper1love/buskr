package booking

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, b *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	HasOverlapByLocation(ctx context.Context, locID string, start, end time.Time) (bool, error)
	HasOverlapByUser(ctx context.Context, userID string, start, end time.Time) (bool, error)
	Update(ctx context.Context, b *Booking) error
	GetByLocationAndTime(ctx context.Context, locID string, start, end time.Time) ([]*Booking, error)
	GetUpcoming(ctx context.Context, from, to time.Time) ([]*Booking, error)
	GetPendingBefore(ctx context.Context, deadline time.Time) ([]*Booking, error)
	GetActiveEndedBefore(ctx context.Context, deadline time.Time) ([]*Booking, error)
	GetPendingOrActiveByUser(ctx context.Context, userID string) ([]*Booking, error)
	CountActiveFuture(ctx context.Context, userID string) (int, error)
	CountByUserAndLocationAndDay(ctx context.Context, userID, locID string, date time.Time) (int, error)
	HasNoisyNeighbor(ctx context.Context, locID string, start, end time.Time, radiusMeters float64) (bool, error)
	GetLastBookingByUser(ctx context.Context, userID string) (*Booking, error)
}

