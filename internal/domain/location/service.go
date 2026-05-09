package location

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

func (s *Service) CreateLocation(ctx context.Context, name, desc string, lat, lon float64, noise NoiseLimit) (*Location, error) {
	loc := NewLocation(name, desc, lat, lon, noise)

	if err := s.repo.Create(ctx, loc); err != nil {
		return nil, err
	}

	return loc, nil
}

func (s *Service) UpdateLocation(ctx context.Context, id string, name, desc string, lat, lon float64, noise NoiseLimit) error {
	loc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	loc.Name = name
	loc.Description = desc
	loc.Coords = Point{Lat: lat, Lon: lon}
	loc.MaxNoise = noise
	loc.UpdatedAt = time.Now()

	return s.repo.Update(ctx, loc)
}

func (s *Service) ChangeStatus(ctx context.Context, id string, status Status) error {
	loc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	loc.Status = status
	loc.UpdatedAt = time.Now()

	return s.repo.Update(ctx, loc)
}

func (s *Service) GetLocationsForAdmin(ctx context.Context) ([]*Location, error) {
	return s.repo.List(ctx, true)
}

func (s *Service) GetLocationsForMusicians(ctx context.Context) ([]*Location, error) {
	return s.repo.List(ctx, false)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Location, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) FindNearby(ctx context.Context, lat, lon float64, radiusMeters float64) ([]*LocationWithDist, error) {
	return s.repo.FindNearby(ctx, lat, lon, radiusMeters)
}