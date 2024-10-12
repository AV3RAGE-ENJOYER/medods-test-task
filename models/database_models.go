package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (u User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err
}

func (u User) String() string {
	return fmt.Sprintf("User {Email: %s, HashedPassword: %s}", u.Email, u.HashedPassword)
}

type RefreshToken struct {
	RefreshTokenHash string    `json:"refresh_token_hash"`
	ExpiresAt        time.Time `json:"expires_at"`
	IpAddress        string    `json:"ip_address"`
}
