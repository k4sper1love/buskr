package adminloc

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
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

type Locations interface {
	CreateLocation(ctx context.Context, name, desc string, lat, lon float64, noise location.NoiseLimit) (*location.Location, error)
	UpdateLocation(ctx context.Context, id string, name, desc string, lat, lon float64, noise location.NoiseLimit) error
	GetLocationsForAdmin(ctx context.Context) ([]*location.Location, error)
	GetByID(ctx context.Context, id string) (*location.Location, error)
	ChangeStatus(ctx context.Context, id string, status location.Status) error
	FindNearby(ctx context.Context, lat, lon float64, radiusMeters float64) ([]*location.LocationWithDist, error)
	DeleteLocation(ctx context.Context, id string) error
}

type Bookings interface {
	GetScheduleForLocation(ctx context.Context, locID string, date time.Time) ([]*booking.Booking, error)
}

type Users interface {
	GetByID(ctx context.Context, id string) (*user.User, error)
}

type MapGenerator interface {
	Generate(locs []*location.Location) ([]byte, string, error)
}

type Usecase struct {
	state          StateStore
	locs           Locations
	bookings       Bookings
	users          Users
	maps           MapGenerator
	ttl            time.Duration
	maxAdvanceDays int
	webAppURL      string
	botUsername    string
}

func NewUsecase(state StateStore, locs Locations, bookings Bookings, users Users, maps MapGenerator, ttl time.Duration, maxAdvanceDays int, webAppURL string, botUsername string) *Usecase {
	return &Usecase{
		state:          state,
		locs:           locs,
		bookings:       bookings,
		users:          users,
		maps:           maps,
		ttl:            ttl,
		maxAdvanceDays: maxAdvanceDays,
		webAppURL:      webAppURL,
		botUsername:    botUsername,
	}
}
