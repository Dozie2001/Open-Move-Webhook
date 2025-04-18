package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
)

// specifically for the zklogin handlers

func GenerateTokens(sub, email string) (string, string, error) {
	// access token (15 min)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})

	accessString, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// refresh token (7 days)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
}
