package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/k4sper1love/buskr/internal/domain/location"
)

type LocationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) Create(ctx context.Context, loc *location.Location) error {
	// ST_SetSRID(ST_MakePoint(longitude, latitude), 4326) creates a standard GPS point
	query := `
		INSERT INTO locations (id, name, description, coords, max_noise, status, created_at, updated_at)
		VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326), $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		loc.ID,
		loc.Name,
		loc.Description,
		loc.Coords.Lon, // X
		loc.Coords.Lat, // Y
		string(loc.MaxNoise),
		string(loc.Status),
		loc.CreatedAt,
		loc.UpdatedAt,
	)

	return err
}

func (r *LocationRepository) GetByID(ctx context.Context, id string) (*location.Location, error) {
	// ST_Y extracts latitude, ST_X extracts longitude
	query := `
		SELECT 
			id, name, description, 
			ST_Y(coords::geometry) as lat, ST_X(coords::geometry) as lon, 
			max_noise, status, created_at, updated_at
		FROM locations
		WHERE id = $1
	`

	var loc location.Location
	var maxNoise, status string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&loc.ID,
		&loc.Name,
		&loc.Description,
		&loc.Coords.Lat,
		&loc.Coords.Lon,
		&maxNoise,
		&status,
		&loc.CreatedAt,
		&loc.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, location.ErrLocationNotFound
		}
		return nil, err
	}

	loc.MaxNoise = location.NoiseLimit(maxNoise)
	loc.Status = location.Status(status)

	return &loc, nil
}

func (r *LocationRepository) List(ctx context.Context, includeInactive bool) ([]*location.Location, error) {
	query := `
		SELECT 
			id, name, description, 
			ST_Y(coords::geometry) as lat, ST_X(coords::geometry) as lon, 
			max_noise, status, created_at, updated_at
		FROM locations
		WHERE status = 'active' OR $1 = true
	`

	rows, err := r.db.QueryContext(ctx, query, includeInactive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*location.Location

	for rows.Next() {
		var loc location.Location
		var maxNoise, status string

		err := rows.Scan(
			&loc.ID,
			&loc.Name,
			&loc.Description,
			&loc.Coords.Lat,
			&loc.Coords.Lon,
			&maxNoise,
			&status,
			&loc.CreatedAt,
			&loc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		loc.MaxNoise = location.NoiseLimit(maxNoise)
		loc.Status = location.Status(status)

		locations = append(locations, &loc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) Update(ctx context.Context, loc *location.Location) error {
	query := `
		UPDATE locations
		SET 
			name = $1, 
			description = $2, 
			coords = ST_SetSRID(ST_MakePoint($3, $4), 4326), 
			max_noise = $5, 
			status = $6, 
			updated_at = $7
		WHERE id = $8
	`

	res, err := r.db.ExecContext(ctx, query,
		loc.Name,
		loc.Description,
		loc.Coords.Lon,
		loc.Coords.Lat,
		string(loc.MaxNoise),
		string(loc.Status),
		loc.UpdatedAt,
		loc.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return location.ErrLocationNotFound
	}

	return nil
}
