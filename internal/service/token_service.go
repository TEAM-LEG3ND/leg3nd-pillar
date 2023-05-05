package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"leg3nd-pillar/internal/config"
	"strconv"
	"time"
)

func GetAccessToken(id int64, duration time.Duration) (*string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(id, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret, err := config.GetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	t, err := token.SignedString([]byte(*jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("token generation failed, %w", err)
	}
	return &t, nil
}

func GetRefreshToken(id int64, duration time.Duration) (*string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(id, 10),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtRefreshSecret, err := config.GetEnv("JWT_REFRESH_SECRET")
	if err != nil {
		return nil, err
	}

	t, err := token.SignedString([]byte(*jwtRefreshSecret))
	if err != nil {
		return nil, fmt.Errorf("token generation failed, %w", err)
	}
	return &t, nil
}
