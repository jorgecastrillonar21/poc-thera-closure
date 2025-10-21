package persistence

import (
	"context"
	"crypto/sha256"
	"fmt"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SessionRepository implements session repository using PostgreSQL
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(db *gorm.DB) ports.SessionRepository {
	return &SessionRepository{db: db}
}

// Create creates a new session in the database
func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetByID retrieves a session by ID
func (r *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ? AND expires_at > ?", sessionID, true, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByRefreshTokenHash retrieves a session by refresh token hash
func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.WithContext(ctx).
		Where("refresh_token_hash = ? AND is_active = ? AND expires_at > ?", tokenHash, true, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByAccessTokenJTI retrieves a session by access token JTI
func (r *SessionRepository) GetByAccessTokenJTI(ctx context.Context, jti string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.WithContext(ctx).
		Where("access_token_jti = ? AND is_active = ? AND expires_at > ?", jti, true, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetByUserID retrieves all active sessions for a user
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	var sessions []*domain.Session
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error

	return sessions, err
}

// Update updates an existing session
func (r *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// Delete permanently deletes a session by ID
func (r *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Session{}, "id = ?", sessionID).Error
}

// DeleteByUserID permanently deletes all sessions for a user
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Session{}, "user_id = ?", userID).Error
}

// InvalidateByID marks a session as inactive (soft delete for logout)
func (r *SessionRepository) InvalidateByID(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&domain.Session{}).
		Where("id = ?", sessionID).
		Update("is_active", false).Error
}

// CleanupExpired removes expired and inactive sessions
func (r *SessionRepository) CleanupExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ? OR is_active = ?", time.Now(), false).
		Delete(&domain.Session{})

	return result.RowsAffected, result.Error
}

// HashToken creates a SHA-256 hash of the token for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}
