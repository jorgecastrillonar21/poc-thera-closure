package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// SessionRepository implements session repository using Redis
type SessionRepository struct {
	client *redis.Client
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(client *redis.Client) ports.SessionRepository {
	return &SessionRepository{client: client}
}

// Create creates a new session in Redis
func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.ID.String())
	ttl := time.Until(session.ExpiresAt)

	return r.client.Set(ctx, key, sessionJSON, ttl).Err()
}

// GetByID retrieves a session by ID from Redis
func (r *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID.String())
	sessionJSON, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// GetByRefreshTokenHash retrieves a session by refresh token hash from Redis
func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, tokenHash string) (*domain.Session, error) {
	pattern := "session:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}

	for _, key := range keys {
		sessionJSON, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip invalid sessions
		}

		var session domain.Session
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			continue // Skip invalid sessions
		}

		if session.RefreshTokenHash == tokenHash && session.IsActive {
			return &session, nil
		}
	}

	return nil, fmt.Errorf("session not found for refresh token hash")
}

// GetByAccessTokenJTI retrieves a session by access token JTI from Redis
func (r *SessionRepository) GetByAccessTokenJTI(ctx context.Context, jti string) (*domain.Session, error) {
	pattern := "session:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}

	for _, key := range keys {
		sessionJSON, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip invalid sessions
		}

		var session domain.Session
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			continue // Skip invalid sessions
		}

		if session.AccessTokenJTI == jti && session.IsActive {
			return &session, nil
		}
	}

	return nil, fmt.Errorf("session not found for access token JTI")
}

// Update updates an existing session in Redis
func (r *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.ID.String())
	ttl := time.Until(session.ExpiresAt)

	return r.client.Set(ctx, key, sessionJSON, ttl).Err()
}

// GetByUserID retrieves all sessions for a user ID from Redis
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	pattern := "session:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %w", err)
	}

	var sessions []*domain.Session
	for _, key := range keys {
		sessionJSON, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip invalid sessions
		}

		var session domain.Session
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			continue // Skip invalid sessions
		}

		if session.UserID == userID {
			sessions = append(sessions, &session)
		}
	}

	return sessions, nil
}

// Delete deletes a session by ID from Redis
func (r *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	key := fmt.Sprintf("session:%s", sessionID.String())
	return r.client.Del(ctx, key).Err()
}

// DeleteByUserID deletes all sessions for a user ID from Redis
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	sessions, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := r.Delete(ctx, session.ID); err != nil {
			return err
		}
	}

	return nil
}

// InvalidateByID marks a session as inactive (soft delete for logout)
func (r *SessionRepository) InvalidateByID(ctx context.Context, sessionID uuid.UUID) error {
	session, err := r.GetByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session for invalidation: %w", err)
	}

	session.IsActive = false
	session.UpdatedAt = time.Now()

	return r.Update(ctx, session)
}

// CleanupExpired removes expired and inactive sessions from Redis
func (r *SessionRepository) CleanupExpired(ctx context.Context) (int64, error) {
	// Redis TTL automatically handles expiration, but we can scan for inactive sessions
	pattern := "session:*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get session keys: %w", err)
	}

	var deletedCount int64
	for _, key := range keys {
		sessionJSON, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue // Skip if session doesn't exist
		}

		var session domain.Session
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			continue // Skip invalid sessions
		}

		// Delete if session is inactive or expired
		if !session.IsActive || time.Now().After(session.ExpiresAt) {
			if err := r.client.Del(ctx, key).Err(); err == nil {
				deletedCount++
			}
		}
	}

	return deletedCount, nil
}
