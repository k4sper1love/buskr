package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/k4sper1love/buskr/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryUser := `
		INSERT INTO users (id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var username sql.NullString
	if u.Username != "" {
		username = sql.NullString{String: u.Username, Valid: true}
	} else {
		username = sql.NullString{Valid: false}
	}

	_, err = tx.ExecContext(ctx, queryUser,
		u.ID,
		u.TelegramID,
		username,
		u.Name,
		string(u.Role),
		string(u.Status),
		string(u.NoiseLevel),
		u.Karma,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	queryStats := `
		INSERT INTO user_stats (user_id, total_bookings, successful_checkins, no_shows)
		VALUES ($1, 0, 0, 0)
	`

	_, err = tx.ExecContext(ctx, queryStats, u.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u user.User
	var role, status, noiseLevel string
	var username sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.TelegramID,
		&username,
		&u.Name,
		&role,
		&status,
		&noiseLevel,
		&u.Karma,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if username.Valid {
		u.Username = username.String
	}
	u.Role = user.Role(role)
	u.Status = user.Status(status)
	u.NoiseLevel = user.NoiseLevel(noiseLevel)

	return &u, nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*user.User, error) {
	query := `
		SELECT id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`

	var u user.User
	var role, status, noiseLevel string
	var username sql.NullString

	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&u.ID,
		&u.TelegramID,
		&username,
		&u.Name,
		&role,
		&status,
		&noiseLevel,
		&u.Karma,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if username.Valid {
		u.Username = username.String
	}
	u.Role = user.Role(role)
	u.Status = user.Status(status)
	u.NoiseLevel = user.NoiseLevel(noiseLevel)

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET username = $1, name = $2, role = $3, status = $4, noise_level = $5, karma = $6, updated_at = $7
		WHERE id = $8
	`

	var username sql.NullString
	if u.Username != "" {
		username = sql.NullString{String: u.Username, Valid: true}
	} else {
		username = sql.NullString{Valid: false}
	}

	res, err := r.db.ExecContext(ctx, query,
		username,
		u.Name,
		string(u.Role),
		string(u.Status),
		string(u.NoiseLevel),
		u.Karma,
		u.UpdatedAt,
		u.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetStats(ctx context.Context, userID string) (*user.UserStats, error) {
	query := `
		SELECT user_id, total_bookings, successful_checkins, no_shows
		FROM user_stats
		WHERE user_id = $1
	`

	var stats user.UserStats
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&stats.UserID,
		&stats.TotalBookings,
		&stats.SuccessfulCheckins,
		&stats.NoShows,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// stats should always exist, but if not, return empty rather than crashing
			return &user.UserStats{UserID: userID}, nil
		}
		return nil, err
	}

	return &stats, nil
}

func (r *UserRepository) UpdateStats(ctx context.Context, stats *user.UserStats) error {
	query := `
		UPDATE user_stats
		SET total_bookings = $1, successful_checkins = $2, no_shows = $3
		WHERE user_id = $4
	`

	res, err := r.db.ExecContext(ctx, query,
		stats.TotalBookings,
		stats.SuccessfulCheckins,
		stats.NoShows,
		stats.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// if update fails because row doesn't exist, try to insert it (fallback)
		insertQuery := `
			INSERT INTO user_stats (user_id, total_bookings, successful_checkins, no_shows)
			VALUES ($1, $2, $3, $4)
		`
		_, err = r.db.ExecContext(ctx, insertQuery, stats.UserID, stats.TotalBookings, stats.SuccessfulCheckins, stats.NoShows)
		return err
	}

	return nil
}

