package ports

import (
	"context"
	"theraclosure/auth-service/internal/core/domain"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByCognitoID(ctx context.Context, cognitoID string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// SessionRepository defines the interface for session management
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error)
	GetByAccessTokenJTI(ctx context.Context, jti string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	InvalidateByID(ctx context.Context, sessionID uuid.UUID) error
	CleanupExpired(ctx context.Context) (int64, error)
}

// AuthService defines the interface for authentication business logic
type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.AuthRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
	InitiateCognitoLogin(ctx context.Context) (string, error)
	HandleCognitoCallback(ctx context.Context, code, state string) (*domain.AuthResponse, error)
}

// UserService defines the interface for user management business logic
type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*domain.User, error)
	UpdateSubscriptionStatus(ctx context.Context, userID uuid.UUID, status domain.SubscriptionStatus) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// CognitoService defines the interface for AWS Cognito integration
type CognitoService interface {
	GetLoginURL(ctx context.Context, state string) string
	ExchangeCodeForTokens(ctx context.Context, code string) (*CognitoTokens, error)
	GetUserInfo(ctx context.Context, accessToken string) (*domain.CognitoUser, error)
}

// CognitoTokens represents the tokens returned by Cognito
type CognitoTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// JWTService defines the interface for JWT token operations
type JWTService interface {
	GenerateTokenPair(ctx context.Context, user *domain.User, sessionID uuid.UUID) (*domain.TokenPair, error)
	ValidateAccessToken(ctx context.Context, token string) (*JWTClaims, error)
	ValidateRefreshToken(ctx context.Context, token string) (*JWTClaims, error)
}

// JWTClaims represents the JWT token claims (simplified for ports interface)
type JWTClaims struct {
	UserID    uuid.UUID       `json:"user_id"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	SessionID uuid.UUID       `json:"session_id"`
	Type      string          `json:"type"` // "access" or "refresh"
}

// PasswordService defines the interface for password operations
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}
