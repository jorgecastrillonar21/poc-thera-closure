package persistence

import (
	"context"
	"fmt"
	"theraclosure/users-service/internal/core/domain"
	"theraclosure/users-service/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type enrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new enrollment repository
func NewEnrollmentRepository(db *gorm.DB) ports.EnrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

// Create creates a new enrollment record
func (r *enrollmentRepository) Create(ctx context.Context, enrollment *domain.EnrollmentData) error {
	if err := r.db.WithContext(ctx).Create(enrollment).Error; err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	return nil
}

// GetByID gets an enrollment record by ID
func (r *enrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EnrollmentData, error) {
	var enrollment domain.EnrollmentData
	if err := r.db.WithContext(ctx).First(&enrollment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("enrollment not found")
		}
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return &enrollment, nil
}

// GetByUserID gets an enrollment record by user ID
func (r *enrollmentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.EnrollmentData, error) {
	var enrollment domain.EnrollmentData
	if err := r.db.WithContext(ctx).First(&enrollment, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("enrollment not found")
		}
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	return &enrollment, nil
}

// Update updates an enrollment record
func (r *enrollmentRepository) Update(ctx context.Context, enrollment *domain.EnrollmentData) error {
	if err := r.db.WithContext(ctx).Save(enrollment).Error; err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}
	return nil
}

// Delete deletes an enrollment record (soft delete)
func (r *enrollmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&domain.EnrollmentData{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}
	return nil
}

// UpdateStep updates a specific enrollment step
func (r *enrollmentRepository) UpdateStep(ctx context.Context, userID uuid.UUID, step int, completed bool) error {
	var updateField string
	switch step {
	case 1:
		updateField = "personal_info_complete"
	case 2:
		updateField = "licensure_details_complete"
	case 3:
		updateField = "practice_info_complete"
	case 4:
		updateField = "admin_setup_complete"
	case 5:
		updateField = "schedule_config_complete"
	default:
		return fmt.Errorf("invalid step number: %d", step)
	}

	if err := r.db.WithContext(ctx).Model(&domain.EnrollmentData{}).
		Where("user_id = ?", userID).
		Update(updateField, completed).Error; err != nil {
		return fmt.Errorf("failed to update enrollment step: %w", err)
	}
	return nil
}

// UpdateEnrollmentStatus updates the enrollment status
func (r *enrollmentRepository) UpdateEnrollmentStatus(ctx context.Context, userID uuid.UUID, status string) error {
	if err := r.db.WithContext(ctx).Model(&domain.EnrollmentData{}).
		Where("user_id = ?", userID).
		Update("enrollment_status", status).Error; err != nil {
		return fmt.Errorf("failed to update enrollment status: %w", err)
	}
	return nil
}

// CompleteEnrollment marks the enrollment as complete
func (r *enrollmentRepository) CompleteEnrollment(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	updates := map[string]interface{}{
		"enrollment_status": "completed",
		"completion_date":   now,
	}

	if err := r.db.WithContext(ctx).Model(&domain.EnrollmentData{}).
		Where("user_id = ?", userID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to complete enrollment: %w", err)
	}
	return nil
}

// GetEnrollmentProgress gets the current and total steps for a user
func (r *enrollmentRepository) GetEnrollmentProgress(ctx context.Context, userID uuid.UUID) (int, int, error) {
	var enrollment domain.EnrollmentData
	if err := r.db.WithContext(ctx).
		Select("current_step, total_steps").
		First(&enrollment, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, fmt.Errorf("enrollment not found")
		}
		return 0, 0, fmt.Errorf("failed to get enrollment progress: %w", err)
	}
	return enrollment.CurrentStep, enrollment.TotalSteps, nil
}
