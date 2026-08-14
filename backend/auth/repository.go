package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UserExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, id uuid.UUID, email, passwordHash string) error
	GetUserCredentials(ctx context.Context, email string) (uuid.UUID, string, error)
	CreateSession(ctx context.Context, sessionID, userID uuid.UUID) error
	GetUserIDBySession(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM users WHERE email = $1
		)`,
		email,
	).Scan(&exists)

	return exists, err
}

func (r *repository) CreateUser(
	ctx context.Context,
	id uuid.UUID,
	email, passwordHash string,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users (id, email, password_hash)
		 VALUES ($1, $2, $3)`,
		id,
		email,
		passwordHash,
	)

	return err
}

func (r *repository) GetUserCredentials(
	ctx context.Context,
	email string,
) (uuid.UUID, string, error) {
	var (
		userID       uuid.UUID
		passwordHash string
	)

	err := r.db.QueryRow(
		ctx,
		`SELECT id, password_hash
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&userID, &passwordHash)

	return userID, passwordHash, err
}

func (r *repository) CreateSession(
	ctx context.Context,
	sessionID, userID uuid.UUID,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO sessions (id, user_id, expires_at)
		 VALUES ($1, $2, NOW() + INTERVAL '24 hours')`,
		sessionID,
		userID,
	)

	return err
}

func (r *repository) GetUserIDBySession(
	ctx context.Context,
	sessionID uuid.UUID,
) (uuid.UUID, error) {
	var userID uuid.UUID

	err := r.db.QueryRow(
		ctx,
		`SELECT user_id
		 FROM sessions
		 WHERE id = $1
		   AND expires_at > NOW()`,
		sessionID,
	).Scan(&userID)

	return userID, err
}