package tokens

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenController struct {
	SigningKey      []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenController(signingKey []byte, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) (TokenController, error) {
	if signingKey == nil {
		return TokenController{}, errors.New("empty signing key")
	}

	return TokenController{signingKey, accessTokenTTL, refreshTokenTTL}, nil
}

func (tc *TokenController) NewJWT(email string, ipAddress string) (TokenPair, error) {
	userId := uuid.New()

	accessClaims := jwt.MapClaims{
		"id":  userId.String(),
		"ip":  ipAddress,
		"iss": "auth-server",
		"sub": email,
		"exp": time.Now().Add(tc.AccessTokenTTL).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	signedAccessToken, err := accessToken.SignedString(tc.SigningKey)

	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := tc.NewRefreshToken()

	if err != nil {
		return TokenPair{}, err
	}

	tokenPair := TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

func (tc *TokenController) Parse(accessTokenString string) (jwt.MapClaims, error) {
	accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return tc.SigningKey, nil
	})

	if err != nil {
		return jwt.MapClaims{}, err
	}

	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return jwt.MapClaims{}, err
	}

}

func (tc *TokenController) NewRefreshToken() (string, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", key), nil
}
