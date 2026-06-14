package location

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
)

const NearbyRadius = 100 // meters

type NoiseLimit string

const (
	LimitLight  NoiseLimit = "light"
	LimitMedium NoiseLimit = "medium"
	LimitHard   NoiseLimit = "hard"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Location struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Coords      Point      `json:"coords"`
	MaxNoise    NoiseLimit `json:"max_noise"`
	Status      Status     `json:"status"`
	SuggestedBy *string    `json:"suggested_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type LocationWithDist struct {
	*Location
	DistMeters float64 `json:"dist_meters"`
}

func NewLocation(name, desc string, lat, lon float64, maxNoise NoiseLimit) *Location {
	now := time.Now()
	return &Location{
		ID:          uuid.New().String(),
		Name:        name,
		Description: desc,
		Coords:      Point{Lat: lat, Lon: lon},
		MaxNoise:    maxNoise,
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (n NoiseLimit) Weight() int {
	switch n {
	case LimitLight:
		return 1
	case LimitMedium:
		return 2
	case LimitHard:
		return 3
	default:
		return 0
	}
}
