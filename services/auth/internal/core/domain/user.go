package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin     UserRole = "ADMIN"
	RoleTherapist UserRole = "THERAPIST"
	RoleStaff     UserRole = "STAFF"
)

// SubscriptionStatus represents the subscription status
type SubscriptionStatus string

const (
	SubscriptionActive     SubscriptionStatus = "active"
	SubscriptionCanceled   SubscriptionStatus = "canceled"
	SubscriptionTrialing   SubscriptionStatus = "trialing"
	SubscriptionPastDue    SubscriptionStatus = "past_due"
	SubscriptionIncomplete SubscriptionStatus = "incomplete"
)

// User represents a user in the authentication domain
type User struct {
	ID                 uuid.UUID          `json:"id" gorm:"type:uuid;primary_key"`
	Email              string             `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash       string             `json:"-" gorm:"not null"`
	FirstName          string             `json:"firstName" gorm:"not null"`
	LastName           string             `json:"lastName" gorm:"not null"`
	Role               UserRole           `json:"role" gorm:"not null;default:'THERAPIST'"`
	SubscriptionStatus SubscriptionStatus `json:"subscriptionStatus" gorm:"not null;default:'trialing'"`
	StripeCustomerID   *string            `json:"stripeCustomerId" gorm:"index"`
	CognitoID          *string            `json:"cognitoId" gorm:"index"`
	IsActive           bool               `json:"isActive" gorm:"not null;default:true"`
	EmailVerified      bool               `json:"emailVerified" gorm:"not null;default:false"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

// BeforeCreate sets the ID and timestamps
func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// Session represents a user session stored in PostgreSQL
type Session struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID           uuid.UUID `json:"userId" gorm:"type:uuid;not null;index"`
	RefreshTokenHash string    `json:"-" gorm:"not null;unique_index"` // Hashed refresh token
	AccessTokenJTI   string    `json:"-" gorm:"not null;unique_index"` // JWT ID for access token
	UserAgent        *string   `json:"userAgent" gorm:"type:text"`     // Nullable
	IPAddress        *string   `json:"ipAddress" gorm:"type:inet"`     // Nullable
	IsActive         bool      `json:"isActive" gorm:"not null;default:true;index"`
	ExpiresAt        time.Time `json:"expiresAt" gorm:"not null;index"`
	CreatedAt        time.Time `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"not null;default:now()"`
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	AccessTokenJTI string `json:"-"` // JWT ID for access token tracking
}

// AuthRequest represents a login request
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required,min=2"`
	LastName  string `json:"lastName" binding:"required,min=2"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// CognitoUser represents a user from AWS Cognito
type CognitoUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	EmailVerified bool   `json:"email_verified"`
}

// HasRole checks if the user has a specific role
func (u *User) HasRole(role UserRole) bool {
	return u.Role == role
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsSubscriptionActive checks if the user has an active subscription
func (u *User) IsSubscriptionActive() bool {
	return u.SubscriptionStatus == SubscriptionActive
}
