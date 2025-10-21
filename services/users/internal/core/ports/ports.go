package ports

import (
	"context"
	"theraclosure/users-service/internal/core/domain"

	"github.com/google/uuid"
)

// UserProfileRepository defines the interface for user profile data operations
type UserProfileRepository interface {
	// Profile CRUD operations
	Create(ctx context.Context, profile *domain.UserProfile) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.UserProfile, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error)
	GetByEmail(ctx context.Context, email string) (*domain.UserProfile, error)
	Update(ctx context.Context, profile *domain.UserProfile) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.UserProfile, error)

	// Profile status operations
	UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error
	MarkProfileComplete(ctx context.Context, userID uuid.UUID) error
}

// EnrollmentRepository defines the interface for enrollment data operations
type EnrollmentRepository interface {
	// Enrollment CRUD operations
	Create(ctx context.Context, enrollment *domain.EnrollmentData) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EnrollmentData, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.EnrollmentData, error)
	Update(ctx context.Context, enrollment *domain.EnrollmentData) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Enrollment progress operations
	UpdateStep(ctx context.Context, userID uuid.UUID, step int, completed bool) error
	UpdateEnrollmentStatus(ctx context.Context, userID uuid.UUID, status string) error
	CompleteEnrollment(ctx context.Context, userID uuid.UUID) error
	GetEnrollmentProgress(ctx context.Context, userID uuid.UUID) (int, int, error) // current, total
}

// UserService defines the business logic interface for user operations
type UserService interface {
	// Profile management
	CreateProfile(ctx context.Context, profile *domain.UserProfile) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *domain.UserProfile) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error
	ListProfiles(ctx context.Context, limit, offset int) ([]*domain.UserProfile, error)

	// Profile validation and completion
	ValidateProfile(ctx context.Context, profile *domain.UserProfile) error
	CheckProfileCompletion(ctx context.Context, userID uuid.UUID) (bool, error)

	// Search and filtering
	SearchProfiles(ctx context.Context, query string, limit, offset int) ([]*domain.UserProfile, error)
	GetProfilesByStatus(ctx context.Context, status string, limit, offset int) ([]*domain.UserProfile, error)
}

// EnrollmentService defines the business logic interface for enrollment operations
type EnrollmentService interface {
	// Enrollment management
	StartEnrollment(ctx context.Context, userID uuid.UUID, selectedPlan string) error
	GetEnrollment(ctx context.Context, userID uuid.UUID) (*domain.EnrollmentData, error)
	UpdateEnrollment(ctx context.Context, enrollment *domain.EnrollmentData) error

	// Step management
	CompleteStep(ctx context.Context, userID uuid.UUID, step int) error
	GetCurrentStep(ctx context.Context, userID uuid.UUID) (int, error)
	GetProgress(ctx context.Context, userID uuid.UUID) (float64, error) // percentage completion

	// Enrollment completion
	CompleteEnrollment(ctx context.Context, userID uuid.UUID) error
	ValidateEnrollmentCompletion(ctx context.Context, userID uuid.UUID) error

	// Plan management
	UpdateSelectedPlan(ctx context.Context, userID uuid.UUID, plan string) error
	GetUsersbyPlan(ctx context.Context, plan string, limit, offset int) ([]*domain.EnrollmentData, error)
}
