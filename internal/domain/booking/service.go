package booking

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
)

const UnlimitedBookings = -1

type UserProvider interface {
	GetByID(ctx context.Context, id string) (*user.User, error)
}

type LocationProvider interface {
	GetByID(ctx context.Context, id string) (*location.Location, error)
}

type Config struct {
	MaxActiveBookings           int
	AdminMaxActiveBookings      int
	MaxBookingsPerLocation      int
	AdminMaxBookingsPerLocation int
	MaxAdvanceDays              int
	AdjacencyRadius             int
	AdminBypassNoiseLimits      bool
	AdminBypassUserOverlap      bool
	AdminBypassNoisyNeighbor    bool
}

type Service struct {
	repo  Repository
	users UserProvider
	locs  LocationProvider
	cfg   Config
}

func NewService(repo Repository,
	users UserProvider,
	locs LocationProvider,
	config Config,
) *Service {
	return &Service{
		repo:  repo,
		users: users,
		locs:  locs,
		cfg:   config,
	}
}

func (s *Service) GetScheduleForLocation(ctx context.Context, locID string, date time.Time) ([]*Booking, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()) // 00:00:00 of the day
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)                            // 23:59:59 of the day

	return s.repo.GetByLocationAndTime(ctx, locID, startOfDay, endOfDay)
}

func (s *Service) BookSlot(ctx context.Context, userID, locID string, start, end time.Time) (*Booking, error) {
	if start.Before(time.Now()) || end.Before(start) {
		return nil, ErrInvalidTime
	}

	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u.Status != user.StatusActive {
		return nil, ErrUserNotActive
	}

	loc, err := s.locs.GetByID(ctx, locID)
	if err != nil {
		return nil, err
	}
	if loc.Status != location.StatusActive {
		return nil, ErrLocationInactive
	}

	bypassNoise := u.Role == user.RoleAdmin && s.cfg.AdminBypassNoiseLimits
	if !bypassNoise && !IsNoiseCompatible(u.NoiseLevel, loc.MaxNoise) {
		return nil, ErrNoiseExceeded
	}

	if u.Role != user.RoleAdmin && start.After(time.Now().AddDate(0, 0, s.cfg.MaxAdvanceDays)) {
		return nil, ErrTooFarInFuture
	}

	activeBookings, err := s.repo.CountActiveFuture(ctx, userID)
	if err != nil {
		return nil, err
	}

	maxLimit := s.cfg.MaxActiveBookings
	if u.Role == user.RoleAdmin {
		maxLimit = s.cfg.AdminMaxActiveBookings
	}

	if maxLimit != UnlimitedBookings && activeBookings >= maxLimit {
		return nil, ErrMaxActiveBookings
	}

	bookingsPerLocation, err := s.repo.CountByUserAndLocationAndDay(ctx, userID, locID, start)
	if err != nil {
		return nil, err
	}

	maxPerDay := s.cfg.MaxBookingsPerLocation
	if u.Role == user.RoleAdmin {
		maxPerDay = s.cfg.AdminMaxBookingsPerLocation
	}

	if maxPerDay != UnlimitedBookings && bookingsPerLocation >= maxPerDay {
		return nil, ErrMaxBookingsPerLocation
	}

	bypassOverlap := u.Role == user.RoleAdmin && s.cfg.AdminBypassUserOverlap
	if !bypassOverlap {
		userOverlap, err := s.repo.HasOverlapByUser(ctx, userID, start, end)
		if err != nil {
			return nil, err
		}
		if userOverlap {
			return nil, ErrTimeOverlap
		}
	}

	locOverlap, err := s.repo.HasOverlapByLocation(ctx, locID, start, end)
	if err != nil {
		return nil, err
	}
	if locOverlap {
		return nil, ErrSlotTaken
	}

	bypassNoisy := u.Role == user.RoleAdmin && s.cfg.AdminBypassNoisyNeighbor
	if !bypassNoisy && u.NoiseLevel == user.NoiseHard {
		noisyNeighbor, err := s.repo.HasNoisyNeighbor(ctx, locID, start, end, float64(s.cfg.AdjacencyRadius))
		if err != nil {
			return nil, err
		}
		if noisyNeighbor {
			return nil, ErrNoisyNeighbor
		}
	}

	booking := NewBooking(userID, locID, start, end)
	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingID string, userID string) error {
	b, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if b.UserID != userID {
		return user.ErrAccessDenied
	}

	if b.Status != StatusPending {
		return ErrInvalidStatus
	}

	b.Status = StatusCancelled
	b.UpdatedAt = time.Now()

	return s.repo.Update(ctx, b)
}

func (s *Service) CheckIn(ctx context.Context, bookingID string, userID string) error {
	b, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if b.UserID != userID {
		return user.ErrAccessDenied
	}

	if b.Status == StatusActive {
		return ErrAlreadyCheckedIn
	}

	if b.Status != StatusPending {
		return ErrInvalidStatus
	}

	b.Status = StatusActive
	b.UpdatedAt = time.Now()

	return s.repo.Update(ctx, b)
}

func (s *Service) CompleteBooking(ctx context.Context, bookingID string) error {
	b, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if b.Status != StatusActive {
		return ErrInvalidStatus
	}

	b.Status = StatusCompleted
	b.UpdatedAt = time.Now()

	return s.repo.Update(ctx, b)
}

func (s *Service) MarkNoShow(ctx context.Context, bookingID string) error {
	b, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if b.Status != StatusPending {
		return ErrInvalidStatus
	}

	b.Status = StatusNoShow
	b.UpdatedAt = time.Now()

	return s.repo.Update(ctx, b)
}

func (s *Service) GetUpcomingForReminder(ctx context.Context, timeUntilStart time.Duration) ([]*Booking, error) {
	targetTime := time.Now().Add(timeUntilStart)
	return s.repo.GetUpcoming(ctx, time.Now(), targetTime)
}

func (s *Service) GetPendingForCheckinTimeout(ctx context.Context, timeout time.Duration) ([]*Booking, error) {
	deadline := time.Now().Add(-timeout)
	return s.repo.GetPendingBefore(ctx, deadline)
}

func (s *Service) GetActiveForCompletion(ctx context.Context) ([]*Booking, error) {
	return s.repo.GetActiveEndedBefore(ctx, time.Now())
}

func (s *Service) GetPendingOrActiveBookings(ctx context.Context, userID string) ([]*Booking, error) {
	return s.repo.GetPendingOrActiveByUser(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, bookingID string) (*Booking, error) {
	return s.repo.GetByID(ctx, bookingID)
}

func (s *Service) GetLastBookingByUser(ctx context.Context, userID string) (*Booking, error) {
	return s.repo.GetLastBookingByUser(ctx, userID)
}

func IsNoiseCompatible(userNoise user.NoiseLevel, locLimit location.NoiseLimit) bool {
	switch locLimit {
	case location.LimitLight:
		return userNoise == user.NoiseLight
	case location.LimitMedium:
		return userNoise == user.NoiseLight || userNoise == user.NoiseMedium
	case location.LimitHard:
		return userNoise == user.NoiseMedium || userNoise == user.NoiseHard
	default:
		return false
	}
}
