package database

import (
	"context"
	"medods_test_task/models"
	"medods_test_task/service"
	"time"

	"github.com/jackc/pgx/v5"
)

type MockDBUserRepository struct {
	Users         map[string]models.User
	RefreshTokens map[string]models.RefreshToken
}

var _ service.UserRepository = &MockDBUserRepository{}

func (db *MockDBUserRepository) GetUser(ctx context.Context, email string) (models.User, error) {
	user, ok := db.Users[email]

	if !ok {
		return models.User{}, pgx.ErrNoRows
	}

	return user, nil
}

func (db *MockDBUserRepository) AddUser(ctx context.Context, user models.User) error {
	db.Users[user.Email] = user
	return nil
}

func (db *MockDBUserRepository) AddRefreshToken(ctx context.Context, email string, refreshTokenHash string, ttl time.Duration, ipAddress string) error {
	refreshToken := models.RefreshToken{
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(ttl),
		IpAddress:        ipAddress,
	}

	db.RefreshTokens[email] = refreshToken
	return nil
}

func (db *MockDBUserRepository) GetRefreshTokenProps(ctx context.Context, email string) (models.RefreshToken, error) {
	refreshToken, ok := db.RefreshTokens[email]

	if !ok {
		return models.RefreshToken{}, pgx.ErrNoRows
	}

	return refreshToken, nil
}
