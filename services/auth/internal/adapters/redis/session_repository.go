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

	key := fmt.Sprintf("session:%s", session.ID)
	ttl := time.Until(session.ExpiresAt)

	return r.client.Set(ctx, key, sessionJSON, ttl).Err()
}

// Get retrieves a session by ID from Redis
func (r *SessionRepository) Get(ctx context.Context, sessionID string) (*domain.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
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

// GetByUserID retrieves all sessions for a user ID from Redis
func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	pattern := fmt.Sprintf("session:*")
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
func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
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

// DeleteExpired deletes all expired sessions from Redis
func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	// Redis TTL automatically handles expiration, so this is a no-op
	// In a real implementation, you might want to scan for expired sessions
	return nil
}
