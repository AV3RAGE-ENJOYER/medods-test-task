package database

import (
	"context"
	"medods_test_task/models"
	"medods_test_task/service"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDBUserRepository struct {
	Pool *pgxpool.Pool
}

var _ service.UserRepository = &PostgresDBUserRepository{}

func NewDB(ctx context.Context, DB_URL string) (PostgresDBUserRepository, error) {
	dbPool, err := pgxpool.New(ctx, DB_URL)

	if err != nil {
		return PostgresDBUserRepository{}, err
	}

	postgres := PostgresDBUserRepository{dbPool}
	return postgres, nil
}

func (db *PostgresDBUserRepository) GetUser(ctx context.Context, email string) (models.User, error) {
	args := pgx.NamedArgs{
		"email": email,
	}

	q := `SELECT users.email, auth.password_hash
		FROM users JOIN auth ON users.id=auth.id
		WHERE users.email=@email`

	row := db.Pool.QueryRow(ctx, q, args)

	var user models.User

	err := row.Scan(
		&user.Email,
		&user.HashedPassword,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *PostgresDBUserRepository) AddUser(ctx context.Context, user models.User) error {
	tx, err := db.Pool.Begin(ctx)

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	args := pgx.NamedArgs{
		"email":         user.Email,
		"password_hash": user.HashedPassword,
	}

	q := `WITH add_user AS (
		INSERT INTO users("email")
		VALUES (@email)
		RETURNING id)
		INSERT INTO auth("id", "password_hash")
		SELECT id, @password_hash FROM add_user`

	_, err = tx.Exec(
		ctx,
		q,
		args,
	)

	return err
}

func (db *PostgresDBUserRepository) AddRefreshToken(ctx context.Context, email string, refreshTokenHash string, ttl time.Duration, ipAddress string) error {
	tx, err := db.Pool.Begin(ctx)

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	args := pgx.NamedArgs{
		"email":              email,
		"refresh_token_hash": refreshTokenHash,
		"ttl":                time.Now().Add(ttl),
		"ip":                 ipAddress,
	}

	q := `INSERT INTO
		refresh_tokens(
			"id", "refresh_token_hash", 
			"expires_at", "last_ip")
		SELECT id, @refresh_token_hash, @ttl, @ip
		FROM users
		WHERE users.email = @email
		ON CONFLICT(id)
		DO UPDATE SET 
		refresh_token_hash=@refresh_token_hash,
		expires_at=@ttl, last_ip=@ip`

	_, err = tx.Exec(ctx, q, args)

	return err
}

func (db *PostgresDBUserRepository) GetRefreshTokenProps(ctx context.Context, email string) (models.RefreshToken, error) {
	args := pgx.NamedArgs{
		"email": email,
	}

	q := `SELECT 
		refresh_tokens.refresh_token_hash,
		refresh_tokens.expires_at,
		refresh_tokens.last_ip
		FROM refresh_tokens
		JOIN users ON refresh_tokens.id=users.id
		WHERE users.email=@email`

	var refreshToken models.RefreshToken
	row := db.Pool.QueryRow(ctx, q, args)

	err := row.Scan(
		&refreshToken.RefreshTokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.IpAddress,
	)

	if err != nil {
		return models.RefreshToken{}, err
	}

	return refreshToken, nil
}
