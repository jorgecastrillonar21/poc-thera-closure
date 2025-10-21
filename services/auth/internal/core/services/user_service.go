package services

import (
	"context"
	"fmt"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"

	"github.com/google/uuid"
)

// UserService implements user management business logic
type UserService struct {
	userRepo ports.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// GetByEmail retrieves a user by email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// UpdateProfile updates a user's profile information
func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) (*domain.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Apply updates
	if firstName, ok := updates["firstName"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := updates["lastName"].(string); ok {
		user.LastName = lastName
	}
	if role, ok := updates["role"].(domain.UserRole); ok {
		user.Role = role
	}

	// Update in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpdateSubscriptionStatus updates a user's subscription status
func (s *UserService) UpdateSubscriptionStatus(ctx context.Context, userID uuid.UUID, status domain.SubscriptionStatus) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.SubscriptionStatus = status

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update subscription status: %w", err)
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
