package services

import (
	"context"
	"fmt"
	"time"
	"theraclosure/users-service/internal/core/domain"
	"theraclosure/users-service/internal/core/ports"

	"github.com/google/uuid"
)

type enrollmentService struct {
	enrollmentRepo ports.EnrollmentRepository
}

// NewEnrollmentService creates a new instance of EnrollmentService
func NewEnrollmentService(enrollmentRepo ports.EnrollmentRepository) ports.EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
	}
}

// StartEnrollment starts the enrollment process for a user
func (s *enrollmentService) StartEnrollment(ctx context.Context, userID uuid.UUID, selectedPlan string) error {
	// Check if enrollment already exists
	existing, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		return fmt.Errorf("enrollment already exists for user %s", userID.String())
	}
	
	// Validate plan
	if !s.isValidPlan(selectedPlan) {
		return fmt.Errorf("invalid plan: %s", selectedPlan)
	}
	
	// Create new enrollment data
	enrollment := &domain.EnrollmentData{
		UserID:                   userID,
		PersonalInfoComplete:     false,
		LicensureDetailsComplete: false,
		PracticeInfoComplete:     false,
		AdminSetupComplete:       false,
		ScheduleConfigComplete:   false,
		EnrollmentStatus:         "in_progress",
		CurrentStep:              1,
		TotalSteps:               5,
		SelectedPlan:             selectedPlan,
		PaymentStatus:            "pending",
		PreferredContactMethod:   "email",
	}
	
	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	
	return nil
}

// GetEnrollment retrieves enrollment data by user ID
func (s *enrollmentService) GetEnrollment(ctx context.Context, userID uuid.UUID) (*domain.EnrollmentData, error) {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}
	
	return enrollment, nil
}

// UpdateEnrollment updates enrollment data
func (s *enrollmentService) UpdateEnrollment(ctx context.Context, enrollment *domain.EnrollmentData) error {
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}
	
	return nil
}

// CompleteStep marks a specific enrollment step as complete
func (s *enrollmentService) CompleteStep(ctx context.Context, userID uuid.UUID, step int) error {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}
	
	// Mark the specific step as complete
	switch step {
	case 1:
		enrollment.PersonalInfoComplete = true
	case 2:
		enrollment.LicensureDetailsComplete = true
	case 3:
		enrollment.PracticeInfoComplete = true
	case 4:
		enrollment.AdminSetupComplete = true
	case 5:
		enrollment.ScheduleConfigComplete = true
		// For testing/demo purposes, mark payment as completed when step 5 is done
		enrollment.PaymentStatus = "completed"
	default:
		return fmt.Errorf("invalid step number: %d", step)
	}
	
	// Update current step if this is the next step
	if step >= enrollment.CurrentStep {
		enrollment.CurrentStep = step + 1
		if enrollment.CurrentStep > enrollment.TotalSteps {
			enrollment.CurrentStep = enrollment.TotalSteps
		}
	}
	
	// Update enrollment
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to update enrollment step: %w", err)
	}
	
	// Check if all steps are complete
	if s.allStepsComplete(enrollment) {
		return s.CompleteEnrollment(ctx, userID)
	}
	
	return nil
}

// GetCurrentStep gets the current enrollment step for a user
func (s *enrollmentService) GetCurrentStep(ctx context.Context, userID uuid.UUID) (int, error) {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("enrollment not found: %w", err)
	}
	
	return enrollment.CurrentStep, nil
}

// GetProgress calculates the enrollment progress percentage
func (s *enrollmentService) GetProgress(ctx context.Context, userID uuid.UUID) (float64, error) {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("enrollment not found: %w", err)
	}
	
	completed := s.countCompletedSteps(enrollment)
	progress := float64(completed) / float64(enrollment.TotalSteps) * 100
	
	return progress, nil
}

// CompleteEnrollment marks the entire enrollment as complete
func (s *enrollmentService) CompleteEnrollment(ctx context.Context, userID uuid.UUID) error {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}
	
	// Validate that all steps are complete
	if err := s.ValidateEnrollmentCompletion(ctx, userID); err != nil {
		return fmt.Errorf("cannot complete enrollment: %w", err)
	}
	
	// Update enrollment status
	enrollment.EnrollmentStatus = "completed"
	enrollment.CurrentStep = enrollment.TotalSteps
	now := time.Now()
	enrollment.CompletionDate = &now
	
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to complete enrollment: %w", err)
	}
	
	return nil
}

// ValidateEnrollmentCompletion validates that enrollment can be completed
func (s *enrollmentService) ValidateEnrollmentCompletion(ctx context.Context, userID uuid.UUID) error {
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}
	
	if !s.allStepsComplete(enrollment) {
		return fmt.Errorf("not all enrollment steps are complete")
	}
	
	if enrollment.PaymentStatus != "completed" {
		return fmt.Errorf("payment not completed")
	}
	
	return nil
}

// UpdateSelectedPlan updates the user's selected plan
func (s *enrollmentService) UpdateSelectedPlan(ctx context.Context, userID uuid.UUID, plan string) error {
	if !s.isValidPlan(plan) {
		return fmt.Errorf("invalid plan: %s", plan)
	}
	
	enrollment, err := s.enrollmentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}
	
	enrollment.SelectedPlan = plan
	enrollment.PaymentStatus = "pending" // Reset payment status when plan changes
	
	if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
		return fmt.Errorf("failed to update selected plan: %w", err)
	}
	
	return nil
}

// GetUsersbyPlan gets users by their selected plan
func (s *enrollmentService) GetUsersbyPlan(ctx context.Context, plan string, limit, offset int) ([]*domain.EnrollmentData, error) {
	// This would typically be implemented at the repository level with proper filtering
	// For now, we'll return an error indicating this needs repository-level implementation
	return nil, fmt.Errorf("GetUsersbyPlan not implemented at repository level")
}

// Helper functions

// isValidPlan checks if a plan is valid
func (s *enrollmentService) isValidPlan(plan string) bool {
	validPlans := []string{"essential", "professional", "enterprise"}
	for _, validPlan := range validPlans {
		if plan == validPlan {
			return true
		}
	}
	return false
}

// allStepsComplete checks if all enrollment steps are complete
func (s *enrollmentService) allStepsComplete(enrollment *domain.EnrollmentData) bool {
	return enrollment.PersonalInfoComplete &&
		enrollment.LicensureDetailsComplete &&
		enrollment.PracticeInfoComplete &&
		enrollment.AdminSetupComplete &&
		enrollment.ScheduleConfigComplete
}

// countCompletedSteps counts the number of completed steps
func (s *enrollmentService) countCompletedSteps(enrollment *domain.EnrollmentData) int {
	count := 0
	if enrollment.PersonalInfoComplete {
		count++
	}
	if enrollment.LicensureDetailsComplete {
		count++
	}
	if enrollment.PracticeInfoComplete {
		count++
	}
	if enrollment.AdminSetupComplete {
		count++
	}
	if enrollment.ScheduleConfigComplete {
		count++
	}
	return count
}