package user

import (
	"time"

	"github.com/google/uuid"
)

// Role defines basic access rights
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleMusician Role = "musician"
)

// Status describes the current stage of onboarding
type Status string

const (
	StatusGuest   Status = "guest"   // no access to the booking
	StatusPending Status = "pending" // waiting for approval
	StatusActive  Status = "active"  // approved, can book
	StatusBanned  Status = "banned"  // banned manually or for karma
)

// NoiseLevel classifies a musician by noise level
type NoiseLevel string

const (
	NoiseNone   NoiseLevel = "none"   // for guest/pending
	NoiseLight  NoiseLevel = "light"    // acoustics, violin
	NoiseMedium NoiseLevel = "medium" // lightweight combo amplifiers
	NoiseHard   NoiseLevel = "hard"   // drums, powerful speakers
)

// User represents a registered user in the system
type User struct {
	ID         string     `json:"id"`
	TelegramID int64      `json:"telegram_id"`
	Username   string     `json:"username"`
	Name       string     `json:"name"` // artist name or group name
	Role       Role       `json:"role"`
	Status     Status     `json:"status"`
	NoiseLevel NoiseLevel `json:"noise_level"`
	Karma      int        `json:"karma"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// UserStats stores statistics for shadow analytics
type UserStats struct {
	UserID             string `json:"user_id"`
	TotalBookings      int    `json:"total_bookings"`
	SuccessfulCheckins int    `json:"successful_checkins"`
	NoShows            int    `json:"no_shows"`
}

type Invite struct {
	ID         string
	Token      string
	NoiseLevel NoiseLevel
	IsUsed     bool
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

func NewUser(telegramID int64, username string, startingKarma int) *User {
	now := time.Now()
	return &User{
		ID:         uuid.New().String(),
		TelegramID: telegramID,
		Username:   username,
		Role:       RoleMusician,
		Status:     StatusGuest,
		NoiseLevel: NoiseNone,
		Karma:      startingKarma,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func NewInvite(noiseLevel NoiseLevel, expirationHours int) *Invite {
	now := time.Now()
	token := uuid.New().String() // generate unique token

	return &Invite{
		ID:         uuid.New().String(),
		Token:      token,
		NoiseLevel: noiseLevel,
		IsUsed:     false,
		CreatedAt:  now,
		ExpiresAt:  now.Add(time.Hour * time.Duration(expirationHours)),
	}
}

func (n NoiseLevel) Weight() int {
	switch n {
	case NoiseNone:
		return 0
	case NoiseLight:
		return 1
	case NoiseMedium:
		return 2
	case NoiseHard:
		return 3
	default:
		return 0
	}
}
