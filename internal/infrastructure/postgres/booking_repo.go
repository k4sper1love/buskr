package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/booking"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, b *booking.Booking) error {
	query := `
		INSERT INTO bookings (id, user_id, location_id, start_time, end_time, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		b.ID,
		b.UserID,
		b.LocationID,
		b.StartTime,
		b.EndTime,
		string(b.Status),
		b.CreatedAt,
		b.UpdatedAt,
	)

	return err
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`

	var b booking.Booking
	var status string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.LocationID,
		&b.StartTime,
		&b.EndTime,
		&status,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, booking.ErrBookingNotFound
		}
		return nil, err
	}

	b.Status = booking.Status(status)

	return &b, nil
}

func (r *BookingRepository) Update(ctx context.Context, b *booking.Booking) error {
	query := `
		UPDATE bookings
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	res, err := r.db.ExecContext(ctx, query,
		string(b.Status),
		b.UpdatedAt,
		b.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return booking.ErrBookingNotFound
	}

	return nil
}

func (r *BookingRepository) HasOverlapByLocation(ctx context.Context, locID string, start, end time.Time) (bool, error) {
	// overlap occurs if (existing_start < new_end) AND (existing_end > new_start)
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM bookings 
			WHERE location_id = $1 
			  AND status IN ('pending', 'active') 
			  AND start_time < $3 
			  AND end_time > $2
		)
	`

	var hasOverlap bool
	err := r.db.QueryRowContext(ctx, query, locID, start, end).Scan(&hasOverlap)
	if err != nil {
		return false, err
	}

	return hasOverlap, nil
}

func (r *BookingRepository) HasOverlapByUser(ctx context.Context, userID string, start, end time.Time) (bool, error) {
	// overlap occurs if (existing_start < new_end) AND (existing_end > new_start)
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM bookings 
			WHERE user_id = $1 
			  AND status IN ('pending', 'active') 
			  AND start_time < $3 
			  AND end_time > $2
		)
	`

	var hasOverlap bool
	err := r.db.QueryRowContext(ctx, query, userID, start, end).Scan(&hasOverlap)
	if err != nil {
		return false, err
	}

	return hasOverlap, nil
}

func (r *BookingRepository) GetByLocationAndTime(ctx context.Context, locID string, start, end time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE location_id = $1
		  AND start_time >= $2 
		  AND start_time <= $3
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, locID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepository) GetUpcoming(ctx context.Context, from, to time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE status = 'pending' AND start_time >= $1 AND start_time <= $2
	`

	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepository) GetPendingBefore(ctx context.Context, deadline time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE status = 'pending' AND start_time <= $1
	`

	rows, err := r.db.QueryContext(ctx, query, deadline)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepository) GetActiveEndedBefore(ctx context.Context, deadline time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE status = 'active' AND end_time <= $1
	`

	rows, err := r.db.QueryContext(ctx, query, deadline)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepository) GetPendingOrActiveByUser(ctx context.Context, userID string) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, location_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1 AND status IN ('pending', 'active')
		ORDER BY start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanBookings(rows)
}

func (r *BookingRepository) CountActiveFuture(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT count(*)
		FROM bookings
		WHERE user_id = $1 AND status IN ('pending', 'active') AND end_time >= NOW()
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *BookingRepository) CountByUserAndLocationAndDay(ctx context.Context, userID string, locID string, date time.Time) (int, error) {
	query := `
		SELECT count(*)
		FROM bookings
		WHERE user_id = $1 
			AND location_id = $2
			AND start_time::date = $3::date
			AND status IN ('pending', 'active', 'completed')
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID, locID, date.Format("2006-01-02")).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *BookingRepository) HasNoisyNeighbor(ctx context.Context, locID string, start, end time.Time, radiusMeters float64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM bookings
			JOIN locations nearby ON b.location_id = nearby.id
			JOIN users u ON b.user_id = u.id
			JOIN locations target ON target.id = $1
			WHERE nearby.id != $1
				AND b.status IN ('pending', 'active')
				AND b.start_time < $3
				AND b.end_time > $2
				AND u.noise_level = 'hard'
				AND ST_DWithin(
					nearby.coords::geography,
					target.coords::geography,
					$4
				)
		)
	`

	var has bool
	err := r.db.QueryRowContext(ctx, query, locID, start, end, radiusMeters).Scan(&has)
	return has, err
}

func (r *BookingRepository) scanBookings(rows *sql.Rows) ([]*booking.Booking, error) {
	var bookings []*booking.Booking

	for rows.Next() {
		var b booking.Booking
		var status string

		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.LocationID,
			&b.StartTime,
			&b.EndTime,
			&status,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		b.Status = booking.Status(status)
		bookings = append(bookings, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}
