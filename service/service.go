package service

import (
	"context"
	"medods_test_task/models"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	AddUser(ctx context.Context, user models.User) error
	AddRefreshToken(ctx context.Context, email string, refreshTokenHash string, ttl time.Duration, ipAddress string) error
	GetRefreshTokenProps(ctx context.Context, email string) (models.RefreshToken, error)
}

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}
