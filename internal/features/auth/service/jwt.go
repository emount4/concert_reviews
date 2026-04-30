package auth_service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	core_domain "github.com/emount4/concert_reviews/internal/core/domain"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTManager interface {
	NewAccessToken(userID uuid.UUID, roleID int, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
	Generate(userID uuid.UUID, roleID int, ttl time.Duration) (core_domain.AuthResponse, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) *Manager {
	return &Manager{signingKey: signingKey}
}

func (m *Manager) Generate(userID uuid.UUID, roleID int, ttl time.Duration) (core_domain.AuthResponse, error) {
	now := time.Now()
	accessToken, err := m.NewAccessToken(userID, roleID, ttl)
	if err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := m.NewRefreshToken()
	if err != nil {
		return core_domain.AuthResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	return core_domain.AuthResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresAt:    now.Add(ttl),
	}, nil

}

type JWTClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID
	RoleID int
}

func (m *Manager) NewAccessToken(userID uuid.UUID, roleID int, ttl time.Duration) (string, error) {

	claims := JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Subject:   userID.String(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
		RoleID: roleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("read crypto random: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(m.signingKey), nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("parse access token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}
