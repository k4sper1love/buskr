package booking

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/location"
	"github.com/k4sper1love/buskr/internal/domain/user"
)

type UserProvider interface {
	GetByID(ctx context.Context, id string) (*user.User, error)
}

type LocationProvider interface {
	GetByID(ctx context.Context, id string) (*location.Location, error)
}

type Service struct {
	repo  Repository
	users UserProvider
	locs  LocationProvider
}

func NewService(repo Repository, users UserProvider, locs LocationProvider) *Service {
	return &Service{repo: repo, users: users, locs: locs}
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

	if !isNoiseAllowed(u.NoiseLevel, loc.MaxNoise) {
		return nil, ErrNoiseExceeded
	}

	overlap, err := s.repo.HasOverlap(ctx, locID, start, end)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, ErrSlotTaken
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

func isNoiseAllowed(userNoise user.NoiseLevel, locLimit location.NoiseLimit) bool {
	return userNoise.Weight() <= locLimit.Weight()
}
