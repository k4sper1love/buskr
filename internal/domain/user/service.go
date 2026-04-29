package user

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrCreateUser(ctx context.Context, telegramID int64, username string) (*User, error) {
	u, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		newUser := NewUser(telegramID, username)
		if err := s.repo.Create(ctx, newUser); err != nil {
			return nil, err
		}

		stats := &UserStats{UserID: newUser.ID}
		_ = s.repo.UpdateStats(ctx, stats)

		return newUser, nil
	}

	return u, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) SubmitApplication(ctx context.Context, userID string) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.Status != StatusGuest {
		return ErrInvalidStatus
	}

	u.Status = StatusPending
	u.UpdatedAt = time.Now()

	return s.repo.Update(ctx, u)
}

func (s *Service) ApproveUser(ctx context.Context, actorID, targetUserID string, noise NoiseLevel) error {
	actor, err := s.repo.GetByID(ctx, actorID)
	if err != nil {
		return err
	}

	if actor.Role != RoleAdmin {
		return ErrAccessDenied
	}

	targetUser, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if targetUser.Status != StatusPending {
		return ErrInvalidStatus
	}

	targetUser.Status = StatusActive
	targetUser.NoiseLevel = noise
	targetUser.UpdatedAt = time.Now()

	return s.repo.Update(ctx, targetUser)
}

func (s *Service) ApproveApplication(ctx context.Context, userID string, noiseLevel NoiseLevel) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.Status != StatusPending {
		return ErrInvalidStatus
	}

	u.Status = StatusActive
	u.NoiseLevel = noiseLevel
	u.UpdatedAt = time.Now()

	return s.repo.Update(ctx, u)
}

func (s *Service) RejectApplication(ctx context.Context, userID string) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.Status != StatusPending {
		return ErrInvalidStatus
	}

	u.Status = StatusGuest
	u.UpdatedAt = time.Now()

	return s.repo.Update(ctx, u)
}

func (s *Service) ActiveViaInvite(ctx context.Context, userID string, noise NoiseLevel) error {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Status = StatusActive
	u.NoiseLevel = noise
	u.UpdatedAt = time.Now()

	return s.repo.Update(ctx, u)
}

func (s *Service) RecordSuccessfulCheckin(ctx context.Context, userID string) error {
	stats, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		return err
	}

	stats.SuccessfulCheckins++
	if err := s.repo.UpdateStats(ctx, stats); err != nil {
		return err
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Karma += KarmaReward

	// check karma bounds
	if u.Karma > MaxKarma {
		u.Karma = MaxKarma
	}

	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *Service) RecordNoShow(ctx context.Context, userID string) error {
	stats, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		return err
	}

	stats.NoShows++
	if err := s.repo.UpdateStats(ctx, stats); err != nil {
		return err
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Karma -= KarmaPenalty

	// check karma bounds
	if u.Karma < MinKarma {
		u.Karma = MinKarma
	}

	u.UpdatedAt = time.Now()
	return s.repo.Update(ctx, u)
}

func (s *Service) BanUser(ctx context.Context, actorID, targetUserID string) error {
	actor, err := s.repo.GetByID(ctx, actorID)
	if err != nil {
		return err
	}

	if actor.Role != RoleAdmin {
		return ErrAccessDenied
	}

	targetUser, err := s.repo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if targetUser.Role == RoleAdmin {
		return ErrAccessDenied
	}

	targetUser.Status = StatusBanned
	targetUser.UpdatedAt = time.Now()

	return s.repo.Update(ctx, targetUser)
}

func (s *Service) GenerateInvite(ctx context.Context, noiseLevel NoiseLevel, expirationHours int) (*Invite, error) {
	invite := NewInvite(noiseLevel, expirationHours)

	err := s.repo.CreateInvite(ctx, invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *Service) ProcessInvite(ctx context.Context, userID string, token string) error {
	inv, err := s.repo.GetInviteByToken(ctx, token)
	if err != nil {
		return ErrInviteNotFound
	}

	if inv.IsUsed {
		return ErrInviteAlreadyUsed
	}
	if time.Now().After(inv.ExpiresAt) {
		return ErrInviteExpired
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.Status != StatusGuest && u.Status != StatusPending {
		return ErrInvalidStatus
	}

	u.Status = StatusActive
	u.NoiseLevel = inv.NoiseLevel
	u.UpdatedAt = time.Now()

	err = s.repo.Update(ctx, u)
	if err != nil {
		return err
	}

	inv.IsUsed = true
	err = s.repo.UpdateInvite(ctx, inv)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetActiveUsers(ctx context.Context) ([]*User, error) {
	return s.repo.GetActiveUsers(ctx)
}

func (s *Service) GetStats(ctx context.Context, userID string) (*UserStats, error) {
	return s.repo.GetStats(ctx, userID)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) GetByTelegramID(ctx context.Context, telegramID int64) (*User, error) {
	return s.repo.GetByTelegramID(ctx, telegramID)
}

func (s *Service) Update(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *Service) ForgiveFines(ctx context.Context, userID string) error {
	stats, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		return err
	}
	stats.NoShows = 0
	if err := s.repo.UpdateStats(ctx, stats); err != nil {
		return err
	}

	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	u.Karma = MaxKarma
	u.UpdatedAt = time.Now()

	return s.repo.Update(ctx, u)
}