func (r *UserRepository) CreateInvite(ctx context.Context, invite *user.Invite) error {
	query := `
		INSERT INTO invites (id, token, noise_level, is_used, created_at, expires_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var createdBy sql.NullString
	if invite.CreatedBy != "" {
		createdBy = sql.NullString{String: invite.CreatedBy, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		invite.ID,
		invite.Token,
		string(invite.NoiseLevel),
		invite.IsUsed,
		invite.CreatedAt,
		invite.ExpiresAt,
		createdBy,
	)

	return err
}

func (r *UserRepository) GetInviteByToken(ctx context.Context, token string) (*user.Invite, error) {
	query := `
		SELECT id, token, noise_level, is_used, created_at, expires_at, created_by, used_by
		FROM invites
		WHERE token = $1
	`

	var inv user.Invite
	var noise string
	var createdBy, usedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&inv.ID,
		&inv.Token,
		&noise,
		&inv.IsUsed,
		&inv.CreatedAt,
		&inv.ExpiresAt,
		&createdBy,
		&usedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrInviteNotFound
		}
		return nil, err
	}

	inv.NoiseLevel = user.NoiseLevel(noise)
	if createdBy.Valid {
		inv.CreatedBy = createdBy.String
	}
	if usedBy.Valid {
		inv.UsedBy = usedBy.String
	}
	return &inv, nil
}

func (r *UserRepository) UpdateInvite(ctx context.Context, invite *user.Invite) error {
	query := `
		UPDATE invites
		SET is_used = $1, used_by = $2
		WHERE id = $3
	`

	var usedBy sql.NullString
	if invite.UsedBy != "" {
		usedBy = sql.NullString{String: invite.UsedBy, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query, invite.IsUsed, usedBy, invite.ID)
	return err
}

func (r *UserRepository) GetInviteByUsedByID(ctx context.Context, userID string) (*user.Invite, error) {
	query := `
		SELECT id, token, noise_level, is_used, created_at, expires_at, created_by, used_by
		FROM invites
		WHERE used_by = $1
	`

	var inv user.Invite
	var noise string
	var createdBy, usedBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&inv.ID,
		&inv.Token,
		&noise,
		&inv.IsUsed,
		&inv.CreatedAt,
		&inv.ExpiresAt,
		&createdBy,
		&usedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrInviteNotFound
		}
		return nil, err
	}

	inv.NoiseLevel = user.NoiseLevel(noise)
	if createdBy.Valid {
		inv.CreatedBy = createdBy.String
	}
	if usedBy.Valid {
		inv.UsedBy = usedBy.String
	}
	return &inv, nil
}


func (r *UserRepository) GetActiveUsers(ctx context.Context) ([]*user.User, error) {
	query := `
		SELECT id, telegram_id, role, status, noise_level, created_at, updated_at, username
		FROM users
		WHERE status = 'active'
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		var role, status, noiseLevel string
		var username sql.NullString

		err := rows.Scan(
			&u.ID,
			&u.TelegramID,
			&role,
			&status,
			&noiseLevel,
			&u.CreatedAt,
			&u.UpdatedAt,
			&username,
		)
		if err != nil {
			return nil, err
		}

		u.Role = user.Role(role)
		u.Status = user.Status(status)
		u.NoiseLevel = user.NoiseLevel(noiseLevel)
		if username.Valid {
			u.Username = username.String
		}

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
        SELECT id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	var u user.User
	var role, status, noiseLevel string
	var uname sql.NullString

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.TelegramID, &uname, &u.Name, &role, &status, &noiseLevel, &u.Karma, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	if uname.Valid {
		u.Username = uname.String
	}
	u.Role = user.Role(role)
	u.Status = user.Status(status)
	u.NoiseLevel = user.NoiseLevel(noiseLevel)

	return &u, nil
}

func (r *UserRepository) GetUsersPaginated(ctx context.Context, offset, limit int, sortBy string) ([]*user.User, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "created_at DESC"
	switch sortBy {
	case "karma_asc", "karma":
		orderBy = "karma ASC, created_at DESC"
	case "role":
		orderBy = "role ASC, created_at DESC"
	case "name":
		orderBy = "name ASC"
	}

	query := fmt.Sprintf(`
		SELECT id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at
		FROM users
		ORDER BY %s
		LIMIT $1 OFFSET $2
	`, orderBy)
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		var role, status, noiseLevel string
		var username sql.NullString

		err := rows.Scan(
			&u.ID, &u.TelegramID, &username, &u.Name, &role, &status, &noiseLevel, &u.Karma, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if username.Valid {
			u.Username = username.String
		}
		u.Role = user.Role(role)
		u.Status = user.Status(status)
		u.NoiseLevel = user.NoiseLevel(noiseLevel)

		users = append(users, &u)
	}

	return users, total, rows.Err()
}

func (r *UserRepository) FindByQuery(ctx context.Context, query string, offset, limit int) ([]*user.User, int, error) {
	likePattern := "%" + query + "%"

	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE name ILIKE $1 
		   OR username ILIKE $1 
		   OR CAST(telegram_id AS TEXT) = $2
	`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, likePattern, query).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sqlQuery := `
		SELECT id, telegram_id, username, name, role, status, noise_level, karma, created_at, updated_at
		FROM users
		WHERE name ILIKE $1 
		   OR username ILIKE $1 
		   OR CAST(telegram_id AS TEXT) = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, likePattern, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		var role, status, noiseLevel string
		var username sql.NullString

		err := rows.Scan(
			&u.ID, &u.TelegramID, &username, &u.Name, &role, &status, &noiseLevel, &u.Karma, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if username.Valid {
			u.Username = username.String
		}
		u.Role = user.Role(role)
		u.Status = user.Status(status)
		u.NoiseLevel = user.NoiseLevel(noiseLevel)

		users = append(users, &u)
	}
	return users, total, rows.Err()
}
