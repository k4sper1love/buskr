package booking

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
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
	SuggestLocation(ctx context.Context, name, desc string, lat, lon float64, noise location.NoiseLimit, suggestedBy string) (*location.Location, error)
	FindNearby(ctx context.Context, lat, lon float64, radiusMeters float64) ([]*location.LocationWithDist, error)
}

type Users interface {
	GetByID(ctx context.Context, id string) (*user.User, error)
	RecordNewBooking(ctx context.Context, userID string) error
	RecordSuccessfulCheckin(ctx context.Context, userID string) error
}

type Notifier interface {
	NewLocationSuggestion(ctx context.Context, loc *location.Location, suggestedBy *user.User) error
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
	state                 StateStore
	locs                  Locations
	bookings              Bookings
	users                 Users
	notifier              Notifier
	ttl                   time.Duration
	maxAdvanceDays        int
	webAppURL             string
	botUsername           string
	minHour               int
	maxHourWeekday        int
	maxHourWeekend        int
	maxDurationHours      int
	adminMaxDurationHours int
}

type webAppLocation struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	MaxNoise  string  `json:"max_noise,omitempty"`
	Allowed   *bool   `json:"allowed,omitempty"`
	IsVeteran bool    `json:"is_veteran"`
}

func NewUsecase(
	state StateStore,
	locs Locations,
	bookings Bookings,
	users Users,
	notifier Notifier,
	ttl time.Duration,
	maxAdvanceDays int,
	webAppURL string,
	botUsername string,
	minHour int,
	maxHourWeekday int,
	maxHourWeekend int,
	maxDurationHours int,
	adminMaxDurationHours int,
) *Usecase {
	return &Usecase{
		state:                 state,
		locs:                  locs,
		bookings:              bookings,
		users:                 users,
		notifier:              notifier,
		ttl:                   ttl,
		maxAdvanceDays:        maxAdvanceDays,
		webAppURL:             webAppURL,
		botUsername:           botUsername,
		minHour:               minHour,
		maxHourWeekday:        maxHourWeekday,
		maxHourWeekend:        maxHourWeekend,
		maxDurationHours:      maxDurationHours,
		adminMaxDurationHours: adminMaxDurationHours,
	}
}

func (uc *Usecase) SaveBookingMessageID(ctx context.Context, tgID int64, msgID int) error {
	return uc.state.SetData(ctx, tgID, keys.DataBookingMsgID, msgID, uc.ttl)
}

func (uc *Usecase) GetBookingMessageID(ctx context.Context, tgID int64) (int, error) {
	var msgID int
	err := uc.state.GetData(ctx, tgID, keys.DataBookingMsgID, &msgID)
	return msgID, err
}

func (uc *Usecase) ClearBookingMessageID(ctx context.Context, tgID int64) error {
	return uc.state.ClearData(ctx, tgID, keys.DataBookingMsgID)
}

func (uc *Usecase) SaveSuggestLocMessageID(ctx context.Context, tgID int64, msgID int) error {
	return uc.state.SetData(ctx, tgID, keys.DataSuggestLocMsgID, msgID, uc.ttl)
}

func (uc *Usecase) GetSuggestLocMessageID(ctx context.Context, tgID int64) (int, error) {
	var msgID int
	err := uc.state.GetData(ctx, tgID, keys.DataSuggestLocMsgID, &msgID)
	return msgID, err
}

func (uc *Usecase) ClearSuggestLocMessageID(ctx context.Context, tgID int64) error {
	return uc.state.ClearData(ctx, tgID, keys.DataSuggestLocMsgID)
}
