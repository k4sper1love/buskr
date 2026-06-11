package booking

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
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

type Locations interface {
	GetLocationsForMusicians(ctx context.Context) ([]*location.Location, error)
	GetByID(ctx context.Context, id string) (*location.Location, error)
}

type Bookings interface {
	GetScheduleForLocation(ctx context.Context, locID string, date time.Time) ([]*booking.Booking, error)
	BookSlot(ctx context.Context, userID, locID string, start, end time.Time) (*booking.Booking, error)
	GetPendingOrActiveBookings(ctx context.Context, userID string) ([]*booking.Booking, error)
	GetByID(ctx context.Context, id string) (*booking.Booking, error)
	CancelBooking(ctx context.Context, id, userID string) error
	CheckIn(ctx context.Context, id, userID string) error
	GetLastBookingByUser(ctx context.Context, userID string) (*booking.Booking, error)
}

type SendLocation struct {
	Latitude  float64
	Longitude float64
}

type BookingResult struct {
	Reply    response.Reply
	Location SendLocation
	Callback response.Text
}

type CheckInResult struct {
	SuccessSuffix response.Text
	Callback      response.Text
}

type Usecase struct {
	state          StateStore
	locs           Locations
	bookings       Bookings
	ttl            time.Duration
	maxAdvanceDays int
	webAppURL      string
	botUsername    string
}

func NewUsecase(state StateStore, locs Locations, bookings Bookings, ttl time.Duration, maxAdvanceDays int, webAppURL string, botUsername string) *Usecase {
	return &Usecase{
		state:          state,
		locs:           locs,
		bookings:       bookings,
		ttl:            ttl,
		maxAdvanceDays: maxAdvanceDays,
		webAppURL:      webAppURL,
		botUsername:    botUsername,
	}
}
