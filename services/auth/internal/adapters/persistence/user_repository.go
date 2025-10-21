package persistence

import (
	"context"
	"theraclosure/auth-service/internal/core/domain"
	"theraclosure/auth-service/internal/core/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository implements the user repository using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByCognitoID retrieves a user by Cognito ID
func (r *UserRepository) GetByCognitoID(ctx context.Context, cognitoID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("cognito_id = ?", cognitoID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, err
}
