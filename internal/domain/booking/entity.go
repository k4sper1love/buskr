package booking

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
	StatusNoShow    Status = "no_show"
)

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	LocationID string    `json:"location_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewBooking(userID, locationID string, start, end time.Time) *Booking {
	now := time.Now()
	return &Booking{
		ID:         uuid.New().String(),
		UserID:     userID,
		LocationID: locationID,
		StartTime:  start,
		EndTime:    end,
		Status:     StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}