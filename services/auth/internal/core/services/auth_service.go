package services

import (
	"context"
	"errors"
	"fmt"
	"theraclosure/auth-service/internal/adapters/config"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

// AuthService implements the authentication business logic
type AuthService struct {
	userRepo    ports.UserRepository
	sessionRepo ports.SessionRepository
	jwtService  ports.JWTService
	passwordSvc ports.PasswordService
	config      *config.Config
}

// NewAuthService creates a new AuthService instance
func NewAuthService(
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	config *config.Config,
) ports.AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtService:  NewJWTService(config),
		passwordSvc: NewPasswordService(),
		config:      config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	passwordHash, err := s.passwordSvc.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         domain.RoleTherapist, // Default role
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.jwtService.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.RefreshTokenDuration),
		CreatedAt:    time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *domain.AuthRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Verify password
	if err := s.passwordSvc.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	tokens, err := s.jwtService.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &domain.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(s.config.JWT.RefreshTokenDuration),
		CreatedAt:    time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResponse{
		User:         *user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get session
	session, err := s.sessionRepo.Get(ctx, claims.SessionID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if session.RefreshToken != req.RefreshToken {
		return nil, ErrInvalidToken
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	tokens, err := s.jwtService.GenerateTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update session with new refresh token
	session.RefreshToken = tokens.RefreshToken
	session.ExpiresAt = time.Now().Add(s.config.JWT.RefreshTokenDuration)

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return tokens, nil
}

// Logout invalidates a user session
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

// ValidateToken validates an access token and returns the user
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	claims, err := s.jwtService.ValidateAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// InitiateCognitoLogin returns the Cognito login URL
func (s *AuthService) InitiateCognitoLogin(ctx context.Context) (string, error) {
	// This would integrate with AWS Cognito
	// For now, return a placeholder
	return "https://cognito-login-url.com", nil
}

// HandleCognitoCallback handles the Cognito OAuth callback
func (s *AuthService) HandleCognitoCallback(ctx context.Context, code, state string) (*domain.AuthResponse, error) {
	// This would handle the Cognito OAuth callback
	// For now, return an error
	return nil, errors.New("cognito integration not implemented")
}

// PasswordService implements password hashing and verification
type PasswordService struct{}

// NewPasswordService creates a new PasswordService
func NewPasswordService() ports.PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a password using bcrypt
func (s *PasswordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (s *PasswordService) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
