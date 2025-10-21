package persistence

import (
	"context"
	"fmt"
	"theraclosure/users-service/internal/core/domain"
	"theraclosure/users-service/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

// NewUserProfileRepository creates a new user profile repository
func NewUserProfileRepository(db *gorm.DB) ports.UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

// Create creates a new user profile
func (r *userProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	return nil
}

// GetByID gets a user profile by ID
func (r *userProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	if err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user profile not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return &profile, nil
}

// GetByUserID gets a user profile by user ID
func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	if err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user profile not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return &profile, nil
}

// GetByEmail gets a user profile by email
func (r *userProfileRepository) GetByEmail(ctx context.Context, email string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	if err := r.db.WithContext(ctx).First(&profile, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user profile not found")
		}
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return &profile, nil
}

// Update updates a user profile
func (r *userProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}
	return nil
}

// Delete deletes a user profile (soft delete)
func (r *userProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.UserProfile{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete user profile: %w", err)
	}
	return nil
}

// List lists user profiles with pagination
func (r *userProfileRepository) List(ctx context.Context, limit, offset int) ([]*domain.UserProfile, error) {
	var profiles []*domain.UserProfile
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("failed to list user profiles: %w", err)
	}
	return profiles, nil
}

// UpdateStatus updates the status of a user profile
func (r *userProfileRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserProfile{}).
		Where("user_id = ?", userID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update user profile status: %w", err)
	}
	return nil
}

// MarkProfileComplete marks a user profile as complete
func (r *userProfileRepository) MarkProfileComplete(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserProfile{}).
		Where("user_id = ?", userID).
		Update("profile_complete", true).Error; err != nil {
		return fmt.Errorf("failed to mark profile complete: %w", err)
	}
	return nil
}
