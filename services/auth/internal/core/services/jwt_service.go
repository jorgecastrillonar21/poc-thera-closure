package services

import (
	"context"
	"fmt"
	"theraclosure/auth-service/internal/adapters/config"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims extends jwt.RegisteredClaims with our custom fields
type CustomClaims struct {
	UserID    uuid.UUID       `json:"user_id"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	SessionID uuid.UUID       `json:"session_id"`
	Type      string          `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTService implements JWT token operations
type JWTService struct {
	config *config.Config
}

// NewJWTService creates a new JWTService instance
func NewJWTService(config *config.Config) ports.JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateTokenPair generates access and refresh tokens for a user
func (s *JWTService) GenerateTokenPair(ctx context.Context, user *domain.User, sessionID uuid.UUID) (*domain.TokenPair, error) {
	now := time.Now()

	// Generate unique JTI for access token
	accessJTI := uuid.New().String()
	
	// Generate access token
	accessClaims := &CustomClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		SessionID: sessionID,
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.AccessTokenDuration)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &CustomClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		SessionID: sessionID,
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.JWT.RefreshTokenDuration)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:    accessTokenString,
		RefreshToken:   refreshTokenString,
		AccessTokenJTI: accessJTI,
	}, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *JWTService) ValidateAccessToken(ctx context.Context, tokenString string) (*ports.JWTClaims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, fmt.Errorf("invalid token type: %s", claims.Type)
	}

	// Convert CustomClaims to ports.JWTClaims
	return &ports.JWTClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		SessionID: claims.SessionID,
		Type:      claims.Type,
	}, nil
}

// ValidateRefreshToken validates a refresh token and returns claims
func (s *JWTService) ValidateRefreshToken(ctx context.Context, tokenString string) (*ports.JWTClaims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, fmt.Errorf("invalid token type: %s", claims.Type)
	}

	// Convert CustomClaims to ports.JWTClaims
	return &ports.JWTClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		SessionID: claims.SessionID,
		Type:      claims.Type,
	}, nil
}

// validateToken validates a JWT token and returns claims
func (s *JWTService) validateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
