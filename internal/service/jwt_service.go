package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/middleware"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type JWTService interface {
	GenerateToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
}

type jwtService struct {
	secret           string
	expiryHours      int
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewJWTService(secret string, expiryHours int, refreshTokenRepo repository.RefreshTokenRepository) JWTService {
	return &jwtService{
		secret:           secret,
		expiryHours:      expiryHours,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *jwtService) GenerateToken(userID uint, email string) (string, error) {
	claims := middleware.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) GenerateRefreshToken(userID uint) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	tokenString := base64.URLEncoding.EncodeToString(b)

	refreshToken := &entity.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *jwtService) ValidateToken(tokenString string) (*middleware.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*middleware.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *jwtService) RefreshAccessToken(refreshToken string) (string, string, error) {
	storedToken, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	accessToken, err := s.GenerateToken(storedToken.UserID, "")
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := s.GenerateRefreshToken(storedToken.UserID)
	if err != nil {
		return "", "", err
	}

	_ = s.refreshTokenRepo.RevokeByToken(refreshToken)

	return accessToken, newRefreshToken, nil
}
